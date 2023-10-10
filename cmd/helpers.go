package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
)

// Les status web sont gérés ici

// Le serverError écrit les message d'erreur
// puis envoi 500 Internal Server Error à l'utilisateur
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())

	app.errorLog.Print(trace)
}

// clientError envoi un status spécific et la déscription
// correspondante à l'utilisateur
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// notFound fait la même que clientError mais envoi 404 Not Found
// à l'utilisateur
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

// Alloue de la mémoire pour qu'un template puisse être rendue
// Vérifie si le template dérisé existe avant d´être envoyé
// au http.ResponseWriter
func (app *application) render(w http.ResponseWriter, status int, page string, 
        data *templateData) {

	// Récupère le template approprié du cache
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist",
			page)
		app.serverError(w, err)
		return
	}

	buf := new(bytes.Buffer)

	// Execute les templates et envoie au response body
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, err)
	}

	w.WriteHeader(status)

	buf.WriteTo(w)
}

// newTemplateData retourne un pointeur vers templateData
// non initializé et est utilisé par toutes les fonctions dans
// Handler. Permer une meilleur lisibilité du code
func (app *application) newTemplateData(r *http.Request) *templateData {
	return &templateData{}
}
