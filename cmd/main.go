package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"

	// Every func that send or request someting from PSQL are here
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
	// info and error need to be refactored
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

	// Every struct that communicates to PSQL are parse throw here
	app := &application{
		sources:       &database.Source{DB: db},
		infos:         &database.Info{DB: db},
		templateCache: templateCache,
		csvInfo:       &database.CSVInfo{DB: db},
		csvSource:     &database.CSVSource{DB: db},
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
