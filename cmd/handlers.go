package main

import (
	"context"
	"encoding/json"
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
	"github.com/jackc/pgx/v4/pgxpool"
)

func (app *application) dbConn(ctx context.Context) *pgxpool.Conn {
	conn, err := app.DB.Acquire(ctx)
	if err != nil {
		app.errorLog.Println("Couldn't connect to DB")
		return nil
	}

	return conn
}

//
// Home
//

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Content-Type", "application/json")
	conn := app.dbConn(r.Context())
	defer conn.Release()

	// MenuSource func @ database/sources.go
	sources, err := app.sources.MenuSource(conn)
	if err != nil {
		app.serverError(w, err)
		return
	}

	jData, err := json.Marshal(sources)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// newTemplateData @ cmd/templates.go
	data := app.newTemplateData(r)
	data.Sources = sources
	data.JSource = jData

	app.render(w, http.StatusOK, "home.html.gotpl", data)

}

func (app *application) jsonData(w http.ResponseWriter, r *http.Request) {
	conn := app.dbConn(r.Context())
	defer conn.Release()

	sources, err := app.sources.MenuSource(conn)
	if err != nil {
		app.serverError(w, err)
	}

	jsonGraph, err := json.Marshal(sources)
	if err != nil {
		app.serverError(w, err)
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonGraph)
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
	conn := app.dbConn(r.Context())
	defer conn.Release()

	key := chi.URLParam(r, "id")
	id, err := strconv.Atoi(key)
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	// Fait appel à la fonctin dans database/sources.go
	// On récupère le "id" du source dans le URL
	// créé auparavant
	source, err := app.sources.SourceGet(id, conn)
	if err != nil {
		if errors.Is(err, database.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
	}

	// Fait appel à la fonction dans database/infos.go
	// en prenant en compte le id du source
	info, err := app.infos.InfoList(id, conn)
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
	app.render(w, http.StatusOK, "sourceView.html.gotpl", data)

}

// Génération de la page de création de source
// Celle-ci est faite en deux parties:
// La première est une GET, une fois le nom du source choisi
// la fonction sourceCreatePost prend le relais
func (app *application) sourceCreate(w http.ResponseWriter, r *http.Request) {

	data := app.newTemplateData(r)
	data.Form = sourceCreateForm{}

	app.render(w, http.StatusOK, "sourceCreate.html.gotpl", data)
}

func (app *application) sourceCreatePost(w http.ResponseWriter, r *http.Request) {

	conn := app.dbConn(r.Context())
	defer conn.Release()

	// ParseForm analyse le URL et le rend dispo
	// afin de récuperer des infos avec PostForm.Get
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
	}

	form := sourceCreateForm{
		Name: r.PostForm.Get("name"),
	}

	// Petite vérification que le champ ne soit pas vide.
	// Ceci empêche que la BD génère une erreur
	emptyField := "Ce champ ne doit pas être vide"

	form.CheckField(validator.NotBlank(form.Name),
		"name", emptyField)

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity,
			"sourceCreate.html.gotpl", data)
		return
	}

	// Si pas d'erreur, les données sont envoyés vers la BD
	id, err := app.sources.SourceInsert(form.Name, conn)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/source/view/%d", id),
		http.StatusSeeOther)
}

// La fonction récupère l'id du Source. L'id est un string
// qui sera converti en tant que int.
// Cette fonction n'est dispo que si le Source existe, même si on
// vérifi son existance par acquit de conscience.
// Une fois fini, on redirect vers la page d'accueil
func (app *application) sourceDeletePost(w http.ResponseWriter, r *http.Request) {

	conn := app.dbConn(r.Context())
	defer conn.Release()

	key := chi.URLParam(r, "id")

	id, err := strconv.Atoi(key)
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	err = app.sources.SourceDelete(id, conn)
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

// Même foncionnement que la création et suppréssion de Source.
// On récupère l'id et on extrait les données de la BD
// afin de changer le nom.
func (app *application) sourceUpdate(w http.ResponseWriter, r *http.Request) {
	conn := app.dbConn(r.Context())
	defer conn.Release()

	key := chi.URLParam(r, "id")
	id, err := strconv.Atoi(key)
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	source, err := app.sources.SourceGet(id, conn)
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

	app.render(w, http.StatusOK, "sourceUpdate.html.gotpl", data)
}

// Une fois les changements faites, elles sont reenvoyés vers la BD
func (app *application) sourceUpdatePost(w http.ResponseWriter, r *http.Request) {

	conn := app.dbConn(r.Context())
	defer conn.Release()

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
			"sourceUpdate.html.gotpl", data)
		return
	}

	app.sources.Name = form.Name

	err = app.sources.SourceUpdate(id, conn)
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
	ID       int
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

