package database

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v4/pgxpool"
)

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

type CSVdata struct {
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
	SourceID int
	DB       *pgxpool.Pool
	errorlog *log.Logger
	infoLog  *log.Logger
}

func VerifyCSV(s string) {
	data := CSVdata{}
	file := strings.Split(s, ".")
	length := len(file)

	if file[length-1] != "csv" {
		fmt.Println("Wrong type of file")
	} else {
		data.encodingCSV(s)
	}
}

func (data *CSVdata) encodingCSV(s string) {
	cmd, err := exec.Command("file", "-i", s).Output()
	if err != nil {
		data.infoLog.Println(err)
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
		fmt.Println(iconvErr)
	} else {
		data.dataCSV(s)
	}
}

func (data *CSVdata) dataCSV(s string) {
	file, err := os.Open(s)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	lines, err := csv.NewReader(file).ReadAll()
	if err != nil {
		fmt.Println(err)
	}

	//data.sourceNumber(lines[0][0])
	// fmt.Println(lines[0][0])

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
		tmp := line[j+8]
		priority, _ := strconv.Atoi(tmp)
		data.Priority = priority
		data.Estimate = line[j+9]
		data.Oups = line[j+10]
		data.Brips = line[j+11]
		data.Ameps = line[j+12]

		// insertDB()
	}
}

func (data *CSVdata) insertDB() {
	ctx := context.Background()
	query := `
INSERT INTO infos
  (agent, event, material, pilot, detail, target, day_done,
    priority, estimate, oups, brips, ameps, created, source_id)
VALUES
  ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12,
    (to_date($13, dd/mm/yyyy)), $14)
	`
	err := data.DB.QueryRow(ctx, query, data.Agent, data.Event,
		data.Material, data.Pilot, data.Detail, data.Target,
		data.DayDone, data.Priority, data.Estimate,
		data.Oups, data.Brips, data.Ameps, data.Created,
		data.SourceID)
	if err != nil {
		data.errorlog.Println(err)
	}
}

// func (data *CSVdata) sourceNumber(s string) int {
//	switch s {
//	case ""
//	}
//	return 0
// }
