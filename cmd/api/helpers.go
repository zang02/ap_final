package main

import (
	"app/internal/data"
	"bytes"
	"fmt"
	"net/http"
	"time"
)

func (app *application) render(w http.ResponseWriter, r *http.Request, name string, td *data.TemplateData) {
	// Retrieve the appropriate template set from the cache based on the page name
	// (like 'home.page.gohtml'). If no entry exists in the cache with the provided name,
	// call the serverError helper method.
	ts, ok := app.templateCache[name]
	if !ok {
		app.serverError(w, fmt.Errorf("the template %s does not exist", name))
		return
	}

	// Initialize a new buffer.
	buff := new(bytes.Buffer)

	// Write the template to the buffer, instead of straight to the http.ResponseWriter.
	// If there is an error, call our serverError helper and then return.
	td.CurrentYear = fmt.Sprintf("%v", time.Now().Year())
	// td.Flash = app.session.PopString(r, "flash")
	// td.IsAuthenticated = app.isAuthenticated(r)
	// td.CSRFToken = nosurf.Token(r)

	err := ts.Execute(buff, td)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Write the contents of the buffer to the http.ResponseWriter. Again, this is another place
	// where we pass our http.ResponseWriter to a function that take an io.Writer
	if _, err = buff.WriteTo(w); err != nil {
		app.serverError(w, err)
		return
	}
}

func (app *application) serverError(w http.ResponseWriter, err error) {
	app.logger.PrintError(err.Error(), "server error")

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
