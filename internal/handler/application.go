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
	app.ErrorLog.Printf("%v", err)
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}
