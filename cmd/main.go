package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"

	"E-CURATIFv3/database"

	"github.com/jackc/pgx/v4/pgxpool"
)

type application struct {
	sources       *database.Source
	infos         *database.Info
	templateCache map[string]*template.Template
	errorLog      *log.Logger
	infoLog       *log.Logger
}

const (
	addr    = ":3001"
	dataURL = "postgres://web:pass@localhost:5432/ecuratif"
)

func main() {
	infoLog := log.New(os.Stderr, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t",
		log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(dataURL)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	app := &application{
		sources:       &database.Source{DB: db},
		infos:         &database.Info{DB: db},
		templateCache: templateCache,
	}

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
