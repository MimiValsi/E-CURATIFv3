package database

import (
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Source struct {
	ID        int    `json:"-"`                   // Source ID (PK)
	SID       int    `json:"-"`                   // Infos source_id (FK)
	Average   int    `json:"average"`             // Average 'Réalisée/résolu'
	ARealiser int    `json:"a_realiser"`          // Curatifs 'a réaliser'
	EnCours   int    `json:"en_cours"`            // Curatifs 'en cours/affecté'
	Done      int    `json:"done"`                // Curatifs 'Réalisée/résolu'
	Curatifs  int    `json:"curatifs"`            // Info ouvrage
	Diff      int    `json:"diff"`                // Calculate difference between Done and not done (except archives)
	Name      string `json:"name,omitempty"`      // Source name
	CodeGMAO  string `json:"code_GMAO,omitempty"` // CodeGMAO to be display in the graph
}

// fonction afin de choper tous les postes sources
// pour la page d'accueil
func (src *Source) MenuSource(conn *pgxpool.Conn) ([]*Source, error) {
	ctx := context.Background()
	query := `
SELECT s.id,
       s.name,
       s.code_GMAO,
       COUNT(i.status) FILTER (WHERE i.status = 'a réaliser' OR i.status = 'en attente') as a_realiser,
       COUNT(i.status) FILTER (WHERE i.status = 'en cours' OR i.status = 'affecté') as en_cours,
       COUNT(i.status) FILTER (WHERE i.status = 'réalisée' OR i.status = 'résolu') as done,
       COUNT(i.status) FILTER (WHERE i.status <> 'archivée') as all
  FROM source AS s
       LEFT JOIN info AS i 
       ON i.source_id = s.id
 GROUP BY s.id
 ORDER BY name ASC
`
	// AVG((status = 'Réalisée' OR status = 'résolu')::int)::numeric(2,2)*100 as average,

	rows, err := conn.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	sources := []*Source{}

	for rows.Next() {
		sObj := &Source{}

		var all *int
		err = rows.Scan(&sObj.ID, &sObj.Name, &sObj.CodeGMAO,
			&sObj.ARealiser, &sObj.EnCours, &sObj.Done, &all)
		if err != nil {
			return nil, err
		}

		if *all > 0 {
			sObj.Curatifs = *all
			tmp := float64(sObj.Done) / float64(*all)
			sObj.Average = int(tmp * 100)
		} else {
			sObj.Curatifs = *all
			sObj.Average = 100
		}

		sObj.Diff = *all - (sObj.Done)

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
	_, err := conn.Exec(ctx, query, src.Name, strings.ToUpper(src.CodeGMAO), id)
	if err != nil {
		return err
	}

	return nil
}
