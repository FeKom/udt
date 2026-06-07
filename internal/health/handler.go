package health

import (
	"log"
	"net"
	"net/http"
)

// Handler holds the dependencies needed by healthcheck HTTP handlers.
type Handler struct {
	repo *Repository
}

// Healthcheck handles POST /healthcheck.
// It records the caller IP, a timestamp, and a "health" message in MongoDB.
func (h *Handler) Healthcheck(w http.ResponseWriter, r *http.Request) {
	ip, err := GetRemoteIP(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := h.repo.Save(r.Context(), ip); err != nil {
		log.Println("failed to save healthcheck", err)
		http.Error(w, "failed to save healthcheck", http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// GetRemoteIP extracts just the IP address from r.RemoteAddr ("ip:port").
func GetRemoteIP(r *http.Request) (string, error) {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		log.Println("failed to split host port", err)
		return "", err
	}
	return ip, nil
}
