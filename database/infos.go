package database

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Info struct {
        ID         int `json:"-"` // primary key
        Priority   int `json:"priority,omitempty"`
        SourceID   int `json:"-"` // foreign key en référence au PK de source
        Counter    int `json:"counter,omitempty"`
        Agent      string `json:"agent,omitempty"`
        Material   string `json:"material,omitempty"`
        Target     string `json:"target,omitempty"`
        Rte        string `json:"rte,omitempty"`
	Detail     string `json:"detail,omitempty"`
	Estimate   string `json:"estimate,omitempty"`
	Brips      string `json:"brips,omitempty"`
	Oups       string `json:"oups,omitempty"`
	Ameps      string `json:"ameps,omitempty"`
	Ais        string `json:"ais,omitempty"`
	Status     string `json:"status,omitempty"`
	Event      string `json:"event,omitempty"`
	Doneby     string `json:"doneby,omitempty"`
	Pilot      string `json:"pilot,omitempty"`
	ActionDate string `json:"actionDate,omitempty"`
	DayDone    string `json:"dayDone,omitempty"`

	ZeroTime time.Time `json:"-"`
	Created  time.Time `json:"-"`
	Updated  time.Time `json:"-"`
}

// Gather all priority 1 infos
func (i *Info) PriorityInfos(conn *pgxpool.Conn) ([]*Info, error) {
        ctx := context.Background() 

        query := `
SELECT i.material, 
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

                err := rows.Scan(&iObj.Material, &iObj.Detail)
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
INSERT INTO info (source_id, agent, material, detail, 
	   	  event, priority, oups, ameps,
       		  brips, rte, ais, estimate, 
		  target, status, doneby, created)
VALUES ($1,  $2,  $3,  $4, 
	$5,  $6,  $7,  $8, 
	$9,  $10, $11, $12, 
	$13, $14, $15, $16)
  RETURNING id;
`
	err := conn.QueryRow(ctx, query, id, i.Agent,
		i.Material, i.Detail, i.Event, i.Priority,
		i.Oups, i.Ameps, i.Brips, i.Rte, i.Ais,
		i.Estimate, i.Target, i.Status,
		i.Doneby,
		time.Now().UTC()).Scan(&i.ID)
	if err != nil {
		return 0, err
	}

	return i.ID, nil
}

// Fonction d'obtention de donnée spécific
func (i *Info) InfoGet(id int, conn *pgxpool.Conn) (*Info, error) {
	ctx := context.Background()
	query := `
SELECT id, agent, material, priority, 
       rte, detail, estimate, brips,
       oups, ameps, ais, source_id, 
       created, updated, status, event, target, doneby
  FROM info
 WHERE id = $1 AND 
 status <> 'résolu'
`
	var rte, ameps, ais, brips, oups, estimate, target,
		doneby, pilot, actionDate *string
	var updated *time.Time

	iObj := &Info{}
	err := conn.QueryRow(ctx, query, id).Scan(&iObj.ID, &iObj.Agent,
		&iObj.Material, &iObj.Priority, &rte, &iObj.Detail,
		&estimate, &brips, &oups, &ameps, &ais, &iObj.SourceID,
		&iObj.Created, &updated, &iObj.Status, &iObj.Event,
		&target, &iObj.Doneby)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	iObj.ZeroTime = time.Date(0001, time.January,
		1, 0, 0, 0, 0, time.UTC)

	if target != nil {
		iObj.Target = *target
	}

	if rte != nil {
		iObj.Rte = *rte
	}

	if ameps != nil {
		iObj.Ameps = *ameps
	}

	if ais != nil {
		iObj.Ais = *ais
	}

	if brips != nil {
		iObj.Brips = *brips
	}

	if oups != nil {
		iObj.Oups = *oups
	}

	if updated != nil {
		iObj.Updated = *updated
	}

	if estimate != nil {
		iObj.Estimate = *estimate
	}

	if doneby != nil {
		iObj.Doneby = *doneby
	}

	if pilot != nil {
		iObj.Pilot = *pilot
	}

	if actionDate != nil {
		iObj.ActionDate = *actionDate
	}

	return iObj, nil
}

// Fonction afin de choper plusieurs données et la transférer
// dans un slice
func (i *Info) InfoList(id int, conn *pgxpool.Conn) ([]*Info, error) {
	ctx := context.Background()
	query := `
SELECT id, material, created, 
       status, source_id, priority
  FROM info
 WHERE source_id = $1 AND
 status <> 'archivé'
 ORDER BY priority ASC
`
	rows, err := conn.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	infos := []*Info{}

	for rows.Next() {
		iObj := &Info{}

		err = rows.Scan(&iObj.ID, &iObj.Material,
			&iObj.Created, &iObj.Status,
			&iObj.SourceID, &iObj.Priority)
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

func (i *Info) InfoDelete(id int, conn *pgxpool.Conn) error {
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
func (i *Info) InfoUpdate(id int, conn *pgxpool.Conn) error {
	ctx := context.Background()
	query := `
UPDATE info
   SET agent = $1, material = $2, priority = $3, target = $4, rte = $5,
       detail = $6, estimate = $7, brips = $8, oups = $9, ameps = $10,
       ais = $11, updated = $12, status = $13, event = $14, doneby = $15
 WHERE id = $16
`
	_, err := conn.Exec(ctx, query, i.Agent, i.Material,
		i.Priority, i.Target, i.Rte, i.Detail, i.Estimate,
		i.Brips, i.Oups, i.Ameps, i.Ais, time.Now().UTC(),
		i.Status, i.Event, i.Doneby, id)
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
   SET material = $1, updated = $2
 WHERE id = $3
`
	_, err := conn.Exec(ctx, query, i.Material,
		time.Now().UTC(), id)
	if err != nil {
		return err
	}

	return nil
}
