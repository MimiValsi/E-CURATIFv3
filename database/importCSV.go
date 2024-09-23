package database

import (
	"context"
	"encoding/csv"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/text/encoding/charmap"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type CSVData struct {
	ID        int
	Priorite  int
	SourceID  int
	Agent     string
	Evenement string
	Created   time.Time
	Ouvrage   string
	Detail    string
	Status    string
	Echeance  string
	Entite    string

	DB         *pgxpool.Pool
	ErrorLog   *log.Logger
	InfoLog    *log.Logger
	ZeroTime   time.Time
	SourceName string
}

func (data *CSVData) VerifyCSV(s string, conn *pgxpool.Conn) {
	file := strings.Split(s, ".")
	length := len(file)

	if file[length-1] != "csv" {
		data.ErrorLog.Println("Wrong type of file")
	} else {
		data.encoding_to_UTF8(s, conn)
	}
}

func (data *CSVData) encoding_to_UTF8(s string, conn *pgxpool.Conn) {
	file, err := os.ReadFile(s)
	if err != nil {
		log.Printf("Fichier n'existe pas: %v", s)
		log.Println(err)
		return
	}

	tr, err := charmap.Windows1252.NewDecoder().Bytes(file)
	if err != nil {
		log.Printf("Peut pas décoder fichier", file)
		log.Println(err)
		return
	}

	new_file := "new_utf8.csv"
	err = os.WriteFile(new_file, tr, 0666)
	if err != nil {
		log.Println("Peux pas écrire vers nouveau fichier")
		return
	}

	// log.Printf("new_file: %v", new_file)
	data.sendData(new_file, conn)
}

func (data *CSVData) sendData(s string, conn *pgxpool.Conn) {
	file, err := os.Open(s)
	if err != nil {
		data.ErrorLog.Println(err)
	}
	defer file.Close()

	lines, err := csv.NewReader(file).ReadAll()
	if err != nil {
		data.ErrorLog.Println(err)
	}

	data.ZeroTime = time.Date(0o001, time.January,
		1, 0, 0, 0, 0, time.UTC)

	log.Println("Transfert des données en cours")

	dateTmp := ""
	for i, j := 1, 0; i < len(lines); i++ {
		line := lines[i]

		data.SourceID, _ = data.SourceNumber(line[j], conn)
		data.Evenement = line[j+1]
		dateTmp = line[j+2]

		if dateTmp == "" {
			data.Created = time.Now().UTC()
		} else {
			data.Created, err = time.Parse("02/01/2006", dateTmp)
			if err != nil {
				log.Printf("Format de date invalide: %v", data.Created)
				log.Println(err)
				return
			}
		}

		data.Ouvrage = line[j+3]
		data.Detail = line[j+4]
		data.Priorite, _ = strconv.Atoi(line[j+5])
		data.Status = line[j+6]
		data.Echeance = line[j+7]
		data.Entite = line[j+8]

		data.insertDB(conn)
	}
}

func (data *CSVData) insertDB(conn *pgxpool.Conn) {
	ctx := context.Background()

	query := `
INSERT INTO info
  (source_id, evenement, ouvrage,
    detail, priorite, status, created, 
      echeance, entite)
        VALUES
	  ($1, $2, $3, $4, $5, $6, $7, $8 ,$9)
`
	// (to_date($8, 'DD/MM/YYYY'))
	_, err := conn.Exec(ctx, query, data.SourceID, data.Evenement,
		data.Ouvrage, data.Detail, data.Priorite, data.Status, data.Created,
		data.Echeance, data.Entite)
	if err != nil {
		data.ErrorLog.Println(err)
	}
	// } else {
	// 	data.InfoLog.Println("data sent")
	// }
}

func (data *CSVData) SourceNumber(s string, conn *pgxpool.Conn) (int, error) {
	ctx := context.Background()
	query := `
SELECT id
  FROM source
    WHERE name = $1
`

	var id int
	err := conn.QueryRow(ctx, query, s).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return -1, ErrNoRecord
		} else {
			return -1, err
		}
	}

	// fmt.Printf("%v sourceNumber: id > %v \n\n", s, id)

	return id, nil
}
