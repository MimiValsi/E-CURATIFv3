package database

import (
	// "context"
	"log"
	// "os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	// "golang.org/x/text/encoding/charmap"
)

type Export struct {
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

// func (data *Export) Export_DB_csv(conn *pgxpool.Conn) {
// 	ctx := context.Background()
// 	query := `
// \copy (SELECT
// 	s.name AS "Poste Source",
// 	i.evenement AS "Evènement",
// 	i.created AS "Date de détection",
// 	i.ouvrage AS "Ouvrage",
// 	i.detail AS "Détail",
// 	i.priorite AS "Priorité",
// 	i.status AS "Etat",
// 	i.echeance AS "Échéance",
// 	i.entite AS "Entité"
// FROM info AS i
//   LEFT JOIN source AS s
//   ON i.source_id = s.id)
//   TO '~/Projects/E-CURATIFv3/csvFiles/export_actions.csv'
//   DELIMITER ',' CSV HEADER
// `
//
// }
//
// func (data *Export) decode_from_UTF8(s string) {
// 	file, err := os.ReadFile(s)
// 	if err != nil {
// 		log.Printf("File does not exists: %v", s)
// 		log.Println(err)
// 		return
// 	}
//
// 	tr, err := charmap.Windows1252.NewEncoder().Bytes(file)
// 	if err != nil {
// 		log.Printf("Bad encoded file: %v", file)
// 		log.Println(err)
// 		return
// 	}
//
// 	new_file := "./csvFiles/Actions_exportés.csv"
// 	err = os.WriteFile(new_file, tr, 0666)
// 	if err != nil {
// 		log.Println("Cannot write to file")
// 		return
// 	}
//
// 	// data.export(new_file, conn)
// }
//
// // func (data *Export) export()
