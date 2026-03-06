package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/yourusername/twitter-services-app/services/twitte-service/internal/application"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

type Handler struct {
	service *application.TweetService
}

func NewHandler(service *application.TweetService) *Handler {
	return &Handler{service: service}
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

	userID, _ := r.Context().Value("userId").(string)
	if userID == "" {
		userID = "anonymous"
	}
	correlationID, _ := r.Context().Value("correlationId").(string)
	if correlationID == "" {
		correlationID = "manual-test"
	}

	// Observability: Tweet Latency
	meter := otel.GetMeterProvider().Meter("twitte-service")
	histogram, _ := meter.Float64Histogram("tweet_creation_latency", metric.WithUnit("s"))
	start := time.Now()
	defer func() {
		histogram.Record(r.Context(), time.Since(start).Seconds(), metric.WithAttributes(attribute.String("status", "ok")))
	}()

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

func (h *Handler) SearchTweets(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	hashtag := r.URL.Query().Get("hashtag")
	mention := r.URL.Query().Get("mention")
	user := r.URL.Query().Get("user")
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	sort := r.URL.Query().Get("sort")

	page := 1
	if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
		page = p
	}

	limit := 20
	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
		limit = l
	}

	// Observability: Search Latency
	meter := otel.GetMeterProvider().Meter("twitte-service")
	histogram, _ := meter.Float64Histogram("tweet_search_latency", metric.WithUnit("s"))
	start := time.Now()
	defer func() {
		histogram.Record(r.Context(), time.Since(start).Seconds(), metric.WithAttributes(attribute.String("status", "ok")))
	}()

	tweets, err := h.service.SearchTweets(r.Context(), query, hashtag, mention, user, page, limit, sort)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"tweets": tweets,
		"page":   page,
		"limit":  limit,
	})
}
