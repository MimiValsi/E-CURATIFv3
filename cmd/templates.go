package main

import (
	"path/filepath"
	"text/template"
	"time"

	"E-CURATIFv3/database"
)

// Struct Template qui génère et analyse la data
// vers les fichiers .tmpl.html
type templateData struct {
	Source  *database.Source
	Sources []*database.Source
	Info    *database.Info
	Infos   []*database.Info
	Form    any
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
	pages, err := filepath.Glob("./ui/html/pages/*.tmpl.html")
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
			ParseFiles("./ui/html/base.tmpl.html")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl.html")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// Ajoute le set de template ver le map
		cache[name] = ts
	}

	return cache, nil
}
