package http

import (
	"encoding/json"
	"net/http"

	"github.com/yourusername/twitter-services-app/services/twitte-service/internal/application"
)

type Handler struct {
	service *application.TweetService
}
type CreateTweetRequest struct {
	Content string `json:"content"`
}

func (h *Handler) CreateTweet(w http.ResponseWriter, r *http.Request) {
	idempotencyKey := r.Header.Get("Idempotency-Key")
	if idempotencyKey == "" {
		http.Error(w, "Missing Idempotency-Key", http.StatusBadRequest)
		return
	}

	var req CreateTweetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	if len(req.Content) == 0 || len(req.Content) > 280 {
		http.Error(w, "Invalid tweet length", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("userId").(string)
	correlationID := r.Context().Value("correlationId").(string)

	tweetID, err := h.service.CreateTweet(
		r.Context(),
		userID,
		req.Content,
		idempotencyKey,
		correlationID,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"tweetId": tweetID,
	})
}
