package handler

import (
	"log"
	"net/http"

	"github.com/adocoder12/webDevelopment/internal/repository"
)

type Application struct {
	ErrorLog *log.Logger
	InfoLog  *log.Logger
	UserRepo repository.UserRepository
	Renderer *TemplateRenderer // Specialist nested here
}

// NewApplication now accepts the Specialist Renderer
func NewApplication(
	errorLog *log.Logger,
	infoLog *log.Logger,
	userRepo repository.UserRepository,
	renderer *TemplateRenderer,
) *Application {
	return &Application{
		ErrorLog: errorLog,
		InfoLog:  infoLog,
		UserRepo: userRepo,
		Renderer: renderer,
	}
}

func (app *Application) serverError(w http.ResponseWriter, err error) {
	// 1. Log the detailed error for the developer (private)
	// We use Printf with %+v to get as much detail as possible
	app.ErrorLog.Printf("%+v", err)

	// 2. Send a generic message to the user (public)
	// Never show raw DB or system errors to the user
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *Application) clientError(w http.ResponseWriter, status int) {
	// Standardized way to send 400, 401, 403, etc.
	http.Error(w, http.StatusText(status), status)
}

func (app *Application) notFound(w http.ResponseWriter) {
	// Specialized helper for 404s
	app.clientError(w, http.StatusNotFound)
}
