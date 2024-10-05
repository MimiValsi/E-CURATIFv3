package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"E-CURATIFv3/database"
)

// Struct Template qui génère et analyse la data
// vers les fichiers .tmpl.html
type templateData struct {
	Source  *database.Source
	Sources []*database.Source

	Info  *database.Info
	Infos []*database.Info

	JSource []byte

	Flash string
	Form  any
}

// @ tables sources et infos, colonnes "Created" et "Updated"
// ont un timestamp (UTC)
// SELECT NOW()::timestamp;
// 2023-02-10 19:28:53.116296
// |________| besoin
func humanDate(t time.Time) string {
	return t.Format("02/01/2006")
}

// l'Objet template.FuncMap est stpcké dans une variable global
// afin de faciliter l'utilisation de la fonction humanDate
var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	// filepath.Glob crée une slice de tous les chemins
	pages, err := filepath.Glob("./ui/html/pages/*.gotpl.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		// Extrait le nom du fichier du chemin du fichier
		name := filepath.Base(page)

		// Crée une nouvelle template vide, l'ajout de Funcs()
		// sert à enregistrer le template.FuncMap
		// et analyse le fichier
		ts, err := template.New(name).Funcs(functions).
			ParseFiles("./ui/html/base.gotpl.html")
		if err != nil {
			return nil, err
		}

		// ts, err = ts.ParseGlob("./ui/html/partials/*.html.gotpl")
		// if err != nil {
		// 	return nil, err
		// }

		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// Ajoute le set de template ver le map
		cache[name] = ts
	}

	return cache, nil
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request,
	dst any,
) error {
	// Limit the size of the request body to 1MB
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	// Decode the request body into the target dst.
	err := dec.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		var maxBytesError *http.MaxBytesError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-format JSON (at character %d)", syntaxError.Offset)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-format JSON")

		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)

		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		case strings.HasPrefix(err.Error(), "json: unknown field"):
			fieldName := strings.TrimPrefix(err.Error(),
				"json: unknown field")
			return fmt.Errorf("body contains unknown key %s", fieldName)

		case errors.As(err, &maxBytesError):
			return fmt.Errorf("body must not be larger than %d bytes",
				maxBytesError.Limit)

		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain as single JSON value")
	}

	return nil
}
