package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"E-CURATIFv3/database"
	"E-CURATIFv3/internal/validator"
)

func (app *application) dbConn(ctx context.Context) *pgxpool.Conn {
	conn, err := app.DB.Acquire(ctx)
	if err != nil {
		app.errorLog.Println("Couldn't connect to DB")
		app.errorLog.Println(err)
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
		app.serverError(w, r, err)
		return
	}

	jData, err := json.Marshal(sources)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// newTemplateData @ cmd/templates.go
	data := app.newTemplateData(r)
	data.Sources = sources
	data.JSource = jData

	app.render(w, r, http.StatusOK, "home.gotpl.html", data)
}

func (app *application) priorityData(w http.ResponseWriter, r *http.Request) {
	conn := app.dbConn(r.Context())
	defer conn.Release()

	infos, err := app.infos.PriorityInfos(conn)
	if err != nil {
		app.serverError(w, r, err)
	}

	jsonData, err := json.Marshal(infos)
	if err != nil {
		app.serverError(w, r, err)
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
		app.serverError(w, r, err)
	}

	jsonGraph, err := json.Marshal(sources)
	if err != nil {
		app.serverError(w, r, err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
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
			app.serverError(w, r, err)
		}
	}

	// Fait appel à la fonction dans database/infos.go
	// en prenant en compte le id du source
	info, err := app.infos.List(id, conn)
	if err != nil {
		if errors.Is(err, database.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, r, err)
		}
	}

	// Allocation de mémoire pour création de template
	data := app.newTemplateData(r)
	data.Infos = info
	data.Source = source

	// Génération de la page web
	app.render(w, r, http.StatusOK, "sourceView.gotpl.html", data)
}

// Génération de la page de création de source
// Celle-ci est faite en deux parties:
// La première est une GET, une fois le nom du source choisi
// la fonction sourceCreatePost prend le relais
func (app *application) sourceCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = sourceCreateForm{}

	app.render(w, r, http.StatusOK, "sourceCreate.gotpl.html", data)
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
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "sourceCreate.gotpl.html", data)
		return
	}

	// Si pas d'erreur, les données sont envoyés vers la BD
	id, err := app.sources.Insert(form.Name, form.CodeGmao, conn)
	if err != nil {
		app.serverError(w, r, err)
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
			app.serverError(w, r, err)
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
			app.serverError(w, r, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Source = source

	app.render(w, r, http.StatusOK, "sourceUpdate.gotpl.html", data)
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
		app.render(w, r, http.StatusUnprocessableEntity, "sourceUpdate.gotpl.html", data)
		return
	}

	app.sources.Name = form.Name

	err = app.sources.Update(id, conn)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/source/view/%d", id),
		http.StatusSeeOther)
}

//
// Infos Handlers
//

type infoCreateForm struct {
	ID          int
	Agent       string
	Ouvrage     string
	Priorite    string
	Echeance    string
	Detail      string
	Created     string
	Updated     string
	Status      string
	Evenement   string
	Commentaire string
	Entite      string

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
			app.serverError(w, r, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Form = infoCreateForm{}
	data.Source = source

	app.render(w, r, http.StatusOK, "infoCreate.gotpl.html", data)
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
		Ouvrage:     r.PostForm.Get("ouvrage"),
		Detail:      r.PostForm.Get("detail"),
		Evenement:   r.PostForm.Get("evenement"),
		Priorite:    r.PostForm.Get("priorite"),
		Echeance:    r.PostForm.Get("echeance"),
		Status:      r.PostForm.Get("status"),
		Entite:      r.PostForm.Get("entite"),
		Commentaire: r.PostForm.Get("commentaire"),
		Created:     r.PostForm.Get("created"),
	}

	app.infos.Ouvrage = form.Ouvrage
	app.infos.Detail = form.Detail
	app.infos.Evenement = form.Evenement
	app.infos.Priorite, err = strconv.Atoi(form.Priorite)
	if err != nil {
		app.notFound(w)
		return
	}
	app.infos.Echeance = form.Echeance
	app.infos.Status = form.Status
	app.infos.Entite = form.Entite
	app.infos.Commentaire = form.Commentaire

	if form.Created == "" {
		app.infos.Created = time.Now().UTC()
	} else {
		app.infos.Created, err = time.Parse("02/01/2006", form.Created)
		if err != nil {
			log.Printf("Format de date invalide: %v", form.Created)
			log.Println(err)
			return
		}
	}

	iid, err := app.infos.Insert(sID, conn)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Info crée avec succès!")

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
			app.serverError(w, r, err)
		}
		return
	}

	flash := app.sessionManager.PopString(r.Context(), "flash")

	data := app.newTemplateData(r)
	data.Info = info
	data.Flash = flash

	app.render(w, r, http.StatusOK, "infoView.gotpl.html", data)
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
			app.serverError(w, r, err)
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
			app.serverError(w, r, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Info = info

	app.render(w, r, http.StatusOK, "infoUpdate.gotpl.html", data)
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
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/source/%d/info/view/%d",
		sID, iID), http.StatusSeeOther)
}

