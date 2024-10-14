package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/justinas/alice"
)

// Chaque page commence avec chi.NewRouter()
func (app *application) routes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Fichiers statiques
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	r.Handle("GET /static/*", http.StripPrefix("/static", fileServer))

	chain := chi.Chain(app.sessionManager.LoadAndSave, noSurf, app.authenticate)
	protected := append(chain, app.requireAuthentication)

	// Home page
	r.Handle("GET /", chain.HandlerFunc(app.home))
	r.Handle("GET /jsonGraph", chain.HandlerFunc(app.jsonData))

	// Pages Source
	// Chaque place réservée doit être unique pour chaque router
	r.Handle("GET /source/view/{id}", chain.HandlerFunc(app.sourceView))
	r.Handle("GET /source/create", protected.HandlerFunc(app.sourceCreate))
	r.Handle("POST /source/create", protected.HandlerFunc(app.sourceCreatePost))
	r.Handle("POST /source/delete/{id}", protected.HandlerFunc(app.sourceDeletePost))
	r.Handle("GET /source/update/{id}", protected.HandlerFunc(app.sourceUpdate))
	r.Handle("POST /source/update/{id}", protected.HandlerFunc(app.sourceUpdatePost))

	// Pages Infos
	r.Handle("GET /source/{sid}/info/view/{id}", chain.HandlerFunc(app.infoView))
	r.Handle("GET /source/{id}/info/create", protected.HandlerFunc(app.infoCreate))
	r.Handle("POST /source/{id}/info/create", protected.HandlerFunc(app.infoCreatePost))
	r.Handle("POST /source/{sid}/info/delete/{id}", protected.HandlerFunc(app.infoDeletePost))
	r.Handle("GET /source/{sid}/info/update/{id}", protected.HandlerFunc(app.infoUpdate))
	r.Handle("POST /source/{sid}/info/update/{id}", protected.HandlerFunc(app.infoUpdatePost))

	// Pages user
	r.Handle("GET /user/signup", chain.HandlerFunc(app.userSignup))
	r.Handle("POST /user/signup", chain.HandlerFunc(app.userSignupPost))
	r.Handle("GET /user/login", chain.HandlerFunc(app.userLogin))
	r.Handle("POST /user/login", chain.HandlerFunc(app.userLoginPost))
	r.Handle("POST /user/logout", chain.HandlerFunc(app.userLogoutPost))
	r.Handle("GET /user/profile", chain.HandlerFunc(app.userProfile))

	// Import / Export CSV
	r.Handle("POST /importCSV", protected.HandlerFunc(app.importCSVPost))
	r.Handle("GET /exportCSV", protected.HandlerFunc(app.exportCSVPost))

	std := alice.New(app.recoverPanic, app.logRequest, commonHeaders)
	return std.Then(r)
}
