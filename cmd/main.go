package main

import (
	"context"
	"crypto/tls"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	"github.com/jackc/pgx/v5/pgxpool" // PostgreSQL driver

	"E-CURATIFv3/database" // Database regroupe toutes les fonctions pour communiquer avec PSQL
	"E-CURATIFv3/internal/models"
)

// afin de permettre la vérif les informations et communiquer avec PSQL
type application struct {
	sources *database.Source
	infos   *database.Info
	users   *models.User

	templateCache map[string]*template.Template

	errorLog *log.Logger
	infoLog  *log.Logger
	logger   *slog.Logger

	csvImport *database.Import
	csvExport *database.Export

	DB *pgxpool.Pool

	sessionManager *scs.SessionManager
}

// Const for dev
const (
	addr    = ":8080"
	dataURL = "postgres://postgres:pass@localhost:5432/ecuratif"
	// dataURL = "postgres://ameps:pass@localhost:5432/test"
)

func main() {
	// infoLog and errorLog are own made middleware
	// in complementary of go-chi middleware

	// Ldate = Local data & Ltime = Local time
	infoLog := log.New(os.Stderr, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t",
		log.Ldate|log.Ltime|log.Lshortfile)

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	var sessionManager *scs.SessionManager

	// Send an connection request to DB
	db, err := openDB(dataURL)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	sessionManager = scs.New()
	sessionManager.Store = pgxstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	// Func @ cmd/template.go
	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	// Application dependencies are centralized here
	// Gives the possibility to access the middleware logging everywhere
	app := &application{
		DB:      db,
		sources: &database.Source{},
		infos:   &database.Info{},
		users:   &models.User{},

		templateCache: templateCache,
		csvImport:     &database.Import{InfoLog: infoLog, ErrorLog: errorLog},
		csvExport:     &database.Export{InfoLog: infoLog, ErrorLog: errorLog},

		infoLog:  infoLog,
		errorLog: errorLog,
		logger:   logger,

		sessionManager: sessionManager,
	}

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	// Default parameters values to routes
	// See routers.go
	srv := &http.Server{
		Addr:         addr,
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		TLSConfig:    tlsConfig,
	}

	infoLog.Printf("Starting server on %s", addr)
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	errorLog.Fatal(err)
}

// Create a context background and create a new PSQL connection
func openDB(dataURL string) (*pgxpool.Pool, error) {
	ctx := context.Background()

	db, err := pgxpool.New(ctx, dataURL)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(ctx); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
