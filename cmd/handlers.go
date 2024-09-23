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

	// "time"

	"E-CURATIFv3/database"
	"E-CURATIFv3/internal/validator"

	// package pour les routers
	"github.com/go-chi/chi/v5"

	// pkg pour Psql driver
	"github.com/jackc/pgx/v5/pgxpool"
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
	data := app.newTemplateData()
	data.Sources = sources
	data.JSource = jData

	app.render(w, http.StatusOK, "home.gotpl.html", data)
}

func (app *application) priorityData(w http.ResponseWriter, r *http.Request) {
	conn := app.dbConn(r.Context())
	defer conn.Release()

	infos, err := app.infos.PriorityInfos(conn)
	if err != nil {
		app.serverError(w, err)
	}

	jsonData, err := json.Marshal(infos)
	if err != nil {
		app.serverError(w, err)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
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

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonGraph)
}

// func (app *application) charts(w http.ResponseWriter, r *http.Request) {
// 	app.render(w, http.StatusOK, "charts.gotpl.html", nil)
// }

func (app *application) curatifDone(w http.ResponseWriter, r *http.Request) {
	conn := app.dbConn(r.Context())
	defer conn.Release()

	sources, err := app.sources.CuratifsDone(conn)
	if err != nil {
		app.serverError(w, err)
	}

	jsonGraph, err := json.Marshal(sources)
	if err != nil {
		app.serverError(w, err)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonGraph)
}

//
// Sources Handlers
//

