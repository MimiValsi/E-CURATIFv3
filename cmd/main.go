package main

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool" // PostgreSQL driver

	"E-CURATIFv3/database" // Database regroupe toutes les fonctions pour communiquer avec PSQL
)

// afin de permettre la vérif les informations et communiquer avec PSQL
type application struct {
	sources *database.Source
	infos   *database.Info

	templateCache map[string]*template.Template

	errorLog *log.Logger
	infoLog  *log.Logger

	// csvSource *database.CSVSource
	csvImport *database.Import
	csvExport *database.Export

	DB *pgxpool.Pool

	// jSource *database.JsonSource
}

// Ces 2 variables ne sont pas sensé être ni modifiés ni pour la prod
const (
	addr    = ":3001"
	dataURL = "postgres://ameps:pass@localhost:5432/ecuratif"
	// dataURL = "postgres://ameps:pass@localhost:5432/test"
)

func main() {
	// infoLog et errorLog permettent d'avoir un peu plus d'info
	// de ce qui se passe en cas d'erreur

	// Ldate = Local data & Ltime = Local time
	infoLog := log.New(os.Stderr, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t",
		log.Ldate|log.Ltime|log.Lshortfile)

	// Execute la fonction de connection de la BD
	db, err := openDB(dataURL)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	// Fontion @ cmd/template.go
	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	// Permet la comm avec PSQL et autres fonctions
	// Tout passe par ici
	app := &application{
		DB:      db,
		sources: &database.Source{},
		infos:   &database.Info{},

		templateCache: templateCache,
		csvImport:     &database.Import{InfoLog: infoLog, ErrorLog: errorLog},
		csvExport:     &database.Export{InfoLog: infoLog, ErrorLog: errorLog},

		infoLog:  infoLog,
		errorLog: errorLog,

		// jSource: &database.JsonSource{DB: db},
	}

	// Default parameters values to routes
	// See routers.go
	srv := &http.Server{
		Addr:         addr,
		Handler:      app.routes(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("Starting server on %s", addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

// Fonction qui permet la connection avec PSQL via pgx.pgxpool
func openDB(dataURL string) (*pgxpool.Pool, error) {
	ctx := context.Background()

	db, err := pgxpool.New(ctx, dataURL)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(ctx); err != nil {
		return nil, err
	}

	return db, nil
}
