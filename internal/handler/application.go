package handler

import (
	"log"
	"net/http"

	"github.com/adocoder12/webDevelopment/internal/repository"
	"github.com/alexedwards/scs/v2"
)

type Application struct {
	ErrorLog       *log.Logger
	InfoLog        *log.Logger
	UserRepo       repository.UserRepository
	Renderer       *TemplateRenderer
	SessionManager *scs.SessionManager
}

// The struct constructor
func NewApplication(
	errorLog *log.Logger,
	infoLog *log.Logger,
	userRepo repository.UserRepository,
	renderer *TemplateRenderer,
	sessionManager *scs.SessionManager,
) *Application {
	return &Application{
		ErrorLog:       errorLog,
		InfoLog:        infoLog,
		UserRepo:       userRepo,
		Renderer:       renderer,
		SessionManager: sessionManager,
	}
}

func (app *Application) serverError(w http.ResponseWriter, err error) {
	// 1. Log the detailed error for the developer (private)
	app.ErrorLog.Printf("%+v", err)

	// Never show raw DB or system errors to the user
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *Application) clientError(w http.ResponseWriter, status int) {
	// Standardized way to send 400, 401, 403, etc.
	http.Error(w, http.StatusText(status), status)
}

func (app *Application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}
