package database

import (
	"context"
	"encoding/csv"
	"errors"
	// "fmt"
	"log"
	"os"

	// "github.com/jackc/pgx/v4"
	// "github.com/jackc/pgx/v4/pgxpool"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CSVData struct {
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
        Entity   string
        Comment  string

        DB       *pgxpool.Pool

        ErrorLog *log.Logger
        InfoLog  *log.Logger

        srcID    int
        srcName  string

}

func (data *CSVData) Import(s string) {
        file, err := os.Open(s)
        if err != nil {
                data.ErrorLog.Println(err)
        }
        defer file.Close()

        // lines, err := csv.NewReader(file).ReadAll()
        l := csv.NewReader(file)

        l.Comma = ';'

        lines, err := l.ReadAll()
        if err != nil {
                data.ErrorLog.Println(err)
        }

        for i := 1; i < len(lines); i++ {
                line := lines[i]
                j := 0

                data.fetchSourceID(toUTF8([]byte(line[j])))
                data.Material = toUTF8([]byte(line[j+1]))
                data.Detail = toUTF8([]byte(line[j+2]))
                data.Created = toUTF8([]byte(line[j+3]))
                data.Agent = toUTF8([]byte(line[j+4]))
                data.Entity = toUTF8([]byte(line[j+5]))
                data.Status = toUTF8([]byte(line[j+6]))
                data.DayDone = toUTF8([]byte(line[j+8]))
                data.Comment = toUTF8([]byte(line[j+9]))
                // fmt.Println(data)
                data.insert()
                // break
        }

}

func (data *CSVData) insert() {
        ctx := context.Background()
        query := `
INSERT INTO info
  (source_id, agent, material, detail, created, entity, status, day_done, comment, priority, event)
       VALUES
        ($1, $2, $3, $4, (to_date($5, 'DD-MM-YYYY')), $6, $7, $8, $9, $10, $11)
`
        args := []any{data.srcID, data.Agent, data.Material, 
                data.Detail, data.Created, data.Entity, 
                data.Status, data.DayDone, data.Comment, 3, "test"}

        _, err := data.DB.Exec(ctx, query, args...)
        if err != nil {
                data.ErrorLog.Println(err)
        }
}

func (data *CSVData) fetchSourceID(s string) (int, error) {
        ctx := context.Background()
        query := `
SELECT id 
  FROM source 
 WHERE name = $1
`

        err := data.DB.QueryRow(ctx, query, s).Scan(&data.srcID)
        if err != nil {
                if errors.Is(err, pgx.ErrNoRows) {
                        return -1, ErrNoRecord
                } else {
                        return -1, err
                }
        }
        
        return data.srcID, nil
}

func toUTF8(iso []byte) string {
        buf := make([]rune, len(iso))
        for i, b := range iso {
                buf[i] = rune(b)
        }

        return string(buf)
}
