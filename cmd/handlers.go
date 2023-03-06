package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"E-CURATIFv3/database"
	"E-CURATIFv3/internal/validator"

	// package pour les routers
	"github.com/go-chi/chi/v5"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// MenuSource func @ database/sources.go
	sources, err := app.sources.MenuSource()
	if err != nil {
		app.serverError(w, err)
		return
	}

	// newTemplateData @ cmd/templates.go
	data := app.newTemplateData(r)
	data.Sources = sources

	app.render(w, http.StatusOK, "home.tmpl.html", data)

}

//
// Sources Handlers
//

type sourceCreateForm struct {
	Name string
	validator.Validator
}

// Ici on génère la page de visu d'un poste source avec ces curatifs
func (app *application) sourceView(w http.ResponseWriter, r *http.Request) {

	key := chi.URLParam(r, "id")
	id, err := strconv.Atoi(key)
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	// Fait appel à la fonctin dans database/sources.go
	// On récupère le "id" du source dans le URL
	// créé auparavant
	source, err := app.sources.SourceGet(id)
	if err != nil {
		if errors.Is(err, database.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
	}

	// Fait appel à la fonction dans database/infos.go
	// en prenant en compte le id du source
	info, err := app.infos.InfoList(id)
	if err != nil {
		if errors.Is(err, database.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
	}

	// Allocation de mémoire
	data := app.newTemplateData(r)
	data.Infos = info
	data.Source = source

	// Génération de la page web
	app.render(w, http.StatusOK, "sourceView.tmpl.html", data)

}

// Génération de la page de création de source
// Celle-ci est faite en deux parties:
// La première est une GET, une fois le nom du source choisi
// on passe à la fonction sourceCreatePost
func (app *application) sourceCreate(w http.ResponseWriter, r *http.Request) {

	data := app.newTemplateData(r)
	data.Form = sourceCreateForm{}

	app.render(w, http.StatusOK, "sourceCreate.tmpl.html", data)
}

func (app *application) sourceCreatePost(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
	}

	form := sourceCreateForm{
		Name: r.PostForm.Get("name"),
	}

	emptyField := "Ce champ ne doit pas être vide"

	form.CheckField(validator.NotBlank(form.Name),
		"name", emptyField)
	// form.CheckField(validator.MaxChars(form.Name, 20),
	// "name", "Nom de Poste Source trop grand")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity,
			"sourceCreate.tmpl.html", data)
		return
	}

	id, err := app.sources.SourceInsert(form.Name)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/source/view/%d", id),
		http.StatusSeeOther)
}

func (app *application) sourceDeletePost(w http.ResponseWriter, r *http.Request) {

	key := chi.URLParam(r, "id")

	id, err := strconv.Atoi(key)
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	err = app.sources.SourceDelete(id)
	if err != nil {
		if errors.Is(err, database.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func (app *application) sourceUpdate(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "id")
	id, err := strconv.Atoi(key)
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	source, err := app.sources.SourceGet(id)
	if err != nil {
		if errors.Is(err, database.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Source = source

	app.render(w, http.StatusOK, "sourceUpdate.tmpl.html", data)
}

func (app *application) sourceUpdatePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
	}

	key := chi.URLParam(r, "id")
	id, err := strconv.Atoi(key)
	if err != nil {
		app.notFound(w)
		return
	}

	form := sourceCreateForm{
		Name: r.PostForm.Get("name"),
	}

	emptyField := "Ce champ ne doit pas être vide"

	form.CheckField(validator.NotBlank(form.Name),
		"name", emptyField)

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity,
			"sourceUpdate.tmpl.html", data)
		return
	}

	app.sources.Name = form.Name

	err = app.sources.SourceUpdate(id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/source/view/%d", id),
		http.StatusSeeOther)

}

//
// Infos Handlers
//

type infoCreateForm struct {
	Agent    string
	Material string
	Priority string
	Target   string
	Detail   string
	Created  string
	Updated  string
	Status   string
	Event    string
	Rte      string
	Estimate string
	Brips    string
	Ais      string
	Oups     string
	Ameps    string
	Doneby   string

	validator.Validator
}

