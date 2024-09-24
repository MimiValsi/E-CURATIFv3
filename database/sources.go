package database

import (
	"context"
	"errors"
	"strings"
	// "time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Source struct {
	ID       int    `json:"-"`                 // Source ID (PK)
	Name     string `json:"name,omitempy"`     // Source name
	Curatifs int    `json:"curatifs,omitempy"` // Info ouvrage
	CodeGMAO string `json:"code_GMAO,omitempy"`
	SID      int    `json:"-"` // Infos source_id (FK)
}

// fonction afin de choper tous les postes sources
// pour la page d'accueil
func (src *Source) MenuSource(conn *pgxpool.Conn) ([]*Source, error) {
	ctx := context.Background()
	query := `
SELECT s.id,
       s.name,
       s.code_GMAO,
       COUNT(i.status) FILTER (WHERE i.status <> 'Archivés')
  FROM source AS s
       LEFT JOIN info AS i 
       ON i.source_id = s.id
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

		err = rows.Scan(&sObj.ID, &sObj.Name, &sObj.CodeGMAO, &sObj.Curatifs)
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

func (src *Source) CuratifsDone(conn *pgxpool.Conn) ([]*Source, error) {
	ctx := context.Background()
	query := `
SELECT s.id,
       s.name,
       s.code_GMAO,
       COUNT(i.status) FILTER (WHERE i.status = 'Réalisée' AND i.status = 'résolu')
  FROM source AS s
       LEFT JOIN info AS i 
       ON i.source_id = s.id
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

		err = rows.Scan(&sObj.ID, &sObj.Name, &sObj.CodeGMAO, &sObj.Curatifs)
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
func (src *Source) Get(id int, conn *pgxpool.Conn) (*Source, error) {
	ctx := context.Background()
	query := `
SELECT id, name
  FROM source
 WHERE id = $1
`
	sObj := &Source{}
	err := conn.QueryRow(ctx, query, id).Scan(&sObj.ID, &sObj.Name)
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
func (src *Source) Insert(name string, codeGmao string, conn *pgxpool.Conn) (int, error) {
	ctx := context.Background()
	query := `
INSERT INTO source (name, code_gmao)
VALUES ($1, $2)
  RETURNING id
`
	err := conn.QueryRow(ctx, query, strings.ToUpper(name), strings.ToUpper(codeGmao)).Scan(&src.ID)
	if err != nil {
		return 0, nil
	}

	return src.ID, nil
}

// fonction de suppréssion source
func (src *Source) Delete(id int, conn *pgxpool.Conn) error {
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
func (src *Source) Update(id int, conn *pgxpool.Conn) error {
	ctx := context.Background()
	query := `
UPDATE source
   SET name = $1,
       code_gmao = $2
 WHERE id = $3
`
	tmp := strings.ToUpper(src.CodeGMAO)
	_, err := conn.Exec(ctx, query, src.Name, tmp, id)
	if err != nil {
		return err
	}

	return nil
}
