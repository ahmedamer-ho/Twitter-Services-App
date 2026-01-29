package observability

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

type contextKey string

const CorrelationIDKey contextKey = "correlation_id"

// Middleware checks for X-Correlation-ID header.
// If present, it puts it in context.
// If not, it checks if there is an OTel TraceID and uses that.
// Finally, generates a new one.
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		correlationID := r.Header.Get("X-Correlation-ID")
		
		// If no correlation ID, try to get from OTel Span
		if correlationID == "" {
			spanCtx := trace.SpanContextFromContext(r.Context())
			if spanCtx.HasTraceID() {
				correlationID = spanCtx.TraceID().String()
			}
		}

		// Fallback
		if correlationID == "" {
			correlationID = uuid.New().String()
		}

		ctx := context.WithValue(r.Context(), CorrelationIDKey, correlationID)
		w.Header().Set("X-Correlation-ID", correlationID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func FromContext(ctx context.Context) string {
	if v := ctx.Value(CorrelationIDKey); v != nil {
		return v.(string)
	}
	return ""
}
