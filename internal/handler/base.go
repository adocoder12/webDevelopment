package handler

import "net/http"

// Home handles the landing page
func (app *Application) Home(w http.ResponseWriter, r *http.Request) {

	// Passing nil if no dynamic data is needed for Home
	app.renderTemplate(w, "index.html", nil)
}

// About handles the informational page
func (app *Application) About(w http.ResponseWriter, r *http.Request) {
	app.renderTemplate(w, "about.html", nil)
}
