// Package health exposes liveness and readiness probes.
package health

import (
	"context"
	"net/http"
	"time"

	"github.com/fekom/udt/internal/platform/httpx"
)

type Pinger interface {
	Ping(ctx context.Context) error
}

type status struct {
	Status string `json:"status"`
}

// Routes registers the probes. Liveness never touches the database, so a
// database outage does not make the orchestrator restart a healthy process.
func Routes(mux *http.ServeMux, db Pinger) {
	mux.HandleFunc("GET /healthz", handleLiveness)
	mux.HandleFunc("GET /readyz", handleReadiness(db))
}

func handleLiveness(w http.ResponseWriter, _ *http.Request) {
	httpx.WriteJSON(w, http.StatusOK, status{Status: "ok"})
}

func handleReadiness(db Pinger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		if err := db.Ping(ctx); err != nil {
			httpx.WriteJSON(w, http.StatusServiceUnavailable, status{Status: "database unavailable"})
			return
		}
		httpx.WriteJSON(w, http.StatusOK, status{Status: "ready"})
	}
}
