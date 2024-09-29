package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Chaque page commence avec chi.NewRouter()
func (app *application) routes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// Home page
	r.Get("/", app.home)
	r.Get("/jsonGraph", app.jsonData)
	r.Get("/prioData", app.priorityData)

	// Pages Source
	// Chaque place réservée doit être unique pour chaque router
	r.Get("/source/view/{id}", app.sourceView)
	r.Get("/source/create", app.sourceCreate)
	r.Post("/source/create", app.sourceCreatePost)
	r.Post("/source/delete/{id}", app.sourceDeletePost)
	r.Get("/source/update/{id}", app.sourceUpdate)
	r.Post("/source/update/{id}", app.sourceUpdatePost)

	// Pages Infos
	r.Get("/source/{sid}/info/view/{id}", app.infoView)
	r.Get("/source/{id}/info/create", app.infoCreate)
	r.Post("/source/{id}/info/create", app.infoCreatePost)
	r.Post("/source/{sid}/info/delete/{id}", app.infoDeletePost)
	r.Get("/source/{sid}/info/update/{id}", app.infoUpdate)
	r.Post("/source/{sid}/info/update/{id}", app.infoUpdatePost)

	// En cours de création
	r.Post("/importCSV", app.importCSVPost)
	r.Post("/exportCSV", app.exportCSVPost)

	// Fichiers statiques
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	r.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return r
}
