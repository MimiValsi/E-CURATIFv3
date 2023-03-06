package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// routers file
// each page starts with chi.NewRouter()
// See github.com/go-chi/chi for better infos
func (app *application) routes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// Home page
	r.Get("/", app.home)

	// Sources pages
	// Each placeholder must be unique for each route
	r.Get("/source/view/{id}", app.sourceView)
	r.Get("/source/create", app.sourceCreate)
	r.Post("/source/create", app.sourceCreatePost)
	r.Post("/source/delete/{id}", app.sourceDeletePost)
	r.Get("/source/update/{id}", app.sourceUpdate)
	r.Post("/source/update/{id}", app.sourceUpdatePost)

	// Infos page
	r.Get("/source/{sid}/info/view/{id}", app.infoView)
	r.Get("/source/{id}/info/create", app.infoCreate)
	r.Post("/source/{id}/info/create", app.infoCreatePost)
	r.Post("/source/{sid}/info/delete/{id}", app.infoDeletePost)
	r.Get("/source/{sid}/info/update/{id}", app.infoUpdate)
	r.Post("/source/{sid}/info/update/{id}", app.infoUpdatePost)

	r.Get("/info/upload", app.infoUpload)
	r.Post("/info/upload", app.infoUploadPost)

	// Statics files
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	r.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return r
}
