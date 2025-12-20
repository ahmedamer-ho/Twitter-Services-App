package middlewares

import (
	"context"
	"net/http"

	"github.com/Twitter-Services-App/user-service/internal/logger"
	"github.com/google/uuid"
)

type ctxKey string

const CorrelationIDKey ctxKey = "correlation_id"

func CorrelationID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		correlationID := r.Header.Get("X-Correlation-ID")
		if correlationID == "" {
			correlationID = uuid.New().String()
		}

		ctx := context.WithValue(r.Context(), CorrelationIDKey, correlationID)
		w.Header().Set("X-Correlation-ID", correlationID)

		logger.Log.Info().
			Str("correlation_id", correlationID).
			Str("route", r.URL.Path).
			Msg("Request received")

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