// GET form that fetch sourceID which is sent to URL and retrieved
// by infoCreatePost func, so a info can be created.
// Infos table id has a FK to Sources table id.
func (app *application) infoCreate(w http.ResponseWriter, r *http.Request) {

	key := chi.URLParam(r, "id")

	id, err := strconv.Atoi(key)
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	source, err := app.sources.SourceGet(id)
	if err != nil {
		if errors.Is(err, database.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Form = infoCreateForm{}
	data.Source = source

	app.render(w, http.StatusOK, "infoCreate.tmpl.html", data)
}

// POST form that fetch, control and sends data to psql server
func (app *application) infoCreatePost(w http.ResponseWriter, r *http.Request) {

	// Fetch input
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
	}

	// Fetch SourceID
	key := chi.URLParam(r, "id")

	sID, err := strconv.Atoi(key)
	if err != nil || sID < 1 {
		app.notFound(w)
		return
	}

	form := infoCreateForm{
		Agent:    r.PostForm.Get("agent"),
		Material: r.PostForm.Get("material"),
		Detail:   r.PostForm.Get("detail"),
		Event:    r.PostForm.Get("event"),
		Priority: r.PostForm.Get("priority"),
		Oups:     r.PostForm.Get("oups"),
		Ameps:    r.PostForm.Get("ameps"),
		Brips:    r.PostForm.Get("brips"),
		Rte:      r.PostForm.Get("rte"),
		Ais:      r.PostForm.Get("ais"),
		Estimate: r.PostForm.Get("estimate"),
		Target:   r.PostForm.Get("target"),
		Status:   r.PostForm.Get("status"),
		Doneby:   r.PostForm.Get("doneby"),
	}

	emptyField := "Ce champ ne doit pas être vide"

	form.CheckField(validator.NotBlank(form.Agent),
		"agent", emptyField)
	form.CheckField(validator.NotBlank(form.Material),
		"material", emptyField)
	form.CheckField(validator.NotBlank(form.Detail),
		"detail", emptyField)
	form.CheckField(validator.NotBlank(form.Event),
		"event", emptyField)
	form.CheckField(validator.NotBlank(form.Priority),
		"priority", emptyField)
	form.CheckField(validator.NotBlank(form.Status),
		"status", emptyField)

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity,
			"infoCreate.tmpl.html", data)
		return
	}

	// By copying form.xxx into app.infos.xxx, it will
	// send the data to Infos struct @ database/infos.go
	app.infos.Agent = form.Agent
	app.infos.Material = form.Material
	app.infos.Detail = form.Detail
	app.infos.Event = form.Event
	app.infos.Oups = form.Oups
	app.infos.Ameps = form.Ameps
	app.infos.Brips = form.Brips
	app.infos.Rte = form.Rte
	app.infos.Ais = form.Ais
	app.infos.Estimate = form.Estimate
	app.infos.Target = form.Target
	app.infos.Status = form.Status
	app.infos.Doneby = form.Doneby
	app.infos.Priority, err = strconv.Atoi(form.Priority)
	if err != nil {
		app.notFound(w)
		return
	}

	id, err := app.infos.Insert(sID)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/source/%d/info/view/%d",
		sID, id), http.StatusSeeOther)
}

