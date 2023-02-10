package main

import (
	"path/filepath"
	"text/template"
	"time"

	"E-CURATIFv3/database"
)

// Template struct to generate and parse data to xxx.tmpl.html files
type templateData struct {
	Source  *database.Source
	Sources []*database.Source
	Info    *database.Info
	Infos   []*database.Info
	Form    any
}

// @ sources and infos tables, "Created" and "Updated" column are
// timestamp (UTC)
// SELECT NOW()::timestamp;
// 2023-02-10 19:28:53.116296
// |________| needed
func humanDate(t time.Time) string {
	return t.Format("02/01/2006")
}

// template.FuncMap object is stored in a global variable
// it facilitate the use of humanDate function
var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	// filepath.Glob get a slice of all paths
	pages, err := filepath.Glob("./ui/html/pages/*.tmpl.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		// Extract file name from full filepath
		name := filepath.Base(page)

		// Create new empty template, use Funcs()
		// to register the template.FuncMap and
		// then parse the file.
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

		// Add template set to the map
		cache[name] = ts
	}

	return cache, nil
}
