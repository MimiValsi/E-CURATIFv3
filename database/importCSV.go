package database

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	_"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// TODO:
// Put everything back to normal, test it.
//

// We start by verifying if the input file is a .csv
// if true than proceed to EncodingCSV()
//
// Execute command to fetch encoding type
//
// Exemple:

// file.csv
// file.csv: text/csv; charset=iso-8859-1
//                            ^
//
// 1st run the command with Output() to fetch the string
// 2nd split @ "="
//
// We'll get:
// str[0] = file.csv: text/csv; charset
// str[1] = iso-8859-1
//
// Copy str[1] to uppercase into a tmp variable
// As we don't know which encoding type a file might have
// for every file scanned, we verify it's encoding
// if it's not UTF-8 then run cmd to change to it.

// 2 structs are created to separate each PSQL tables
// and for better readability
type CSVInfo struct {
	ID       int
	Agent    string
	Event    string
	Created  string // Cast to date with PSQL
	Material string
	Pilot    string
	Detail   string
	Target   string
	DayDone  string
	Priority int
	Estimate string
	Oups     string
	Brips    string
	Ameps    string
	Status   string
	SourceID int
	DB       *pgxpool.Pool
	ErrorLog *log.Logger
	InfoLog  *log.Logger

	srcID    int
	srcName  string

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

	// // Check if encoding type is UTF-8
	// // if false then run encoding cmd
	// // \n at the end not good...
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

	//data.sourceNumber(lines[0][0])
	// fmt.Println(lines[0][0])

	for i := 2; i < 3; i++ {
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
		// fmt.Printf("line[j+8] >%v\n\n", line[j+8])
		// fmt.Printf("line[j+8] > %T\n\n", line[j+8])
		// data.Priority, err = strconv.Atoi(line[j+8])
		// if err != nil {
		//	data.errorlog.Println(err)
		// }
		data.Priority = 1
		data.Estimate = line[j+9]
		data.Oups = line[j+10]
		data.Brips = line[j+11]
		data.Ameps = line[j+12]
		data.SourceID = 20
		// data.Agent = "Bob"
		// data.Event = "Inc Bat"
		// data.Created = "20/12/2017"
		// data.Material = "TR 611"
		// data.Pilot = "AMEPS CE"
		// data.Detail = "HS"
		// data.Target = "20/01/2018"
		// data.DayDone = ""
		// data.Priority = 1
		// data.Estimate = "10EUR"
		// data.Oups = "X"
		// data.Brips = "X"
		// data.Ameps = "X"
		// data.SourceID, err = data.sourceNumber("Novion")
		// if err != nil {
		//	fmt.Println(err)
		// }
		if data.Status == "" && data.DayDone == "" &&
			data.Target == ""{
			data.Status = "en attente"
		} else if data.Target != "" && data.DayDone == "" {
			data.Status = "affecté"
		} else if data.Target != "" && data.DayDone != "" {
			data.Status = "résolu"
		}


		// fmt.Println(data)
		data.insertDB()
		// fmt.Printf("%v\n", data.Agent)
		// fmt.Printf("%v\n", data.Event)
		// fmt.Printf("%v\n", data.Created)
		// fmt.Printf("%v\n", data.Material)
		// fmt.Printf("%v\n", data.Pilot)
		// fmt.Printf("%v\n", data.Detail)
		// fmt.Printf("%v\n", data.Target)
		// fmt.Printf("%v\n", data.DayDone)
		// fmt.Printf("%v\n", line[j+8])
		// fmt.Printf("%v\n", data.Estimate)
		// fmt.Printf("%v\n", data.Oups)
		// fmt.Printf("%v\n", data.Brips)
		// fmt.Printf("%v\n", data.Ameps)
		// fmt.Printf("%v\n", data.SourceID)
	}
}

// add source_id manualy for testing
func (data *CSVInfo) insertDB() {
	ctx := context.Background()
	query := `
INSERT INTO infos
  (source_id, agent, event, material, pilot, detail, target, day_done,
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
		data.ErrorLog.Println("data sent")
	}
	// fmt.Print(data)
}

func (csv *CSVInfo) sourceNumber(s string) (int, error) {
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
