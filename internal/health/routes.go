package health

import (
	"net/http"

	"github.com/fekom/udt/pkg/database"
)

// HealthRoutes registers all healthcheck routes on the given mux.
// It wires the repository and handler internally — callers only need the mux and db.
func HealthRoutes(mux *http.ServeMux, db *database.DB) {
	repo := NewRepository(db.Collection("healthchecks"))
	h := Handler{repo: repo}

	mux.HandleFunc("POST /healthcheck", h.Healthcheck)
}
