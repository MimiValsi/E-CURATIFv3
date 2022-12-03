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
	ZeroTime time.Time
	DB       *pgxpool.Pool
}

func (i *Info) Insert(id int) (int, error) {
	ctx := context.Background()
	query := `
INSERT INTO infos
    (source_id, agent, material, detail, event, priority, status, oups ,created)
        VALUES
            ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	        RETURNING id;
`
	// i := &Info{}
	err := i.DB.QueryRow(ctx, query, id, i.Agent,
		i.Material, i.Detail, i.Event,
		i.Priority, i.Status, i.Oups, time.Now().UTC()).Scan(&i.ID)
	if err != nil {
		return 0, err
	}

	return i.ID, nil
}

func (i *Info) InfoGet(id int) (*Info, error) {
	ctx := context.Background()
	query := `
SELECT  id, agent, material, priority,
	target, rte, detail,
	ameps, ais, brips,
	oups, status, event,
	source_id, created, updated, estimate
    	    FROM infos
		WHERE id = $1
`

	row := i.DB.QueryRow(ctx, query, id)

	var agent, material, rte, detail, ameps, ais, brips, oups, status, event, estimate, target *string
	var info_id, priority, source_id *int
	var created, updated *time.Time

	iObj := &Info{}

	err := row.Scan(&info_id, &agent, &material, &priority,
		&target, &rte, &detail,
		&ameps, &ais, &brips,
		&oups, &status, &event,
		&source_id, &created, &updated, &estimate)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	iObj.ZeroTime = time.Date(0001, time.January, 1, 0, 0, 0, 0, time.UTC)

	if info_id != nil {
		iObj.ID = *info_id
	}

	if agent != nil {
		iObj.Agent = *agent
	}

	if material != nil {
		iObj.Material = *material
	}

	if priority != nil {
		iObj.Priority = *priority
	}

	if target != nil {
		iObj.Target = *target
	}

	if rte != nil {
		iObj.Rte = *rte
	}

	if detail != nil {
		iObj.Detail = *detail
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

	if status != nil {
		iObj.Status = *status
	}

	if event != nil {
		iObj.Event = *event
	}

	if source_id != nil {
		iObj.SourceID = *source_id
	}

	if created != nil {
		iObj.Created = *created
	}

	if updated != nil {
		iObj.Updated = *updated
	}

	if estimate != nil {
		iObj.Estimate = *estimate
	}

	return iObj, nil
}

func (i *Info) InfoList(sId int) ([]*Info, error) {
	ctx := context.Background()
	stmt := `
SELECT id, material, created, status, source_id 
    FROM infos
        WHERE source_id = $1
            ORDER BY
	        created ASC;
`
	rows, err := i.DB.Query(ctx, stmt, sId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	infos := []*Info{}

	for rows.Next() {
		iObj := &Info{}

		// var infoID *int

		err = rows.Scan(&iObj.ID, &iObj.Material, &iObj.Created, &iObj.Status, &iObj.SourceID)
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
