package database

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"E-CURATIFv3/internal/validator"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/text/encoding/charmap"
)

type Export struct {
	// source table
	SourceName string
	SourceID   int

	// info table
	ID        int
	Priorite  int
	Agent     string
	Evenement string
	Created   time.Time
	Ouvrage   string
	Detail    string
	Status    string
	Echeance  string
	Entite    string

	// db conn
	DB *pgxpool.Pool

	// app self log
	ErrorLog *log.Logger
	InfoLog  *log.Logger

	// needed to create starting point date
	ZeroTime time.Time
}

func (data *Export) Export_DB_csv(conn *pgxpool.Conn) {
	ctx := context.Background()
	query := `
SELECT i.id AS "Info ID",
       s.name AS "Poste Source",
       i.evenement AS "Evènement",
       i.created AS "Date de détection",
       i.ouvrage AS "Ouvrage",
       i.detail AS "Détail",
       i.priorite AS "Priorité",
       i.status AS "Etat",
       i.echeance AS "Échéance",
       i.entite AS "Entité"
  FROM info AS i
    LEFT JOIN source AS s
    ON i.source_id = s.id
`

	// path := "./csvFiles/test_export.csv"
	path := "./csvFiles/test/test_export.csv"
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		data.ErrorLog.Println("Cannot create or open file")
		return // err file open
	}

	defer file.Close()

	rows, err := conn.Query(ctx, query)
	if err != nil {
		data.ErrorLog.Println("Couldn't fetch from DB")
		return // Add err return
	}
	defer rows.Close()

	for rows.Next() {
		line := &Export{}

		var evenement, ouvrage, detail, priorite, status, echeance, entite *string

		err = rows.Scan(&line.ID, &line.SourceName, &evenement, &line.Created, &ouvrage,
			&detail, &priorite, &status, &echeance, &entite)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return // error here pgx
			}
			data.ErrorLog.Println("Couldn't copy to var")
			data.ErrorLog.Println(err)
			return // Add err return
		}

		// line.Created, err = time.Parse("02/01/2006", *created)
		// if err != nil {
		// 	data.ErrorLog.Printf("Format de date invalide: %v", created)
		// 	data.ErrorLog.Println(err)
		// }

		if evenement != nil {
			line.Evenement = *evenement
		}

		if ouvrage != nil {
			line.Ouvrage = *ouvrage
		}

		if detail != nil {
			line.Detail = *detail
		}

		if priorite != nil {
			line.Priorite, _ = strconv.Atoi(*priorite)
		}

		if status != nil {
			line.Status = validator.ToCapital(*status)
		}

		if echeance != nil {
			line.Echeance = *echeance
		}

		if entite != nil {
			line.Entite = *entite
		}

		s := fmt.Sprintf("\"%v\",\"%v\",\"%v\",\"%v\",\"%v\",\"%v\",\"%v\",\"%v\",\"%v\",\"%v\"\n",
			line.ID, line.SourceName, line.Evenement, line.Created.Format("02/01/2006"), line.Ouvrage,
			line.Detail, line.Priorite, line.Status, line.Echeance, line.Entite)

		_, err := io.WriteString(file, s)
		if err != nil {
			data.ErrorLog.Println("io couldn't write to file")
			return // add io err
		}

	}

	data.decode_from_UTF8(path)
}

func (data *Export) decode_from_UTF8(s string) {
	file, err := os.ReadFile(s)
	if err != nil {
		log.Printf("File does not exists: %v", s)
		log.Println(err)
		return
	}

	tr, err := charmap.Windows1252.NewEncoder().Bytes(file)
	if err != nil {
		log.Printf("Bad encoded file: %v", file)
		log.Println(err)
		return
	}

	new_file := "./csvFiles/test/Actions_exportés.csv"
	err = os.WriteFile(new_file, tr, 0o666)
	if err != nil {
		log.Println("Cannot write to file")
		return
	}
}
