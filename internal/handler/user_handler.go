package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

func (app *Application) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	// 1. GET Request: Just show the registration page
	if r.Method == http.MethodGet {
		app.renderTemplate(w, "register.html", nil)
		return
	}

	// 2. POST Request: Process the form
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Extract values from the form
	name := r.PostForm.Get("name")
	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")
	confirm := r.PostForm.Get("confirmPassword")

	// 3. Call your Repository
	_, err = app.UserRepo.CreateUser(r.Context(), name, email, password, confirm, "default_avatar.png")
	if err != nil {
		// Render the page again with the error message
		app.renderTemplate(w, "register.html", err.Error())
		return
	}

	// 4. Redirect to login on success
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (app *Application) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		app.renderTemplate(w, "login.html", nil)
		return
	}

	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")

	// 1. Fetch user by email
	user, err := app.UserRepo.GetUserByEmail(r.Context(), email)
	if err != nil {
		app.renderTemplate(w, "login.html", "Invalid credentials")
		return
	}

	// 2. Compare Bcrypt Hash
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		app.renderTemplate(w, "login.html", "Invalid credentials")
		return
	}

	// 3. Success (Session establishment would go here)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
func (app *Application) ListUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := app.UserRepo.GetUsers(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (app *Application) GetByIDHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	user, err := app.UserRepo.GetUserByID(r.Context(), id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	if user == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
