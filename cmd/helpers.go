package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
)

// Web status are handle here

// The serverError helper writes and error message
// then sends a generic 500 Internal server Error response to the user.
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())

	app.errorLog.Print(trace)
}

// clientError sends a specific status code and corresponding
// description to the user.
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// notFound do the same as clientError but sends a 404 Not Found
// response to the user.
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

// Allocates memory so that a template can be rendered.
// It checks if the desired template exists before beeing sent
// to http.ResponseWriter.
func (app *application) render(w http.ResponseWriter, status int, page string, data *templateData) {
	// Retrieve the appropriate template set from cache
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist",
			page)
		app.serverError(w, err)
		return
	}

	buf := new(bytes.Buffer)

	// Executes the template set and write the response body
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, err)
	}

	w.WriteHeader(status)

	buf.WriteTo(w)

}

// newTemplateData returns a pointer to a templateData not initialized
// and is used on all Handler functions. Makes code more readable
func (app *application) newTemplateData(r *http.Request) *templateData {
	return &templateData{}
}
