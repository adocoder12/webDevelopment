package repository

import (
	"context"
	"database/sql"

	"github.com/adocoder12/webDevelopment/internal/model"
)

var UserSchema = `
CREATE TABLE IF NOT EXISTS Users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    hashedPassword TEXT NOT NULL,
    createdAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
);`

var ProfileSchema = `
CREATE TABLE IF NOT EXISTS Profile (
    user_id INTEGER PRIMARY KEY REFERENCES Users(id) ON DELETE CASCADE,
    avatar TEXT NOT NULL,
    active BOOLEAN NOT NULL DEFAULT 1,
    verified BOOLEAN NOT NULL DEFAULT 0
);`

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

type UserRepository interface {
	CreateUser(ctx context.Context, name, email, password, confirmPassword, avatar string) (*model.User, error)
	GetUsers(ctx context.Context) ([]*model.User, error)
	GetUserByID(ctx context.Context, id int64) (*model.User, error) // Notice: I removed db DBTX here! (See note below)
}
