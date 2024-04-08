package database

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CSVData struct {
	ID        int
	Priorite  int
	SourceID  int
	Agent     string
	Evenement string
	Created   string // Cast to date with PSQL
	Ouvrage   string
	Detail    string
	Target    string
	Devis     string
	Ameps     string
	Status    string
	DB        *pgxpool.Pool
	ErrorLog  *log.Logger
	InfoLog   *log.Logger

	SourceName string
}

func (data *CSVData) VerifyCSV(s string, conn *pgxpool.Conn) {
	file := strings.Split(s, ".")
	length := len(file)

	if file[length-1] != "csv" {
		data.ErrorLog.Println("Wrong type of file")
	} else {
		data.encodingCSV(s, conn)
	}
}

func (data *CSVData) encodingCSV(s string, conn *pgxpool.Conn) {
	cmd, err := exec.Command("file", "-i", s).Output()
	if err != nil {
		data.ErrorLog.Println(err)
	}

	strSplit := []string{}
	tmp := strings.Split(string(cmd), "=")
	strSplit = append(strSplit, tmp...)

	tmp2 := strings.ToUpper(strSplit[1])

	// Vérif si encodage est en UTF-8
	// si faux, on lance la commande de changement
	if tmp2 != "UTF-8\n" {
		cmd := exec.Command("iconv", "-f", tmp2,
			"-t", "UTF-8", s, "-o", s)
		iconvErr := cmd.Run()
		data.ErrorLog.Println(iconvErr)
	}

	data.dataCSV(s, conn)
}

func (data *CSVData) dataCSV(s string, conn *pgxpool.Conn) {
	file, err := os.Open(s)
	if err != nil {
		data.ErrorLog.Println(err)
	}
	defer file.Close()

	lines, err := csv.NewReader(file).ReadAll()
	if err != nil {
		data.ErrorLog.Println(err)
	}

	// Le fichier attendue a en première ligne et colonne
	// le nom du poste source.
	// Ce positionnement ne doit pas être changé!
	source, err := data.SourceNumber(lines[0][0], conn)
	if err != nil {
		data.ErrorLog.Println(err)
	}

	// i = 0 -> Nom du poste
	// i = 1 -> Nom des colonnes
	for i, j := 2, 0; i < len(lines); i++ {
		line := lines[i]

		data.Agent = line[j]
		data.Evenement = line[j+1]
		data.Created = line[j+2]
		data.Ouvrage = line[j+3]
		data.Detail = line[j+4]
		data.Priorite, _ = strconv.Atoi(line[j+5])
		data.Devis = line[j+6]
		data.SourceID = source
		data.Status = "en attente"

		data.insertDB(conn)
	}
}

func (data *CSVData) insertDB(conn *pgxpool.Conn) {
	ctx := context.Background()

	query := `
INSERT INTO info
  (source_id, agent, evenement, ouvrage, 
    detail, priorite, status, created)
      VALUES
	($1, $2, $3, $4, $5, $6, $7,
	  (to_date($8, 'DD/MM/YYYY')))
`
	_, err := conn.Exec(ctx, query, data.SourceID, data.Agent,
		data.Evenement, data.Ouvrage, data.Detail,
		data.Priorite, data.Status, data.Created)
	if err != nil {
		data.ErrorLog.Println(err)
	} else {
		data.InfoLog.Println("data sent")
	}
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

	fmt.Printf("@ sourceNumber: id > %v \n\n", id)

	return id, nil
}
