package database

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"E-CURATIFv3/internal/validator"
)

type Info struct {
	ID          int    `json:"-"` // primary key
	Priorite    int    `json:"priorite"`
	SourceID    int    `json:"-"` // foreign key en référence au PK de source
	Agent       string `json:"agent"`
	Ouvrage     string `json:"ouvrage"`
	Echeance    string `json:"echeance"`
	Detail      string `json:"detail"`
	Status      string `json:"status"`
	Evenement   string `json:"evenement"`
	Commentaire string `json:"commentaire"`
	Entite      string `json:"entite"`

	ZeroTime time.Time `json:"-"`
	Created  time.Time `json:"-"`
	Updated  time.Time `json:"-"`
}

// Gather all priorite 1 infos
func (i *Info) PriorityInfos(conn *pgxpool.Conn) ([]*Info, error) {
	ctx := context.Background()

	query := `
SELECT i.ouvrage,
       i.detail
  FROM info AS i
 WHERE status NOT LIKE 'réalis%' AND status NOT LIKE 'résol%' AND 
       status NOT LIKE 'archiv%'
`

	rows, err := conn.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	infos := []*Info{}

	for rows.Next() {
		iObj := &Info{}

		err = rows.Scan(&iObj.Ouvrage, &iObj.Detail)
		if err != nil {
			return nil, err
		}

		infos = append(infos, iObj)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return infos, nil
}

// Fonction de création donnée info
func (i *Info) Insert(id int, conn *pgxpool.Conn) (int, error) {
	ctx := context.Background()
	query := `
INSERT INTO info (source_id, ouvrage, detail,
		  evenement, priorite, echeance, status,
		  created, entite, commentaire)
VALUES ($1,  $2,  $3,  $4, $5,
	$6,  $7,  $8, $9, $10)
  RETURNING id;
`
	err := conn.QueryRow(ctx, query, id,
		i.Ouvrage, i.Detail, i.Evenement, i.Priorite,
		i.Echeance, strings.ToLower(i.Status), i.Created,
		i.Entite, i.Commentaire).Scan(&i.ID)
	if err != nil {
		return 0, err
	}

	return i.ID, nil
}

// Fonction d'obtention de donnée spécific
func (i *Info) Get(id int, conn *pgxpool.Conn) (*Info, error) {
	ctx := context.Background()
	query := `
SELECT id, ouvrage, priorite, 
       detail, source_id, 
       created, updated, status, evenement, echeance,
       commentaire
  FROM info
 WHERE id = $1
`
	var ouvrage, priorite, detail, status, evenement, echeance, commentaire *string
	var updated *time.Time

	iObj := &Info{}
	err := conn.QueryRow(ctx, query, id).Scan(&iObj.ID, &ouvrage,
		&priorite, &detail, &iObj.SourceID, &iObj.Created, &updated, &status,
		&evenement, &echeance, &commentaire)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	iObj.ZeroTime = time.Date(0o001, time.January,
		1, 0, 0, 0, 0, time.UTC)

	if ouvrage != nil {
		iObj.Ouvrage = *ouvrage
	}

	if priorite != nil {
		iObj.Priorite, _ = strconv.Atoi(*priorite)
	}
	if detail != nil {
		iObj.Detail = *detail
	}

	if status != nil {
		iObj.Status = validator.ToCapital(*status)
	}

	if evenement != nil {
		iObj.Evenement = *evenement
	}

	if echeance != nil {
		iObj.Echeance = *echeance
	}

	if updated != nil {
		iObj.Updated = *updated
	}

	if commentaire != nil {
		iObj.Commentaire = *commentaire
	}

	return iObj, nil
}

// Fonction afin de choper plusieurs données et la transférer
// dans un slice
func (i *Info) List(id int, conn *pgxpool.Conn) ([]*Info, error) {
	ctx := context.Background()
	query := `
SELECT id, ouvrage, detail, 
       status, source_id, priorite
  FROM info
 WHERE source_id = $1 AND
 status NOT LIKE 'archiv%' AND status NOT LIKE 'réalis%' AND status NOT LIKE 'résol%'
 ORDER BY priorite ASC
`
	rows, err := conn.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	infos := []*Info{}

	for rows.Next() {
		iObj := &Info{}

		var status, ouvrage, priorite, detail *string
		err = rows.Scan(&iObj.ID, &ouvrage,
			&detail, &status,
			&iObj.SourceID, &priorite)
		if err != nil {
			return nil, err
		}

		if status != nil {
			iObj.Status = validator.ToCapital(*status)
		}

		if ouvrage != nil {
			iObj.Ouvrage = *ouvrage
		}

		if priorite != nil {
			iObj.Priorite, _ = strconv.Atoi(*priorite)
		}

		if detail != nil {
			iObj.Detail = *detail
		}

		infos = append(infos, iObj)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return infos, nil
}

func (i *Info) Delete(id int, conn *pgxpool.Conn) error {
	ctx := context.Background()
	query := `
DELETE FROM info
 WHERE id = $1
`
	_, err := conn.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

// Fonction de mise à jour donnée
func (i *Info) Update(id int, conn *pgxpool.Conn) error {
	ctx := context.Background()
	query := `
UPDATE info
   SET ouvrage = $1, detail = $2, evenement = $3, priorite = $4, 
     echeance = $5, status = $6, updated = $7, entite = $8, commentaire = $9
 WHERE id = $10
`
	_, err := conn.Exec(ctx, query, i.Ouvrage, i.Detail, i.Evenement,
		i.Priorite, i.Echeance, i.Status, time.Now().UTC(),
		i.Entite, i.Commentaire, id)
	if err != nil {
		return err
	}

	return nil
}
