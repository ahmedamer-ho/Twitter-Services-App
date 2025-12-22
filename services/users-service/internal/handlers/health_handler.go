package handlers

import "net/http"

type HealthHandler struct{}

func (h *HealthHandler) Live(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (h *HealthHandler) Ready(w http.ResponseWriter, r *http.Request) {
	// later: DB, Keycloak ping
	w.WriteHeader(http.StatusOK)
}
