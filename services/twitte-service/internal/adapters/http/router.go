package http

import (
	"encoding/json"
	"net/http"
)

// NewRouter wires HTTP routes
func NewRouter(h *Handler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/health/live", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("alive"))
	})

	mux.HandleFunc("/health/ready", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ready"))
	})

	mux.HandleFunc("/tweets", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			h.CreateTweet(w, r)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Tweet endpoint placeholder",
		})
	})

	return mux
}
