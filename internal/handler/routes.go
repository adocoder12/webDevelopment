package handler

import "net/http"

// SetupRoutes attaches all your URLs to the Application
func (app *Application) SetupRoutes() http.Handler {
	mux := http.NewServeMux()

	// Static Files
	publicDir := "./templates/html/ui/public"
	fileServer := http.FileServer(http.Dir(publicDir))
	mux.Handle("/public/", http.StripPrefix("/public", fileServer))

	// Public routes
	mux.HandleFunc("GET /{$}", app.Home)
	mux.HandleFunc("GET /register", app.RegisterHandler)
	mux.HandleFunc("POST /register", app.RegisterHandler)
	mux.HandleFunc("GET /login", app.LoginHandler)
	mux.HandleFunc("POST /login", app.LoginHandler)
	mux.HandleFunc("POST /logout", app.LogoutHandler)

	// Protected routes (add these later)
	// mux.HandleFunc("GET /dashboard", app.AuthRequired(app.Dashboard))

	return app.RecoverPanic(app.SessionManager.LoadAndSave(app.Logger(mux)))
}
