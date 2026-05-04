package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/adocoder12/webDevelopment/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

func (app *Application) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		app.renderTemplate(w, "register.html", nil)
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
		app.renderTemplate(w, "register.html", err.Error())
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (app *Application) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		app.renderTemplate(w, "login.html", nil)
		return
	}

	if err := r.ParseForm(); err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	user, err := app.UserRepo.GetUserByEmail(r.Context(), r.PostForm.Get("email"))
	if err != nil {
		app.renderTemplate(w, "login.html", "Invalid email or password")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(r.PostForm.Get("password"))); err != nil {
		app.renderTemplate(w, "login.html", "Invalid email or password")
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
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
		// FIX: use sentinel error to return 404 instead of 500 for missing users
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
