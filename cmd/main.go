package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"

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
	db, err := sql.Open("sqlite3", "app.db?_foreign_keys=on")
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
	ctx := context.Background()
	if _, err := db.ExecContext(ctx, repository.UserSchema); err != nil {
		errorLog.Fatal(fmt.Errorf("failed to create users table: %w", err))
	}
	if _, err := db.ExecContext(ctx, repository.ProfileSchema); err != nil {
		errorLog.Fatal(fmt.Errorf("failed to create profile table: %w", err))
	}

	// 5. Specialist Initialization
	// We pass 'true' for dev mode so you don't have to restart for HTML edits.
	renderer, err := handler.NewTemplateRenderer("./templates/html", true)
	if err != nil {
		errorLog.Fatal(err)
	}

	// 6. Dependency Injection
	userRepo := repository.NewUserRepository(db)
	app := handler.NewApplication(errorLog, infoLog, userRepo, renderer)

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
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}
