package observability

import (
	"context"

	"go.opentelemetry.io/otel"
)

// InjectNATSHeaders injects the current context's span into the NATS message headers.
// Call this in the Producer before writing a message.
func InjectNATSHeaders(ctx context.Context, headers map[string][]string) {
	propagator := otel.GetTextMapPropagator()
	propagator.Inject(ctx, &natsHeaderCarrier{headers})
}

// ExtractNATSHeaders extracts the span context (if any) from the NATS message headers.
// Call this in the Consumer before processing a message.
func ExtractNATSHeaders(ctx context.Context, headers map[string][]string) context.Context {
	propagator := otel.GetTextMapPropagator()
	return propagator.Extract(ctx, &natsHeaderCarrier{headers})
}

// natsHeaderCarrier implements propagation.TextMapCarrier for NATS headers (map[string][]string)
type natsHeaderCarrier struct {
	headers map[string][]string
}

func (c *natsHeaderCarrier) Get(key string) string {
	if values, ok := c.headers[key]; ok && len(values) > 0 {
		return values[0]
	}
	return ""
}

func (c *natsHeaderCarrier) Set(key string, value string) {
	c.headers[key] = []string{value}
}

func (c *natsHeaderCarrier) Keys() []string {
	keys := make([]string, 0, len(c.headers))
	for k := range c.headers {
		keys = append(keys, k)
	}
	return keys
}
