package main

import (
	"net/http"

	"github.com/fekom/udt/internal/health"
	"github.com/fekom/udt/pkg/database"
)

// SetupRoutes registers all application routes grouped by context.
// To add a new module, create its XxxRoutes function and call it here.
func SetupRoutes(mux *http.ServeMux, db *database.DB) {
	health.HealthRoutes(mux, db)
	// auth.AuthRoutes(mux, db)
	// media.MediaRoutes(mux, db)
	// room.RoomRoutes(mux, db)
}
