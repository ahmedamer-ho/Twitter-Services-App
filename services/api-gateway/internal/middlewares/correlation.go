package middlewares

import (
	"context"
	"log"
	"net/http"

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

		log.Printf("Request received correlation_id=%s route=%s", correlationID, r.URL.Path)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
