// Package identity owns users. Other modules use the exported functions in
// service.go; they never query the users table directly.
package identity

import (
	"errors"
	"time"
)

type User struct {
	ID        string    `db:"id"         json:"id"`
	Name      string    `db:"name"       json:"name"`
	Email     string    `db:"email"      json:"email"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type NewUser struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UserFilter struct {
	Name  string
	Email string
	Limit uint64
}

var (
	ErrNotFound     = errors.New("identity: user not found")
	ErrEmailTaken   = errors.New("identity: email already registered")
	ErrInvalidInput = errors.New("identity: invalid input")
)