type sourceCreateForm struct {
	Name     string
	CodeGmao string

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
	source, err := app.sources.Get(id, conn)
	if err != nil {
		if errors.Is(err, database.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
	}

	// Fait appel à la fonction dans database/infos.go
	// en prenant en compte le id du source
	info, err := app.infos.List(id, conn)
	if err != nil {
		if errors.Is(err, database.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
	}

	// Allocation de mémoire pour création de template
	data := app.newTemplateData()
	data.Infos = info
	data.Source = source

	// Génération de la page web
	app.render(w, http.StatusOK, "sourceView.gotpl.html", data)
}

// Génération de la page de création de source
// Celle-ci est faite en deux parties:
// La première est une GET, une fois le nom du source choisi
// la fonction sourceCreatePost prend le relais
func (app *application) sourceCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData()
	data.Form = sourceCreateForm{}

	app.render(w, http.StatusOK, "sourceCreate.gotpl.html", data)
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
		Name:     r.PostForm.Get("name"),
		CodeGmao: r.PostForm.Get("code_gmao"),
	}

	// Petite vérification que le champ ne soit pas vide.
	// Ceci empêche que la BD génère une erreur
	emptyField := "Ce champ ne doit pas être vide"

	form.CheckField(validator.NotBlank(form.Name),
		"name", emptyField)
	form.CheckField(validator.NotBlank(form.CodeGmao),
		"code_gmao", emptyField)

	if !form.Valid() {
		data := app.newTemplateData()
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity,
			"sourceCreate.gotpl.html", data)
		return
	}

	// Si pas d'erreur, les données sont envoyés vers la BD
	id, err := app.sources.Insert(form.Name, form.CodeGmao, conn)
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

	err = app.sources.Delete(id, conn)
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

	source, err := app.sources.Get(id, conn)
	if err != nil {
		if errors.Is(err, database.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	data := app.newTemplateData()
	data.Source = source

	app.render(w, http.StatusOK, "sourceUpdate.gotpl.html", data)
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
		data := app.newTemplateData()
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity,
			"sourceUpdate.gotpl.html", data)
		return
	}

	app.sources.Name = form.Name

	err = app.sources.Update(id, conn)
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
	ID        int
	Agent     string
	Ouvrage   string
	Priorite  string
	Echeance  string
	Detail    string
	Created   string
	Updated   string
	Status    string
	Evenement string
	Devis     string
	// Oups        string
	FaitPar     string
	Commentaire string

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
	source, err := app.sources.Get(id, conn)
	if err != nil {
		if errors.Is(err, database.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	data := app.newTemplateData()
	data.Form = infoCreateForm{}
	data.Source = source

	app.render(w, http.StatusOK, "infoCreate.gotpl.html", data)
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
		Agent:     r.PostForm.Get("agent"),
		Ouvrage:   r.PostForm.Get("ouvrage"),
		Detail:    r.PostForm.Get("detail"),
		Evenement: r.PostForm.Get("evenement"),
		Priorite:  r.PostForm.Get("priorite"),
		Devis:     r.PostForm.Get("devis"),
		Echeance:  r.PostForm.Get("echeance"),
		Status:    r.PostForm.Get("status"),
		// FaitPar:    r.PostForm.Get("fait_par"),
	}

	// Certains champs ne doivent pas être vides.
	// Afin de ne pas recevoir une erreur venant de la BD
	// On la vérifie en avance.
	// Ceci sera fait en JS en amont en plus
	emptyField := "Ce champ ne doit pas être vide"

	form.CheckField(validator.NotBlank(form.Agent),
		"agent", emptyField)
	form.CheckField(validator.NotBlank(form.Ouvrage),
		"ouvrage", emptyField)
	form.CheckField(validator.NotBlank(form.Detail),
		"detail", emptyField)
	form.CheckField(validator.NotBlank(form.Evenement),
		"evenement", emptyField)
	form.CheckField(validator.NotBlank(form.Priorite),
		"priorite", emptyField)
	form.CheckField(validator.NotBlank(form.Status),
		"status", emptyField)

	if !form.Valid() {
		data := app.newTemplateData()
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity,
			"infoCreate.gotpl.html", data)
		return
	}

	app.infos.Agent = form.Agent
	app.infos.Ouvrage = form.Ouvrage
	app.infos.Detail = form.Detail
	app.infos.Evenement = form.Evenement
	app.infos.Echeance = form.Echeance
	app.infos.Status = form.Status
	app.infos.Priorite, err = strconv.Atoi(form.Priorite)
	if err != nil {
		app.notFound(w)
		return
	}

	iid, err := app.infos.Insert(sID, conn)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/source/%d/info/view/%d",
		sID, iid), http.StatusSeeOther)
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

	info, err := app.infos.Get(id, conn)
	if err != nil {
		if errors.Is(err, database.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	data := app.newTemplateData()
	data.Info = info

	app.render(w, http.StatusOK, "infoView.gotpl.html", data)
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

	err = app.infos.Delete(id, conn)
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
	info, err := app.infos.Get(id, conn)
	if err != nil {
		if errors.Is(err, database.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	data := app.newTemplateData()
	data.Info = info

	app.render(w, http.StatusOK, "infoUpdate.gotpl.html", data)
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
		Agent:     r.PostForm.Get("agent"),
		Ouvrage:   r.PostForm.Get("ouvrage"),
		Detail:    r.PostForm.Get("detail"),
		Evenement: r.PostForm.Get("evenement"),
		Priorite:  r.PostForm.Get("priorite"),
		Echeance:  r.PostForm.Get("echeance"),
		Status:    r.PostForm.Get("status"),
	}

	app.infos.Agent = form.Agent
	app.infos.Ouvrage = form.Ouvrage
	app.infos.Detail = form.Detail
	app.infos.Evenement = form.Evenement
	app.infos.Echeance = form.Echeance
	app.infos.Status = form.Status
	app.infos.Priorite, err = strconv.Atoi(form.Priorite)
	if err != nil {
		app.notFound(w)
		return
	}

	err = app.infos.Update(iID, conn)
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
	data := app.newTemplateData()
	app.render(w, http.StatusOK, "importCSV.gotpl.html", data)
}

func (app *application) importCSVPost(w http.ResponseWriter, r *http.Request) {
	conn := app.dbConn(r.Context())
	defer conn.Release()

	// Taille max du fichier: 2MB
	r.ParseMultipartForm(2_000_000)

	// Crée handler pour filename, size et headers
	file, handler, err := r.FormFile("inpt")
	if err != nil {
		app.errorLog.Println("Error Retrieving the File")
		app.errorLog.Println(err)
		return
	}

	defer file.Close()

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
	app.csvData.VerifyCSV("csvFiles/"+handler.Filename, conn)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// func (app *application) pageTest(w http.ResponseWriter, r *http.Request) {
// 	conn := app.dbConn(r.Context())
// 	defer conn.Release()
//
// 	csv, err := app.csvData.Export(conn)
// 	if err != nil {
// 		app.serverError(w, err)
// 	}
//
// 	jsonData, err := json.Marshal(csv)
// 	if err != nil {
// 		app.serverError(w, err)
// 	}
//
// 	w.WriteHeader(http.StatusOK)
// 	w.Header().Set("Content-Type", "application/json")
// 	w.Write(jsonData)
// }