// Import CSV data functionality
// Copy the file content, check it's encoding
// and send data to DB
// func (app *application) importCSV(w http.ResponseWriter, r *http.Request) {
// 	data := app.newTemplateData()
// 	app.render(w, http.StatusOK, "importCSV.gotpl.html", data)
// }

func (app *application) importCSVPost(w http.ResponseWriter, r *http.Request) {
	conn := app.dbConn(r.Context())
	defer conn.Release()

	// Max file size: 2MB
	r.ParseMultipartForm(2_000_000)

	// Create handler for file name, size and headers
	file, handler, err := r.FormFile("inpt")
	if err != nil {
		app.errorLog.Println("Error Retrieving the File")
		app.errorLog.Println(err)
		return
	}

	defer file.Close()

	// Create file
	dst, err := os.Create("csvFiles/" + handler.Filename)
	if err != nil {
		app.errorLog.Println(w, err.Error(),
			http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copy transfered file to the system
	if _, err := io.Copy(dst, file); err != nil {
		app.errorLog.Println(w, err.Error(),
			http.StatusInternalServerError)
		return
	}

	// Start the file encoding verification and if all good
	// send data to DB
	app.csvImport.VerifyCSV("csvFiles/"+handler.Filename, conn)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) exportCSVPost(w http.ResponseWriter, r *http.Request) {
	conn := app.dbConn(r.Context())
	defer conn.Release()

	path, err := app.csvExport.Export_DB_csv(conn)
	if err != nil {
		app.errorLog.Println("Couldn't fetch file")
		return
	}

	app.infoLog.Println("Export en cours")

	file, err := os.ReadFile(path)
	if err != nil {
		app.errorLog.Println("File doesn't exist")
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename=export.csv")
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Length", fmt.Sprint(len(file)))

	flusher, ok := w.(http.Flusher)
	if !ok {
		app.errorLog.Println("Cannot use flusher")
	}
	w.Write(file)

	flusher.Flush()

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

type userSignupForm struct {
	Name     string
	Email    string
	Password string

	validator.Validator
}

// users
func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userSignupForm{}

	app.render(w, r, http.StatusOK, "signup.gotpl.html", data)
}

func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	conn := app.dbConn(r.Context())
	defer conn.Release()

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
	}

	form := userSignupForm{
		Name:     r.PostForm.Get("name"),
		Email:    r.PostForm.Get("email"),
		Password: r.PostForm.Get("password"),
	}

	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "signup.gotpl.html", data)
		return
	}

	// app.users.Name = form.Name
	// app.users.Email = form.Email
	// app.users.HashedPassword = []byte(form.Password)

	err = app.users.Insert(form.Name, form.Email, form.Password, conn)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Vous vous êtes bien enregistré")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

type userLoginForm struct {
	Email    string
	Password string

	validator.Validator
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userLoginForm{}

	app.render(w, r, http.StatusOK, "login.gotpl.html", data)
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	conn := app.dbConn(r.Context())
	defer conn.Release()

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
	}

	form := userLoginForm{
		Email:    r.PostForm.Get("email"),
		Password: r.PostForm.Get("password"),
	}

	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "login.gotpl.html", data)
		return
	}

	// app.users.Email = form.Email
	// app.users.HashedPassword = []byte(form.Password)

	id, err := app.users.Authenticate(form.Email, form.Password, conn)
	if err != nil {
		if errors.Is(err, database.ErrInvalidCredentials) {
			form.AddNonFieldError("Email or password is incorrect")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, r, http.StatusUnprocessableEntity, "login.gotpl.html", data)
			return
		} else {
			app.serverError(w, r, err)
		}

		return
	}

	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	err := app.sessionManager.RenewToken(ctx)
	if err != nil {
		app.serverError(w, r, err)
	}

	app.sessionManager.Remove(ctx, "authenticatedUserID")

	app.sessionManager.Put(ctx, "flash", "Déconnecté avec succès!")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
