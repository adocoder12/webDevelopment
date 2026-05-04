package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func (app *Application) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	// In a real app, you'd decode r.Body into a struct here
	user, err := app.UserRepo.CreateUser(r.Context(), "Jon Snow", "jon@wall.com", "ghost123", "ghost123", "wolf.png")
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
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
