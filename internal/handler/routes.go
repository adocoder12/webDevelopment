package handler

import "net/http"

// SetupRoutes attaches all your URLs to the Application
func (app *Application) SetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	// Static Files
	publicDir := "./templates/html/ui/public"

	fileServer := http.FileServer(http.Dir(publicDir))

	// This tells Go: when the browser asks for "/public/...",
	// look inside the publicDir folder.
	mux.Handle("/public/", http.StripPrefix("/public", fileServer))

	// Web Pages
	mux.HandleFunc("GET /{$}", app.Home)

	// Registration Routes
	mux.HandleFunc("GET /register", app.RegisterHandler)  // Show form
	mux.HandleFunc("POST /register", app.RegisterHandler) // Process form

	// Login Routes
	mux.HandleFunc("GET /login", app.LoginHandler)  // Show form
	mux.HandleFunc("POST /login", app.LoginHandler) // Process form

	return mux
}
