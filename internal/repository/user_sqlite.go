package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/adocoder12/webDevelopment/internal/model"
	"golang.org/x/crypto/bcrypt"
)

type SqlUserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &SqlUserRepository{db: db}
}

func (r *SqlUserRepository) CreateUser(ctx context.Context, name, email, password, confirmPassword, avatar string) (*model.User, error) {
	if password != confirmPassword {
		return nil, errors.New("passwords do not match")
	}

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	res, err := tx.ExecContext(ctx,
		"INSERT INTO Users (name, email, hashedPassword) VALUES (?, ?, ?)",
		name, email, string(hashedBytes),
	)
	if err != nil {
		return nil, fmt.Errorf("user insert: %w", err)
	}

	// FIX: error from LastInsertId is now checked
	userID, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("last insert id: %w", err)
	}

	if _, err = tx.ExecContext(ctx,
		"INSERT INTO Profile (user_id, avatar) VALUES (?, ?)",
		userID, avatar,
	); err != nil {
		return nil, fmt.Errorf("profile insert: %w", err)
	}

	user, err := getUserByIDTx(ctx, tx, userID)
	if err != nil {
		return nil, err
	}

	return user, tx.Commit()
}

func (r *SqlUserRepository) GetUsers(ctx context.Context) ([]*model.User, error) {
	query := `
		SELECT u.id, u.name, u.email, u.hashedPassword, u.createdAt, u.updatedAt,
		       p.avatar, p.active, p.verified
		FROM   Users u
		LEFT JOIN Profile p ON u.id = p.user_id`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*model.User
	for rows.Next() {
		var u model.User
		if err := rows.Scan(
			&u.Id, &u.Name, &u.Email, &u.Password,
			&u.CreatedAt, &u.UpdatedAt,
			&u.Profile.Avatar, &u.Profile.Active, &u.Profile.Verified,
		); err != nil {
			return nil, err
		}
		users = append(users, &u)
	}
	return users, rows.Err()
}

func (r *SqlUserRepository) GetUserByID(ctx context.Context, id int64) (*model.User, error) {
	return getUserByIDTx(ctx, r.db, id)
}

// getUserByIDTx is an unexported helper that accepts any DBTX so it works
// both inside a transaction (CreateUser) and against the plain *sql.DB.
func getUserByIDTx(ctx context.Context, db DBTX, id int64) (*model.User, error) {
	query := `
		SELECT u.id, u.name, u.email, u.hashedPassword, u.createdAt, u.updatedAt,
		       p.avatar, p.active, p.verified
		FROM   Users u
		LEFT JOIN Profile p ON u.id = p.user_id
		WHERE  u.id = ?`

	var u model.User
	err := db.QueryRowContext(ctx, query, id).Scan(
		&u.Id, &u.Name, &u.Email, &u.Password,
		&u.CreatedAt, &u.UpdatedAt,
		&u.Profile.Avatar, &u.Profile.Active, &u.Profile.Verified,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user %d: %w", id, ErrNotFound)
		}
		return nil, err
	}
	return &u, nil
}

func (r *SqlUserRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `
		SELECT u.id, u.name, u.email, u.hashedPassword, u.createdAt, u.updatedAt,
		       p.avatar, p.active, p.verified
		FROM   Users u
		LEFT JOIN Profile p ON u.id = p.user_id
		WHERE  u.email = ?`

	var u model.User
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&u.Id, &u.Name, &u.Email, &u.Password,
		&u.CreatedAt, &u.UpdatedAt,
		&u.Profile.Avatar, &u.Profile.Active, &u.Profile.Verified,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user %s: %w", email, ErrNotFound)
		}
		return nil, err
	}
	return &u, nil
}
