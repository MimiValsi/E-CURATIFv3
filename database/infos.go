package database

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Info struct {
	ID         int // primary key
	Agent      string
	Material   string
	Priority   int
	Target     string
	Rte        string
	Detail     string
	Estimate   string
	Brips      string
	Oups       string
	Ameps      string
	Ais        string
	SourceID   int // foreign key en référence au PK de source
	Created    time.Time
	Updated    time.Time
	Status     string
	Event      string
	Doneby     string
	Pilot      string
	ActionDate string
	DayDone    string
	ZeroTime   time.Time
	DB         *pgxpool.Pool
}

// Fonction de création donnée info
func (i *Info) Insert(id int) (int, error) {
	ctx := context.Background()
	query := `
INSERT INTO infos
    (source_id, agent, material, detail, event, priority, oups, ameps,
       brips, rte, ais, estimate, target, status, doneby, created)
	  VALUES
	    ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13,
	      $14, $15, $16, $17, $18)
		RETURNING id;
`
	err := i.DB.QueryRow(ctx, query, id, i.Agent,
		i.Material, i.Detail, i.Event, i.Priority,
		i.Oups, i.Ameps, i.Brips, i.Rte, i.Ais,
		i.Estimate, i.Target, i.Status,
		i.Doneby, i.ActionDate,
		time.Now().UTC()).Scan(&i.ID)
	if err != nil {
		return 0, err
	}

	return i.ID, nil
}

// Fonction d'obtention de donnée spécific
func (i *Info) InfoGet(id int) (*Info, error) {
	ctx := context.Background()
	query := `
SELECT id, agent, material, priority, rte, detail, estimate, brips,
  oups, ameps, ais, source_id, created, updated, status, event, target
    FROM infos
      WHERE id = $1
`
	var rte, ameps, ais, brips, oups, estimate, target,
		doneby, pilot, actionDate *string
	var updated *time.Time

	iObj := &Info{}
	err := i.DB.QueryRow(ctx, query, id).Scan(&iObj.ID, &iObj.Agent,
		&iObj.Material, &iObj.Priority, &rte, &iObj.Detail,
		&estimate, &brips, &oups, &ameps, &ais, &iObj.SourceID,
		&iObj.Created, &updated, &iObj.Status, &iObj.Event,
		&target)
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
func (i *Info) InfoList(id int) ([]*Info, error) {
	ctx := context.Background()
	query := `
SELECT id, material, created, status, source_id, priority
  FROM infos
    WHERE source_id = $1
      ORDER BY
	priority ASC
`
	rows, err := i.DB.Query(ctx, query, id)
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

func (i *Info) InfoDelete(id int) error {
	ctx := context.Background()
	query := `
DELETE FROM infos
    WHERE id = $1
`
	_, err := i.DB.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

// Fonction de mise à jour donnée
func (i *Info) InfoUpdate(id int) error {
	ctx := context.Background()
	query := `
UPDATE infos
  SET agent = $1, material = $2, priority = $3, target = $4, rte = $5,
    detail = $6, estimate = $7, brips = $8, oups = $9, ameps = $10,
      ais = $11, updated = $12, status = $13, event = $14, doneby = $15,
	pilot = $16, action_date = $17
	  WHERE id = $18
`
	_, err := i.DB.Exec(ctx, query, i.Agent, i.Material,
		i.Priority, i.Target, i.Rte, i.Detail, i.Estimate,
		i.Brips, i.Oups, i.Ameps, i.Ais, time.Now().UTC(),
		i.Status, i.Event, i.Doneby, i.Pilot, i.ActionDate, id)
	if err != nil {
		return err
	}

	return nil
}

// Fonction de test
func (i *Info) InfoUp(id int) error {
	ctx := context.Background()
	query := `
UPDATE infos
  SET material = $1, updated = $2
    WHERE id = $3
`
	_, err := i.DB.Exec(ctx, query, i.Material,
		time.Now().UTC(), id)
	if err != nil {
		return err
	}

	return nil
}