// Même fonctionnement que "sourceCreate", la différence principale
// est qu'on doit récuperer l'id du Source.
//
// Dans la tableau "Infos" de la BD, source_id = FK de id (Source)
func (app *application) infoCreate(w http.ResponseWriter, r *http.Request) {

	conn := app.dbConn(r.Context())
	defer conn.Release()

	// Ici on récupère l'id du Source
	key := chi.URLParam(r, "id")

	id, err := strconv.Atoi(key)
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	// La fonction attend comme paramètre un "int" dont l'id Source
	source, err := app.sources.SourceGet(id, conn)
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

	app.render(w, http.StatusOK, "infoCreate.html.gotpl", data)
}

func (app *application) infoCreatePost(w http.ResponseWriter, r *http.Request) {

	conn := app.dbConn(r.Context())
	defer conn.Release()

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
	}

	key := chi.URLParam(r, "id")

	sID, err := strconv.Atoi(key)
	if err != nil || sID < 1 {
		app.notFound(w)
		return
	}

	// Les données récupérés depuis la page HTML sont envoyées
	// vers la BD
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

	// Certains champs ne doivent pas être vides.
	// Afin de ne pas recevoir une erreur venant de la BD
	// On la vérifie en avance.
	// Ceci sera fait en JS en amont en plus
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
			"infoCreate.html.gotpl", data)
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
	app.infos.Doneby = form.Doneby
	app.infos.Priority, err = strconv.Atoi(form.Priority)
	if err != nil {
		app.notFound(w)
		return
	}

	_, err = app.infos.Insert(sID, conn)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/source/%d/info/create",
		sID), http.StatusSeeOther)
}

// Page permettant de visualiser en détails la fiche curatif
// de l'ouvrage concerné
func (app *application) infoView(w http.ResponseWriter, r *http.Request) {
	conn := app.dbConn(r.Context())
	defer conn.Release()

	iKey := chi.URLParam(r, "id")

	id, err := strconv.Atoi(iKey)
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	info, err := app.infos.InfoGet(id, conn)
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

	app.render(w, http.StatusOK, "infoView.html.gotpl", data)
}

// HTML POST afin de supprimer le curatif(info)
func (app *application) infoDeletePost(w http.ResponseWriter, r *http.Request) {

	conn := app.dbConn(r.Context())
	defer conn.Release()

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

	err = app.infos.InfoDelete(id, conn)
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

// Même fonctionnement que "sourceUpdate"
func (app *application) infoUpdate(w http.ResponseWriter, r *http.Request) {
	conn := app.dbConn(r.Context())
	defer conn.Release()

	key := chi.URLParam(r, "id")
	id, err := strconv.Atoi(key)
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	// id ~> Info id
	info, err := app.infos.InfoGet(id, conn)
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

	app.render(w, http.StatusOK, "infoUpdate.html.gotpl", data)
}

func (app *application) infoUpdatePost(w http.ResponseWriter, r *http.Request) {

	conn := app.dbConn(r.Context())
	defer conn.Release()

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
	}

	// Etant donnée qu'un curatif(Info) est attaché à un source,
	// on doit le récuperer
	sKey := chi.URLParam(r, "sid")
	sID, err := strconv.Atoi(sKey)
	if err != nil || sID < 1 {
		app.notFound(w)
		return
	}

	// ainsi que l'id (forcément)
	iKey := chi.URLParam(r, "id")
	iID, err := strconv.Atoi(iKey)
	if err != nil || iID < 1 {
		app.notFound(w)
		return
	}

	// Même fonctionnement que "infoCreate", on récupère tout
	// et on envoi ce qui a été modifié
	// Ceci sera traîtré dans "database.go"
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

	err = app.infos.InfoUpdate(iID, conn)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/source/%d/info/view/%d",
		sID, iID), http.StatusSeeOther)
}

// Page HTML en cours de création
// en soit permet de d'ouvrir et lire des fichiers .csv
// A l'heure actuelle les données ne sont pas importés correctement
func (app *application) importCSV(w http.ResponseWriter, r *http.Request) {

	data := app.newTemplateData(r)
	app.render(w, http.StatusOK, "importCSV.html.gotpl", data)
}

func (app *application) importCSVPost(w http.ResponseWriter, r *http.Request) {
	conn := app.dbConn(r.Context())
	defer conn.Release()

	// Taille max du fichier: 2MB
	r.ParseMultipartForm(2 << 20)

	// Crée handler pour filename, size et headers
	file, handler, err := r.FormFile("inpt")
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
	if err != nil {
		app.errorLog.Println(w, err.Error(),
			http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copie le fichier transféré dans le système
	if _, err := io.Copy(dst, file); err != nil {
		app.errorLog.Println(w, err.Error(),
			http.StatusInternalServerError)
		return
	}

	// Lance la verification de l'extension et encodage du fichier,
	// si concluant, les données seront transférées dans la BD
	app.csvInfo.VerifyCSV("csvFiles/" + handler.Filename)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) pageTest(w http.ResponseWriter, r *http.Request) {

	app.render(w, http.StatusOK, "pageTest.html.gotpl", nil)
}
