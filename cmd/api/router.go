package main

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/fekom/udt/internal/health"
	"github.com/fekom/udt/internal/identity"
	"github.com/fekom/udt/internal/platform/httpx"
)

// newRouter is the single place where modules get wired. Adding a module
// means one Routes call here and nothing else in main.
func newRouter(pool *pgxpool.Pool) http.Handler {
	mux := http.NewServeMux()

	health.Routes(mux, pool)
	identity.Routes(mux, pool)
	// family.Routes(mux, pool)
	// media.Routes(mux, pool, storage)
	// room.Routes(mux, pool)

	return httpx.Chain(mux, httpx.Recover, httpx.LogRequests)
}
