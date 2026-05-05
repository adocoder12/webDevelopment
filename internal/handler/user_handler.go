package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/adocoder12/webDevelopment/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// -- Flash & auth helpers ---------------------------------------------------

func (app *Application) setFlash(r *http.Request, msg string) {
	//Put to store
	app.SessionManager.Put(r.Context(), "flash", msg)
}

func (app *Application) getFlash(r *http.Request) string {
	//PopString to read-and-delete in one call.
	return app.SessionManager.PopString(r.Context(), "flash")
}

func (app *Application) isAuthenticated(r *http.Request) bool {
	return app.SessionManager.GetInt64(r.Context(), "userID") != 0
}

// -- Handlers ---------------------------------------------------------------

func (app *Application) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		app.renderTemplate(w, r, "register.html", nil)
		return
	}

	if err := r.ParseForm(); err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	_, err := app.UserRepo.CreateUser(
		r.Context(),
		r.PostForm.Get("name"),
		r.PostForm.Get("email"),
		r.PostForm.Get("password"),
		r.PostForm.Get("confirmPassword"),
		"default_avatar.png",
	)
	if err != nil {
		app.renderTemplate(w, r, "register.html", &TemplateData{
			Error: err.Error(),
		})
		return
	}

	app.setFlash(r, "Account created! Please log in.")
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (app *Application) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		app.renderTemplate(w, r, "login.html", nil)
		return
	}

	if err := r.ParseForm(); err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	user, err := app.UserRepo.GetUserByEmail(r.Context(), r.PostForm.Get("email"))
	if err != nil {
		app.renderTemplate(w, r, "login.html", &TemplateData{
			Error: "Invalid email or password",
		})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(r.PostForm.Get("password"))); err != nil {
		app.renderTemplate(w, r, "login.html", &TemplateData{
			Error: "Invalid email or password",
		})
		return
	}
	//This prevents "session fixation" — an attack where someone steals your token before you log in and then reuses it after
	if err := app.SessionManager.RenewToken(r.Context()); err != nil {
		app.serverError(w, err)
		return
	}

	app.SessionManager.Put(r.Context(), "userID", user.Id)
	app.SessionManager.Put(r.Context(), "userName", user.Name)
	app.setFlash(r, "Welcome back, "+user.Name+"!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *Application) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	//This prevents "session fixation" — an attack where someone steals your token before you log in and then reuses it after
	if err := app.SessionManager.RenewToken(r.Context()); err != nil {
		app.serverError(w, err)
		return
	}

	app.SessionManager.Remove(r.Context(), "userID")
	app.SessionManager.Remove(r.Context(), "userName")
	app.setFlash(r, "You have been logged out.")
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (app *Application) ListUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := app.UserRepo.GetUsers(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(users); err != nil {
		app.serverError(w, err)
	}
}

func (app *Application) GetByIDHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	user, err := app.UserRepo.GetUserByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}
		app.serverError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		app.serverError(w, err)
	}
}
