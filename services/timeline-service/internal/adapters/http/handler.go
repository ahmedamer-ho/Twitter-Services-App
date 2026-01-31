package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/yourusername/twitter-services-app/services/timeline-service/internal/application"
)

type Handler struct {
	service *application.TimelineService
}

func NewHandler(service *application.TimelineService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetTimeline(w http.ResponseWriter, r *http.Request) {
	// Pattern: /timeline/:userId
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "User ID required", http.StatusBadRequest)
		return
	}
	userID := parts[2]

	timeline, err := h.service.GetTimeline(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(timeline)
}
