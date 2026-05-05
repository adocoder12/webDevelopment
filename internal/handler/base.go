package handler

import "net/http"

// Home handles the landing page
func (app *Application) Home(w http.ResponseWriter, r *http.Request) {
	app.renderTemplate(w, r, "index.html", nil)
}

func (app *Application) About(w http.ResponseWriter, r *http.Request) {
	app.renderTemplate(w, r, "about.html", nil)
}
