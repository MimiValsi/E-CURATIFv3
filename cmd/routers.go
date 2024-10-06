package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Chaque page commence avec chi.NewRouter()
func (app *application) routes() http.Handler {
	// r := http.NewServeMux()
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Fichiers statiques
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	r.Handle("GET /static/*", http.StripPrefix("/static", fileServer))

	chain := chi.Chain(app.sessionManager.LoadAndSave).HandlerFunc

	// Home page
	r.Handle("GET /", chain(app.home))
	r.Handle("GET /jsonGraph", chain(app.jsonData))

	// Pages Source
	// Chaque place réservée doit être unique pour chaque router
	r.Handle("GET /source/view/{id}", chain(app.sourceView))
	r.Handle("GET /source/create", chain(app.sourceCreate))
	r.Handle("POST /source/create", chain(app.sourceCreatePost))
	r.Handle("POST /source/delete/{id}", chain(app.sourceDeletePost))
	r.Handle("GET /source/update/{id}", chain(app.sourceUpdate))
	r.Handle("POST /source/update/{id}", chain(app.sourceUpdatePost))

	// Pages Infos
	r.Handle("GET /source/{sid}/info/view/{id}", chain(app.infoView))
	r.Handle("GET /source/{id}/info/create", chain(app.infoCreate))
	r.Handle("POST /source/{id}/info/create", chain(app.infoCreatePost))
	r.Handle("POST /source/{sid}/info/delete/{id}", chain(app.infoDeletePost))
	r.Handle("GET /source/{sid}/info/update/{id}", chain(app.infoUpdate))
	r.Handle("POST /source/{sid}/info/update/{id}", chain(app.infoUpdatePost))

	// Pages user
	r.Handle("GET /user/signup", chain(app.userSignup))
	r.Handle("POST /user/signup", chain(app.userSignupPost))
	r.Handle("GET /user/login", chain(app.userLogin))
	r.Handle("GET /user/login", chain(app.userLoginPost))
	r.Handle("GET /user/logout", chain(app.userLogoutPost))

	// Import / Export CSV
	r.Handle("POST /importCSV", chain(app.importCSVPost))
	r.Handle("POST /exportCSV", chain(app.exportCSVPost))

	return r
}
