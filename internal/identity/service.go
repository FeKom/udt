package identity

import (
	"context"
	"fmt"
	"net/mail"
	"strings"

	"github.com/fekom/udt/internal/platform/postgres"
)

const (
	defaultListLimit = 50
	maxListLimit     = 200
)

// CreateUser validates and normalizes the input before persisting it.
func CreateUser(ctx context.Context, db postgres.DBTX, in NewUser) (User, error) {
	in, err := normalizeNewUser(in)
	if err != nil {
		return User{}, err
	}
	return insertUser(ctx, db, in)
}

func GetUser(ctx context.Context, db postgres.DBTX, id string) (User, error) {
	return findUserByID(ctx, db, id)
}

func ListUsers(ctx context.Context, db postgres.DBTX, f UserFilter) ([]User, error) {
	f.Email = strings.ToLower(strings.TrimSpace(f.Email))
	f.Limit = clampLimit(f.Limit)
	return listUsers(ctx, db, f)
}

// normalizeNewUser is pure: no I/O, so it is trivial to unit test.
func normalizeNewUser(in NewUser) (NewUser, error) {
	in.Name = strings.TrimSpace(in.Name)
	// Lowercasing here keeps the UNIQUE constraint effectively case-insensitive.
	in.Email = strings.ToLower(strings.TrimSpace(in.Email))

	if in.Name == "" {
		return NewUser{}, fmt.Errorf("%w: name is required", ErrInvalidInput)
	}
	if _, err := mail.ParseAddress(in.Email); err != nil {
		return NewUser{}, fmt.Errorf("%w: email is invalid", ErrInvalidInput)
	}
	return in, nil
}

func clampLimit(limit uint64) uint64 {
	switch {
	case limit == 0:
		return defaultListLimit
	case limit > maxListLimit:
		return maxListLimit
	default:
		return limit
	}
}
