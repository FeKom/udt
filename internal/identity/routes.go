package identity

import (
	"net/http"

	"github.com/fekom/udt/internal/platform/postgres"
)

func Routes(mux *http.ServeMux, db postgres.DBTX) {
	mux.HandleFunc("POST /users", handleCreateUser(db))
	mux.HandleFunc("GET /users", handleListUsers(db))
	mux.HandleFunc("GET /users/{id}", handleGetUser(db))
}
