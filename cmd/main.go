package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/alexedwards/scs/sqlite3store"
	"github.com/alexedwards/scs/v2"

	"github.com/adocoder12/webDevelopment/internal/handler"
	"github.com/adocoder12/webDevelopment/internal/repository"
)

func main() {
	// 1. Structured Logging
	// os.Stdout for general info, os.Stderr for errors.
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// 2. Database Initialization
	// SQLite with foreign keys enabled is vital for data integrity.
	db, err := sql.Open("sqlite3", "./app.db?_foreign_keys=on")
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	// 3. Fail-Fast Database Ping
	// Never start a server if the database is unreachable.
	ctxPing, cancelPing := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelPing()

	if err := db.PingContext(ctxPing); err != nil {
		errorLog.Fatal(fmt.Errorf("database unreachable: %w", err))
	}
	infoLog.Println("Database connection established!")

	// 4. Schema Migration/Verification
	ctxMigrate, cancelMigrate := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelMigrate()

	if _, err := db.ExecContext(ctxMigrate, repository.UserSchema); err != nil {
		errorLog.Fatal(fmt.Errorf("failed to create users table: %w", err))
	}
	if _, err := db.ExecContext(ctxMigrate, repository.ProfileSchema); err != nil {
		errorLog.Fatal(fmt.Errorf("failed to create profile table: %w", err))
	}
	if _, err := db.ExecContext(ctxMigrate, repository.SessionSchema); err != nil {
		errorLog.Fatal(fmt.Errorf("failed to create sessions table: %w", err))
	}

	// 5. Specialist Initialization
	// We pass 'true' for dev mode so you don't have to restart for HTML edits.
	renderer, err := handler.NewTemplateRenderer("./templates/html", true)
	if err != nil {
		errorLog.Fatal(err)
	}

	// 6. Session Manager
	sessionStore := sqlite3store.New(db)
	defer sessionStore.StopCleanup()
	sessionManager := scs.New()
	sessionManager.Store = sessionStore
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = false // set true in production (HTTPS only)
	sessionManager.Cookie.SameSite = http.SameSiteLaxMode
	// after creating sessionStore

	// 6. Dependency Injection
	userRepo := repository.NewUserRepository(db)
	app := handler.NewApplication(errorLog, infoLog, userRepo, renderer, sessionManager)

	// 7. Configured HTTP Server
	// Using the http.Server struct allows for timeouts that prevent DDoS/leaks.
	srv := &http.Server{
		Addr:         ":8280",
		ErrorLog:     errorLog,
		Handler:      app.SetupRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Println("Server running on :8280")
	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		errorLog.Fatal(err)
	}
}
