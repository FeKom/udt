package identity

import (
	"context"
	"errors"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"

	"github.com/fekom/udt/internal/platform/postgres"
)

// Repository functions are unexported on purpose: SQL stays inside the
// module, and the service is the only door in.

var userColumns = []string{"id", "name", "email", "created_at", "updated_at"}

func insertUser(ctx context.Context, db postgres.DBTX, in NewUser) (User, error) {
	query, args, err := postgres.SQL.Insert("users").
		Columns("name", "email").
		Values(in.Name, in.Email).
		Suffix("RETURNING " + strings.Join(userColumns, ", ")).
		ToSql()
	if err != nil {
		return User{}, err
	}

	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		return User{}, err
	}
	u, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[User])
	if postgres.IsUniqueViolation(err) {
		return User{}, ErrEmailTaken
	}
	return u, err
}

func findUserByID(ctx context.Context, db postgres.DBTX, id string) (User, error) {
	query, args, err := postgres.SQL.Select(userColumns...).
		From("users").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return User{}, err
	}

	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		return User{}, err
	}
	u, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[User])
	if errors.Is(err, pgx.ErrNoRows) {
		return User{}, ErrNotFound
	}
	return u, err
}

func listUsers(ctx context.Context, db postgres.DBTX, f UserFilter) ([]User, error) {
	q := postgres.SQL.Select(userColumns...).
		From("users").
		OrderBy("created_at DESC").
		Limit(f.Limit)

	if f.Name != "" {
		q = q.Where(sq.ILike{"name": "%" + f.Name + "%"})
	}
	if f.Email != "" {
		q = q.Where(sq.Eq{"email": f.Email})
	}

	query, args, err := q.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[User])
}
