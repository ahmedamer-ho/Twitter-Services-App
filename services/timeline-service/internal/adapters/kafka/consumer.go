package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
	"github.com/yourusername/twitter-services-app/services/timeline-service/internal/application"
	"github.com/yourusername/twitter-services-app/services/timeline-service/internal/domain"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
)

type Consumer struct {
	reader  *kafka.Reader
	service *application.TimelineService
}

func NewConsumer(brokers []string, topic, groupID string, service *application.TimelineService) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:  brokers,
			Topic:    topic,
			GroupID:  groupID,
			MinBytes: 10e3, // 10KB
			MaxBytes: 10e6, // 10MB
		}),
		service: service,
	}
}

func (c *Consumer) Start(ctx context.Context) {
	log.Println("Timeline Service Kafka Consumer started...")
	for {
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			log.Printf("Error reading message: %v", err)
			continue
		}

		// Extract Trace Context
		carrier := propagation.MapCarrier{}
		for _, h := range m.Headers {
			carrier[h.Key] = string(h.Value)
		}
		ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)

		// Start Span
		tracer := otel.Tracer("timeline-service")
		ctx, span := tracer.Start(ctx, "kafka.consume")
		span.SetAttributes(
			attribute.String("kafka.topic", m.Topic),
			attribute.Int64("kafka.offset", m.Offset),
		)

		var tweet domain.TimelineTweet
		if err := json.Unmarshal(m.Value, &tweet); err != nil {
			log.Printf("Error unmarshaling message: %v", err)
			span.RecordError(err)
			span.End()
			continue
		}

		err = c.service.AddTweetToTimeline(ctx, tweet)
		if err != nil {
			log.Printf("Error adding tweet to timeline: %v", err)
			span.RecordError(err)
		}
		span.End()
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