// Func that retrive Info table data and send it to be displayed
func (app *application) infoView(w http.ResponseWriter, r *http.Request) {
	iKey := chi.URLParam(r, "id")

	id, err := strconv.Atoi(iKey)
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	info, err := app.infos.InfoGet(id)
	if err != nil {
		if errors.Is(err, database.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Info = info

	app.render(w, http.StatusOK, "infoView.tmpl.html", data)
}

// POST form to delete selected info
func (app *application) infoDeletePost(w http.ResponseWriter, r *http.Request) {

	sKey := chi.URLParam(r, "sid")
	iKey := chi.URLParam(r, "id")

	id, err := strconv.Atoi(iKey)
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	sID, err := strconv.Atoi(sKey)
	if err != nil || sID < 1 {
		app.notFound(w)
		return
	}

	err = app.infos.InfoDelete(id)
	if err != nil {
		if errors.Is(err, database.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/source/view/%d", sID),
		http.StatusSeeOther)

}

func (app *application) infoUpdate(w http.ResponseWriter, r *http.Request) {

	key := chi.URLParam(r, "id")
	id, err := strconv.Atoi(key)
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	// id ~> Info id
	info, err := app.infos.InfoGet(id)
	if err != nil {
		if errors.Is(err, database.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Info = info

	app.render(w, http.StatusOK, "infoUpdate.tmpl.html", data)
}



func (app *application) infoUpdatePost(w http.ResponseWriter, r *http.Request) {
	// Fetch input
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
	}

	// Fetch SourceID
	sKey := chi.URLParam(r, "sid")
	sID, err := strconv.Atoi(sKey)
	if err != nil || sID < 1 {
		app.notFound(w)
		return
	}

	// Fetch InfoID
	iKey := chi.URLParam(r, "id")
	iID, err := strconv.Atoi(iKey)
	if err != nil || iID < 1 {
		app.notFound(w)
		return
	}

	form := infoCreateForm{
		Agent:    r.PostForm.Get("agent"),
		Material: r.PostForm.Get("material"),
		Detail:   r.PostForm.Get("detail"),
		Event:    r.PostForm.Get("event"),
		Priority: r.PostForm.Get("priority"),
		Oups:     r.PostForm.Get("oups"),
		Ameps:    r.PostForm.Get("ameps"),
		Brips:    r.PostForm.Get("brips"),
		Rte:      r.PostForm.Get("rte"),
		Ais:      r.PostForm.Get("ais"),
		Estimate: r.PostForm.Get("estimate"),
		Target:   r.PostForm.Get("target"),
		Status:   r.PostForm.Get("status"),
		Doneby:   r.PostForm.Get("doneby"),
	}

	// app.infos.__ go fetch the Info struct @ database/infos.go
	app.infos.Agent = form.Agent
	app.infos.Material = form.Material
	app.infos.Detail = form.Detail
	app.infos.Event = form.Event
	app.infos.Oups = form.Oups
	app.infos.Ameps = form.Ameps
	app.infos.Brips = form.Brips
	app.infos.Rte = form.Rte
	app.infos.Ais = form.Ais
	app.infos.Estimate = form.Estimate
	app.infos.Target = form.Target
	app.infos.Status = form.Status
	app.infos.Doneby = form.Doneby
	app.infos.Priority, err = strconv.Atoi(form.Priority)
	if err != nil {
		app.notFound(w)
		return
	}

	err = app.infos.InfoUpdate(iID)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/source/%d/info/view/%d",
		sID, iID), http.StatusSeeOther)
}

func (app *application) infoUpload(w http.ResponseWriter, r *http.Request) {

	data := app.newTemplateData(r)
	app.render(w, http.StatusOK, "importCSV.tmpl.html", data)
}

func (app *application) infoUploadPost(w http.ResponseWriter, r *http.Request) {
	// Taille max du fichier: 2MB
	r.ParseMultipartForm(2 << 20)

	// Crée handler pour filename, size et headers
	file, handler, err := r.FormFile("myFile")
	if err != nil {
		app.errorLog.Println("Error Retrieving the File")
		app.errorLog.Println(err)
		return
	}

	defer file.Close()
	app.infoLog.Printf("Uploaded File: %+v\n", handler.Filename)
	app.infoLog.Printf("File Size: %+v\n", handler.Size)
	app.infoLog.Printf("MIME Header: %+v\n", handler.Header)

	// Creation du fichier
	dst, err := os.Create("csvFiles/" + handler.Filename)
	defer dst.Close()
	if err != nil {
		app.errorLog.Println(w, err.Error(),
			http.StatusInternalServerError)
		return
	}

	// Copie le fichier transféré dans le système
	if _, err := io.Copy(dst, file); err != nil {
		app.errorLog.Println(w, err.Error(),
			http.StatusInternalServerError)
		return
	}

	// Lance la verification de l'extension et encodage du fichier,
	// si concluant, les données seront transférées dans la BD
	app.csvInfo.VerifyCSV("csvFiles/" + handler.Filename)

	http.Redirect(w, r, fmt.Sprintf("/"), http.StatusSeeOther)
}
