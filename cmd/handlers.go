package main

import (
	"E-CURATIFv3/database"
	"E-CURATIFv3/internal/validator"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	sources, err := app.sources.MenuSource()
	if err != nil {
		app.serverError(w, err)
		return
	}

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

func (app *application) sourceView(w http.ResponseWriter, r *http.Request) {
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

	infos, err := app.infos.InfoList(id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Infos = infos
	data.Source = source

	app.render(w, http.StatusOK, "sourceView.tmpl.html", data)

}

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

	champVide := "Ce champ ne doit pas être vide"

	form.CheckField(validator.NotBlank(form.Name), "name", champVide)
	form.CheckField(validator.MaxChars(form.Name, 20), "name", "Nom de Poste Source trop grand")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "sourceCreate.tmpl.html", data)
		return
	}

	id, err := app.sources.SourceInsert(form.Name)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/source/view/%d", id), http.StatusSeeOther)
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
	// target   string
	validator.Validator
}

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
	}

	champVide := "Ce champ ne doit pas être vide"

	form.CheckField(validator.NotBlank(form.Agent), "agent", champVide)
	form.CheckField(validator.NotBlank(form.Material), "material", champVide)
	form.CheckField(validator.NotBlank(form.Detail), "detail", champVide)
	form.CheckField(validator.NotBlank(form.Event), "event", champVide)
	form.CheckField(validator.NotBlank(form.Priority), "priority", champVide)
	form.CheckField(validator.NotBlank(form.Status), "status", champVide)

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "infoCreate.tmpl.html", data)
		return
	}

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

	http.Redirect(w, r, fmt.Sprintf("/source/%d/info/view/%d", sID, id), http.StatusSeeOther)
}

func (app *application) infoView(w http.ResponseWriter, r *http.Request) {
	iKey := chi.URLParam(r, "id")

	iID, err := strconv.Atoi(iKey)
	if err != nil || iID < 1 {
		app.notFound(w)
		return
	}

	info, err := app.infos.InfoGet(iID)
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

	http.Redirect(w, r, fmt.Sprintf("/source/view/%d", sID), http.StatusSeeOther)

}
