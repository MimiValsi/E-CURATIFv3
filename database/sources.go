package database

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Source struct {
	ID      int
	Name    string
	Created time.Time
	DB      *pgxpool.Pool
}

func (s *Source) MenuSource() ([]*Source, error) {
	ctx := context.Background()
	query := `
SELECT id, name, created
    FROM sources
	ORDER BY
	    name ASC
`
	rows, err := s.DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sources := []*Source{}

	for rows.Next() {
		sObj := &Source{}

		err := rows.Scan(&sObj.ID, &sObj.Name, &sObj.Created)
		if err != nil {
			return nil, err
		}

		sources = append(sources, sObj)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return sources, nil
}

func (s *Source) SourceGet(id int) (*Source, error) {
	ctx := context.Background()
	stmt := `
SELECT *
    FROM sources
	WHERE
	    id = $1
`

	row := s.DB.QueryRow(ctx, stmt, id)

	sObj := &Source{}

	err := row.Scan(&sObj.ID, &sObj.Name, &sObj.Created)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	return sObj, nil

}

func (s *Source) SourceInsert(name string) (int, error) {
	ctx := context.Background()
	query := `
INSERT INTO sources
    (name, created)
	VALUES ($1, $2)
	    RETURNING id
`
	err := s.DB.QueryRow(ctx, query, name,
		time.Now().UTC()).Scan(&s.ID)
	if err != nil {
		return 0, nil
	}

	return s.ID, nil
}

func (s *Source) SourceDelete(id int) error {
	ctx := context.Background()
	query := `
DELETE FROM sources
    WHERE id = $1
`
	_, err := s.DB.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *Source) SourceUpdate(id int) error {
	ctx := context.Background()
	query := `
UPDATE sources
  SET name = $1
    WHERE id = $2
`
	_, err := s.DB.Exec(ctx, query, s.Name, id)
	if err != nil {
		return err
	}

	return nil
}
