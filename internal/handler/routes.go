package handler

import (
	"net/http"
)

// SetupRoutes attaches all your URLs to the Application
func (app *Application) SetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	// 1. Static File Server
	// Ensure this points to your actual CSS/Assets directory
	publicDir := "./templates/html/ui/public"
	fileServer := http.FileServer(http.Dir(publicDir))
	mux.Handle("/public/", http.StripPrefix("/public", fileServer))

	// 2. The Fix: Add '{$}' to the root pattern
	mux.HandleFunc("GET /{$}", app.Home)

	// 3. API Routes
	mux.HandleFunc("POST /api/v1/register", app.RegisterHandler)
	mux.HandleFunc("GET /api/v1/users", app.ListUsersHandler)
	mux.HandleFunc("GET /api/v1/users/{id}", app.GetByIDHandler)

	return mux
}
