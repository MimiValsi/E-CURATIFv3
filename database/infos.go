package database

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Info struct {
	ID          int    `json:"-"` // primary key
	Priorite    int    `json:"priorite,omitempty"`
	SourceID    int    `json:"-"` // foreign key en référence au PK de source
	Agent       string `json:"agent,omitempty"`
	Ouvrage     string `json:"ouvrage,omitempty"`
	Echeance    string `json:"echeance,omitempty"`
	Detail      string `json:"detail,omitempty"`
	Status      string `json:"status,omitempty"`
	Evenement   string `json:"evenement,omitempty"`
	Commentaire string `json:"commentaire,omitempty"`

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
 WHERE status <> 'résolu' AND 
       status <> 'archivé'
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
INSERT INTO info (source_id, agent, ouvrage, detail, 
	   	  priorite, status, evenement, created)
VALUES ($1,  $2,  $3,  $4, 
	$5,  $6,  $7,  $8)
  RETURNING id;
`
	if i.Created == i.ZeroTime {
		i.Created = time.Now().UTC()
	}

	err := conn.QueryRow(ctx, query, id, i.Agent,
		i.Ouvrage, i.Detail, i.Priorite,
		i.Echeance, i.Status, i.Evenement,
		i.Created).Scan(&i.ID)
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
		iObj.Status = *status
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
SELECT id, ouvrage, created, 
       status, source_id, priorite
  FROM info
 WHERE source_id = $1 AND
 status <> 'archivé'
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

		err = rows.Scan(&iObj.ID, &iObj.Ouvrage,
			&iObj.Created, &iObj.Status,
			&iObj.SourceID, &iObj.Priorite)
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
   SET agent = $1, ouvrage = $2, priorite = $3, echeance = $4,
       detail = $5, updated = $6, status = $7, evenement = $8, commentaire = $9
 WHERE id = $10
`
	_, err := conn.Exec(ctx, query, i.Agent, i.Ouvrage, i.Priorite,
		i.Echeance, i.Detail,
		time.Now().UTC(), i.Status, i.Evenement, i.Commentaire, id)
	if err != nil {
		return err
	}

	return nil
}

// Fonction de test
func (i *Info) InfoUp(id int, conn *pgxpool.Conn) error {
	ctx := context.Background()
	query := `
UPDATE info
   SET ouvrage = $1, updated = $2
 WHERE id = $3
`
	_, err := conn.Exec(ctx, query, i.Ouvrage,
		time.Now().UTC(), id)
	if err != nil {
		return err
	}

	return nil
}
