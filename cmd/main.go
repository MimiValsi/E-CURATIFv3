package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"

	// Database regroupe toutes les fonctions pour communiquer
	// avec PSQL
	"E-CURATIFv3/database"

	// PostgreSQL driver
	"github.com/jackc/pgx/v4/pgxpool"
)

// Main struct, every other struct must be "connected" to this one
// It enables to parse informations and communicate with PSQL
type application struct {
	sources       *database.Source
	infos         *database.Info
	templateCache map[string]*template.Template
	errorLog      *log.Logger
	infoLog       *log.Logger
	csvSource     *database.CSVSource
	csvInfo       *database.CSVInfo
}

// This 2 variables aren't meant to be change nor to production
const (
	addr    = ":3001"
	dataURL = "postgres://web:pass@localhost:5432/ecuratif"
)

func main() {
	// infoLog and errorLog may give a bit more information about
	// errors and/or others.
	// Ldate = Local data & Ltime = Local time
	infoLog := log.New(os.Stderr, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t",
		log.Ldate|log.Ltime|log.Lshortfile)

	// Run DB connexion func
	db, err := openDB(dataURL)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	// Func @ cmd/template.go
	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	// To allow communication with PSQL or others functions,
	// like infoLog / errorLog, they must be parse throw here.
	app := &application{
		sources:       &database.Source{DB: db},
		infos:         &database.Info{DB: db},
		templateCache: templateCache,
		csvInfo:       &database.CSVInfo{DB: db,
			ErrorLog: errorLog, InfoLog: infoLog},
		csvSource:     &database.CSVSource{DB: db,
			Errorlog: errorLog, InfoLog: infoLog},
		infoLog: infoLog,
		errorLog: errorLog,
	}

	// Default parameters values to routes
	// See routers.go
	srv := &http.Server{
		Addr:         addr,
		Handler:      app.routes(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Function test
	// app.csvSource.SourceNumber("Nanterre")
	// app.csvInfo.VerifyCSV("test.csv")
	// app.csvInfo.SourceNumber("Ampère")
	// app.csvInfo.DataCSV("test.csv")

	infoLog.Printf("Starting server on %s", addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)

}

// Function that allows to connect to PSQL via pgxpool
func openDB(dataURL string) (*pgxpool.Pool, error) {
	ctx := context.Background()

	db, err := pgxpool.Connect(ctx, dataURL)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(ctx); err != nil {
		return nil, err
	}

	return db, nil
}
