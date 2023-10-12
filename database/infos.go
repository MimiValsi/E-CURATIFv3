package database

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Info struct {
	ID          int    `json:"-"` // primary key
	Priorite    int    `json:"priorite,omitempty"`
	SourceID    int    `json:"-"` // foreign key en référence au PK de source
	Counteur    int    `json:"counteur,omitempty"`
	Agent       string `json:"agent,omitempty"`
	Ouvrage     string `json:"ouvrage,omitempty"`
	DatePrevue  string `json:"date_prevue,omitempty"`
	Detail      string `json:"detail,omitempty"`
	Devis       string `json:"devis,omitempty"`
	// Oups        string `json:"oups,omitempty"`
	Status      string `json:"status,omitempty"`
	Evenement   string `json:"evenement,omitempty"`
	FaitPar     string `json:"fait_par,omitempty"`
	Pilote      string `json:"pilote,omitempty"`
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

		err := rows.Scan(&iObj.Ouvrage, &iObj.Detail)
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
	   	  evenement, priorite, oups, 
       		  devis, date_prevue, status, fait_par, 
                  commentaire, created)
VALUES ($1,  $2,  $3,  $4, 
	$5,  $6,  $7,  $8, 
	$9,  $10, $11, $12, $13)
  RETURNING id;
`
	err := conn.QueryRow(ctx, query, id, i.Agent,
		i.Ouvrage, i.Detail, i.Evenement, i.Priorite,
		i.Oups, i.Devis, i.DatePrevue, i.Status,
		i.FaitPar, i.Commentaire,
		time.Now().UTC()).Scan(&i.ID)
	if err != nil {
		return 0, err
	}

	return i.ID, nil
}

// Fonction d'obtention de donnée spécific
func (i *Info) Get(id int, conn *pgxpool.Conn) (*Info, error) {
	ctx := context.Background()
	query := `
SELECT id, agent, ouvrage, priorite, 
       detail, devis, oups, source_id, 
       created, updated, status, evenement, fait_par, date_prevue,
       commentaire
  FROM info
 WHERE id = $1
`
	var date_prevue, oups, devis, fait_par, commentaire *string
	var updated *time.Time

	iObj := &Info{}
	err := conn.QueryRow(ctx, query, id).Scan(&iObj.ID, &iObj.Agent,
		&iObj.Ouvrage, &iObj.Priorite, &iObj.Detail,
		&devis, &oups, &iObj.SourceID,
		&iObj.Created, &updated, &iObj.Status, &iObj.Evenement,
		&fait_par, &date_prevue, &commentaire)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	iObj.ZeroTime = time.Date(0001, time.January,
		1, 0, 0, 0, 0, time.UTC)

	if date_prevue != nil {
		iObj.DatePrevue = *date_prevue
	}

	if oups != nil {
		iObj.Oups = *oups
	}

	if updated != nil {
		iObj.Updated = *updated
	}

	if devis != nil {
		iObj.Devis = *devis
	}

	if fait_par != nil {
		iObj.FaitPar = *fait_par
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
   SET agent = $1, ouvrage = $2, priorite = $3, date_prevue = $4,
       detail = $5, devis = $6,oups = $7,
       updated = $8, status = $9, evenement = $10, fait_par = $11,
       commentaire = $12
 WHERE id = $13
`
	_, err := conn.Exec(ctx, query, i.Agent, i.Ouvrage, i.Priorite,
		i.DatePrevue, i.Detail, i.Devis,
		i.Oups, time.Now().UTC(),
		i.Status, i.Evenement, i.FaitPar,
		i.Commentaire, id)
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
