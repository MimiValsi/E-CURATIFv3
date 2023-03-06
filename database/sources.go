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

// fonction afin de choper tous les postes sources
// pour la page d'accueil
func (src *Source) MenuSource() ([]*Source, error) {
	ctx := context.Background()
	query := `
SELECT id, name, created
  FROM sources
    ORDER BY
      name ASC
`

	rows, err := src.DB.Query(ctx, query)
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

// fonction d'obtention de donnée spécific source
func (src *Source) SourceGet(id int) (*Source, error) {
	ctx := context.Background()
	query := `
SELECT *
  FROM sources
    WHERE id = $1
`
	sObj := &Source{}
	err := src.DB.QueryRow(ctx, query, id).Scan(&sObj.ID, &sObj.Name,
		&sObj.Created)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	return sObj, nil
}

// fonction de création donnée source
func (src *Source) SourceInsert(name string) (int, error) {
	ctx := context.Background()
	query := `
INSERT INTO sources
  (name, created)
    VALUES ($1, $2)
      RETURNING id
`
	err := src.DB.QueryRow(ctx, query, name,
		time.Now().UTC()).Scan(&src.ID)
	if err != nil {
		return 0, nil
	}

	return src.ID, nil
}

// fonction de suppréssion source
func (src *Source) SourceDelete(id int) error {
	ctx := context.Background()
	query := `
DELETE FROM sources
    WHERE id = $1
`
	_, err := src.DB.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}


// Fonction de MaJ source
func (src *Source) SourceUpdate(id int) error {
	ctx := context.Background()
	query := `
UPDATE sources
  SET name = $1
    WHERE id = $2
`
	_, err := src.DB.Exec(ctx, query, src.Name, id)
	if err != nil {
		return err
	}

	return nil
}
