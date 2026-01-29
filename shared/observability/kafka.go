package observability

import (
	"context"

	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel"
)

// InjectKafkaHeaders injects the current context's span into the Kafka message headers.
// Call this in the Producer before writing a message.
func InjectKafkaHeaders(ctx context.Context, msg *kafka.Message) {
	propagator := otel.GetTextMapPropagator()
	propagator.Inject(ctx, &kafkaHeaderCarrier{msg})
}

// ExtractKafkaHeaders extracts the span context (if any) from the Kafka message headers.
// Call this in the Consumer before processing a message.
func ExtractKafkaHeaders(ctx context.Context, msg *kafka.Message) context.Context {
	propagator := otel.GetTextMapPropagator()
	return propagator.Extract(ctx, &kafkaHeaderCarrier{msg})
}

// kafkaHeaderCarrier implements propagation.TextMapCarrier for kafka.Message
type kafkaHeaderCarrier struct {
	msg *kafka.Message
}

func (c *kafkaHeaderCarrier) Get(key string) string {
	for _, h := range c.msg.Headers {
		if h.Key == key {
			return string(h.Value)
		}
	}
	return ""
}

func (c *kafkaHeaderCarrier) Set(key string, value string) {
	// Remove existing key if present to avoid duplicates
	newHeaders := make([]kafka.Header, 0, len(c.msg.Headers)+1)
	for _, h := range c.msg.Headers {
		if h.Key != key {
			newHeaders = append(newHeaders, h)
		}
	}
	newHeaders = append(newHeaders, kafka.Header{
		Key:   key,
		Value: []byte(value),
	})
	c.msg.Headers = newHeaders
}

func (c *kafkaHeaderCarrier) Keys() []string {
	keys := make([]string, len(c.msg.Headers))
	for i, h := range c.msg.Headers {
		keys[i] = h.Key
	}
	return keys
}
