package database

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"
	"os"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Source struct {
	ID       int    `json:"-"`        // Source ID (PK)
	Name     string `json:"name"`     // Source name
	Curatifs int    `json:"curatifs"` // Info ouvrage
	SID      int    `json:"-"`        // Infos source_id (FK)

	Created time.Time     `json:"-"`
	DB      *pgxpool.Pool `json:"-"`
}

type JsonSource struct {
	ID       int `json:"id"`
	Curatifs int `json:"curatifs"`

	DB *pgxpool.Pool `json:"-"`
}

// fonction afin de choper tous les postes sources
// pour la page d'accueil
func (src *Source) MenuSource() ([]*Source, error) {
	ctx := context.Background()
	query := `
SELECT s.id,
       s.name,
       COUNT(i.status) FILTER (WHERE i.status <> 'archivé')
  FROM source AS s
       LEFT JOIN info AS i ON i.source_id = s.id
  GROUP BY s.id
  ORDER BY name ASC
`

	rows, err := src.DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	sources := []*Source{}

	for rows.Next() {
		sObj := &Source{}

		err := rows.Scan(&sObj.ID, &sObj.Name, &sObj.Curatifs)
		if err != nil {
			return nil, err
		}

		sources = append(sources, sObj)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// jsonData, err := json.Marshal(sources)
	// if err != nil {
	//	return nil, err
	// }

	fmt.Printf("sources data: %+v\nsources type: %[1]T\n\n", sources)

	return sources, nil
}

func (jsrc *JsonSource) JSource(js []*Source) ([]byte, error) {

	jsonData, err := json.Marshal(js)
	if err != nil {
		return nil, err
	}

	fmt.Printf("json data: %s\njson data type: %[1]T\n\n", jsonData)
	f, err := os.Create("source.json")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	_, err = f.Write(jsonData)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}

// fonction d'obtention de donnée spécific source
func (src *Source) SourceGet(id int) (*Source, error) {
	ctx := context.Background()
	query := `
SELECT id, name, created
  FROM source
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
INSERT INTO source (name, created)
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
DELETE FROM source
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
UPDATE source
  SET name = $1
    WHERE id = $2
`
	_, err := src.DB.Exec(ctx, query, src.Name, id)
	if err != nil {
		return err
	}

	return nil
}
