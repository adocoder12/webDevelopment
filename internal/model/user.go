package model

import "time"

type User struct {
	Id        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Profile   Profile   `json:"profile"`
}

type Profile struct {
	Avatar   string `json:"avatar"`
	Active   bool   `json:"active"`
	Verified bool   `json:"verified"`
}
