package handler

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}

func (app *Application) Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rw, r)
		app.InfoLog.Printf("%s - %s %s %s — %d",
			r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI(), rw.status)
	})
}

func (app *Application) RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")

				// Log the error and the stack trace.
				app.ErrorLog.Printf("%s\n%s", err, debug.Stack())

				// Send a 500 Internal Server Error response.
				app.serverError(w, fmt.Errorf("%v", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (app *Application) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := app.SessionManager.GetInt64(r.Context(), "userID")
		if userID == 0 {
			// Set a flash message then redirect to login
			app.SessionManager.Put(r.Context(), "flash", "You must be logged in to view this page")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})

}
