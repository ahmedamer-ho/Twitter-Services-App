package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
	"github.com/yourusername/twitter-services-app/shared/observability"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type Producer interface {
	Publish(ctx context.Context, topic string, event Event) error
}

type kafkaProducer struct {
	writer *kafka.Writer
}

func NewProducer(brokers []string) Producer {
	return &kafkaProducer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func (p *kafkaProducer) Publish(ctx context.Context, topic string, event Event) error {
	// Start Span
	tracer := otel.Tracer("kafka-producer")
	ctx, span := tracer.Start(ctx, fmt.Sprintf("publish %s", topic), trace.WithAttributes(
		attribute.String("messaging.system", "kafka"),
		attribute.String("messaging.destination", topic),
		attribute.String("messaging.destination_kind", "topic"),
	))
	defer span.End()

	value, err := json.Marshal(event)
	if err != nil {
		return err
	}

	msg := kafka.Message{
		Topic: topic,
		Key:   []byte(event.AggregateID),
		Value: value,
		Headers: []kafka.Header{
			{Key: HeaderCorrelationID, Value: []byte(event.CorrelationID)},
			{Key: EventTypeHeader, Value: []byte(event.EventType)},
		},
	}

	// Propagate Trace Context
	observability.InjectKafkaHeaders(ctx, &msg)

	return p.writer.WriteMessages(ctx, msg)
}
