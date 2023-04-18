package database

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// TODO:
// Put everything back to normal, test it.
//

// On commence par vérifier si le fichier fini par .csv
// si vrai, alors on démarre EncodingCSV()
//
// Exemple:

// file.csv
// file.csv: text/csv; charset=iso-8859-1
//                            ^
//
// D'abord on lance la commande avec Output() afin de choper le string
// puis on la sépare avec "="
//
// On obtient:
// str[0] = file.csv: text/csv; charset
// str[1] = iso-8859-1
//
// Copie str[1] en majuscule dans une variable tmp
// Comme on ne sait pas quelle encodage le fichier peut avoir
// on le vérifie
// Si ce n'est pas en UTF-8, on change

// 2 structs sont créées afin de séparer chaque DB
// pour une meilleure lisibilité
type CSVInfo struct {
	ID       int
	Priority int
	SourceID int
	Agent    string
	Event    string
	Created  string // Cast to date with PSQL
	Material string
	Pilot    string
	Detail   string
	Target   string
	DayDone  string
	Estimate string
	Oups     string
	Brips    string
	Ameps    string
	Status   string
	DB       *pgxpool.Pool
	ErrorLog *log.Logger
	InfoLog  *log.Logger

	// srcID    int
	// srcName  string
}

type CSVSource struct {
	ID int
	Name string
	Created time.Time
	DB *pgxpool.Pool
	Errorlog *log.Logger
	InfoLog *log.Logger
}

func (data *CSVInfo) VerifyCSV(s string) {

	file := strings.Split(s, ".")
	length := len(file)

	if file[length-1] != "csv" {
		data.ErrorLog.Println("Wrong type of file")
	} else {
		data.encodingCSV(s)
	}
}

func (data *CSVInfo) encodingCSV(s string) {
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
	} else {
		data.dataCSV(s)
	}
}

func (data *CSVInfo) dataCSV(s string) {
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
	source, err := data.SourceNumber(lines[0][0])
	if err != nil {
		data.ErrorLog.Println(err)
	}

	// i = 0 -> Nom du poste
	// i = 1 -> Nom des colonnes
	for i := 2; i < len(lines); i++ {
		line := lines[i]
		j := 0

		data.Agent = line[j]
		data.Event = line[j+1]
		data.Created = line[j+2]
		data.Material = line[j+3]
		data.Pilot = line[j+4]
		data.Detail = line[j+5]
		data.Target = line[j+6]
		data.DayDone = line[j+7]
		data.Priority = 1
		data.Estimate = line[j+9]
		data.Oups = line[j+10]
		data.Brips = line[j+11]
		data.Ameps = line[j+12]
		data.SourceID = source
		if data.Status == "" && data.DayDone == "" &&
			data.Target == ""{
			data.Status = "en attente"
		} else if data.Target != "" && data.DayDone == "" {
			data.Status = "affecté"
		} else if data.Target != "" && data.DayDone != "" {
			data.Status = "résolu"
		}

		data.insertDB()
	}
}


// add source_id manualy for testing
func (data *CSVInfo) insertDB() {
	ctx := context.Background()
	query := `
INSERT INTO infos
  (source_id, agent, event, material, pilote, detail, target, day_done,
    priority, estimate, oups, brips, ameps, created, status)
      VALUES
	($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13,
	  (to_date($14, 'DD/MM/YYYY')), $15)
`
	_, err := data.DB.Exec(ctx, query, data.SourceID, data.Agent,
		data.Event, data.Material, data.Pilot, data.Detail,
		data.Target, data.DayDone, data.Priority, data.Estimate,
		data.Oups, data.Brips, data.Ameps, data.Created,
		data.Status)
	if err != nil {
		data.ErrorLog.Println(err)
	} else {
		data.InfoLog.Println("data sent")
	}
}

func (csv *CSVInfo) SourceNumber(s string) (int, error) {
	ctx := context.Background()
	query := `
SELECT id
  FROM sources
    WHERE name = $1
`

	var id int
	err := csv.DB.QueryRow(ctx, query, s).Scan(&id)
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
