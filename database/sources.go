package database

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	// "os"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	// "github.com/jackc/pgx/v4/pgxpool"
)
type Source struct {
	ID       int    `json:"-"`        // Source ID (PK)
	Name     string `json:"name"`     // Source name
	Curatifs int    `json:"curatifs"` // Info ouvrage
	SID      int    `json:"-"`        // Infos source_id (FK)

	Created time.Time     `json:"-"`
	// DB      * `json:"-"`
}

// fonction afin de choper tous les postes sources
// pour la page d'accueil
func (src *Source) MenuSource(conn *pgxpool.Conn) ([]*Source, error) {
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

	rows, err := conn.Query(ctx, query)
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

	// fmt.Printf("sources data: %+v\nsources type: %[1]T\n\n", sources)

	return sources, nil
}

func (jsrc *Source) JSource() ([]byte, error) {
	js := []*Source{}

	jsonData, err := json.Marshal(js)
	if err != nil {
		return nil, err
	}

	fmt.Printf("json data: %s\njson data type: %[1]T\n\n", jsonData)

	return jsonData, nil
}

// fonction d'obtention de donnée spécific source
func (src *Source) SourceGet(id int, conn *pgxpool.Conn) (*Source, error) {
	ctx := context.Background()
	query := `
SELECT id, name, created
  FROM source
    WHERE id = $1
`
	sObj := &Source{}
	err := conn.QueryRow(ctx, query, id).Scan(&sObj.ID, &sObj.Name,
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
func (src *Source) SourceInsert(name string, conn *pgxpool.Conn) (int, error) {
	ctx := context.Background()
	query := `
INSERT INTO source (name, created)
VALUES ($1, $2)
  RETURNING id
`
	err := conn.QueryRow(ctx, query, name,
		time.Now().UTC()).Scan(&src.ID)
	if err != nil {
		return 0, nil
	}

	return src.ID, nil
}

// fonction de suppréssion source
func (src *Source) SourceDelete(id int, conn *pgxpool.Conn) error {
	ctx := context.Background()
	query := `
DELETE FROM source
  WHERE id = $1
`
	_, err := conn.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

// Fonction de MaJ source
func (src *Source) SourceUpdate(id int, conn *pgxpool.Conn) error {
	ctx := context.Background()
	query := `
UPDATE source
  SET name = $1
    WHERE id = $2
`
	_, err := conn.Exec(ctx, query, src.Name, id)
	if err != nil {
		return err
	}

	return nil
}
