package identity

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/fekom/udt/internal/platform/httpx"
	"github.com/fekom/udt/internal/platform/postgres"
)

// Handlers are closures over their dependencies: no struct, no hidden state.

func handleCreateUser(db postgres.DBTX) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var in NewUser
		if err := httpx.DecodeJSON(w, r, &in); err != nil {
			httpx.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}

		u, err := CreateUser(r.Context(), db, in)
		if err != nil {
			writeServiceError(w, r, err)
			return
		}
		httpx.WriteJSON(w, http.StatusCreated, u)
	}
}

func handleGetUser(db postgres.DBTX) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		// Checked before the query: pgx rejects a malformed uuid with an
		// encode error that would otherwise surface as a 500.
		if !isUUID(id) {
			httpx.WriteError(w, http.StatusNotFound, "user not found")
			return
		}

		u, err := GetUser(r.Context(), db, id)
		if err != nil {
			writeServiceError(w, r, err)
			return
		}
		httpx.WriteJSON(w, http.StatusOK, u)
	}
}

func handleListUsers(db postgres.DBTX) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		f := UserFilter{Name: q.Get("name"), Email: q.Get("email")}

		if raw := q.Get("limit"); raw != "" {
			limit, err := strconv.ParseUint(raw, 10, 64)
			if err != nil {
				httpx.WriteError(w, http.StatusBadRequest, "limit must be a positive integer")
				return
			}
			f.Limit = limit
		}

		users, err := ListUsers(r.Context(), db, f)
		if err != nil {
			writeServiceError(w, r, err)
			return
		}
		httpx.WriteJSON(w, http.StatusOK, users)
	}
}

func writeServiceError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, ErrInvalidInput):
		httpx.WriteError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, ErrNotFound):
		httpx.WriteError(w, http.StatusNotFound, "user not found")
	case errors.Is(err, ErrEmailTaken):
		httpx.WriteError(w, http.StatusConflict, "email already registered")
	default:
		httpx.WriteInternalError(w, r, err)
	}
}

func isUUID(s string) bool {
	if len(s) != 36 {
		return false
	}
	for i, c := range s {
		switch i {
		case 8, 13, 18, 23:
			if c != '-' {
				return false
			}
		default:
			isHex := (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')
			if !isHex {
				return false
			}
		}
	}
	return true
}
