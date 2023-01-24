package database

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Info struct {
	ID       int
	Agent    string
	Material string
	Priority int
	Target   string
	Rte      string
	Detail   string
	Estimate string
	Brips    string
	Oups     string
	Ameps    string
	Ais      string
	SourceID int
	Created  time.Time
	Updated  time.Time
	Status   string
	Event    string
	Doneby   string
	ZeroTime time.Time
	DB       *pgxpool.Pool
}

// Func to create new Info of a Source
func (i *Info) Insert(id int) (int, error) {
	ctx := context.Background()
	query := `
INSERT INTO infos
    (source_id, agent, material, detail, event, priority,
       oups, ameps, brips, rte, ais, estimate,
	target, status, doneby, created)
	  VALUES
	    ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
	      $11, $12, $13, $14, $15, $16)
		RETURNING id;
`
	err := i.DB.QueryRow(ctx, query, id, i.Agent,
		i.Material, i.Detail, i.Event, i.Priority,
		i.Oups, i.Ameps, i.Brips, i.Rte, i.Ais,
		i.Estimate, i.Target, i.Status,
		i.Doneby, time.Now().UTC()).Scan(&i.ID)
	if err != nil {
		return 0, err
	}

	return i.ID, nil
}

// Func to gather data from infos table
func (i *Info) InfoGet(id int) (*Info, error) {
	ctx := context.Background()
	query := `
SELECT *
  FROM infos
    WHERE id = $1
`

	row := i.DB.QueryRow(ctx, query, id)

	var rte, ameps, ais, brips, oups, estimate, target,
		doneby *string
	// var priority *int
	var updated *time.Time

	iObj := &Info{}

	err := row.Scan(&iObj.ID, &iObj.Agent, &iObj.Material,
		&iObj.Priority, &rte, &iObj.Detail, &estimate, &brips,
		&oups, &ameps, &ais, &iObj.SourceID, &iObj.Created,
		&updated, &iObj.Status, &iObj.Event, &target, &doneby)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	iObj.ZeroTime = time.Date(0001, time.January,
		1, 0, 0, 0, 0, time.UTC)


	// if agent != nil {
	//	iObj.Agent = *agent
	// }

	// if material != nil {
	//	iObj.Material = *material
	// }

	// if priority != nil {
	//	iObj.Priority = *priority
	// }

	if target != nil {
		iObj.Target = *target
	}

	if rte != nil {
		iObj.Rte = *rte
	}

	// if detail != nil {
	//	iObj.Detail = *detail
	// }

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

	// if status != nil {
	//	iObj.Status = *status
	// }

	// if event != nil {
	//	iObj.Event = *event
	// }

	// if created != nil {
	//	iObj.Created = *created
	// }

	if updated != nil {
		iObj.Updated = *updated
	}

	if estimate != nil {
		iObj.Estimate = *estimate
	}

	if doneby != nil {
		iObj.Doneby = *doneby
	}

	return iObj, nil
}

func (i *Info) InfoList(id int) ([]*Info, error) {
	ctx := context.Background()
	stmt := `
SELECT id, material, created, status, source_id, priority
    FROM infos
	WHERE source_id = $1
	    ORDER BY
		created ASC
`
	rows, err := i.DB.Query(ctx, stmt, id)
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

func(i *Info) InfoUpdate(id int) error {
	ctx := context.Background()
	query := `
UPDATE infos
  SET agent = $1,
      material = $2,
      priority = $3,
      target = $4,
      rte = $5,
      detail = $6,
      estimate = $7,
      brips = $8,
      oups = $9,
      ameps = $10,
      ais = $11,
      updated = $12,
      status = $13,
      event = $14,
      doneby = $15
	WHERE id = $16
`
	_, err := i.DB.Exec(ctx, query, i.Agent, i.Material,
		i.Priority, i.Target, i.Rte, i.Detail, i.Estimate,
		i.Brips, i.Oups, i.Ameps, i.Ais, time.Now().UTC(),
		i.Status, i.Event, i.Doneby, id)
	if err != nil {
		return err
	}

	return nil
}

func(i *Info) InfoUp(id int) error {
	ctx := context.Background()
	query := `
UPDATE infos
  SET material = $1,
      updated = $2
	WHERE id = $3
`
	_, err := i.DB.Exec(ctx, query, i.Material,
		time.Now().UTC(), id)
	if err != nil {
		return err
	}

	return nil
}
