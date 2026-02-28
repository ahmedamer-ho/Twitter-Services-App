package nats

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/yourusername/twitter-services-app/shared/observability"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type HandlerFunc func(ctx context.Context, event Event) error

type Consumer struct {
	js           jetstream.JetStream
	subject      string
	queueGroup   string
	handler      HandlerFunc
	deduplicator *Deduplicator
}

type Deduplicator struct {
	redis *redis.Client
}

func NewConsumer(js jetstream.JetStream, subject, queueGroup string, handler HandlerFunc) *Consumer {
	return &Consumer{
		js:         js,
		subject:    subject,
		queueGroup: queueGroup,
		handler:    handler,
	}
}

func (c *Consumer) Start(ctx context.Context) {
	// Creating a generic stream for the subject if it doesn't exist
	streamName := "stream_" + c.subject // Simple naming convention
	
	streamConfig := jetstream.StreamConfig{
		Name:     streamName,
		Subjects: []string{c.subject},
	}
	
	_, err := c.js.CreateStream(ctx, streamConfig)
	if err != nil {
		// Log and continue, stream might already exist
		log.Printf("Stream creation issue (might already exist): %v", err)
	}

	consumerConfig := jetstream.ConsumerConfig{
		Durable:       c.queueGroup,
		DeliverPolicy: jetstream.DeliverAllPolicy,
		AckPolicy:     jetstream.AckExplicitPolicy,
		FilterSubject: c.subject,
	}

	cons, err := c.js.CreateOrUpdateConsumer(ctx, streamName, consumerConfig)
	if err != nil {
		log.Fatalf("Failed to create/update JetStream consumer: %v", err)
	}

	iter, err := cons.Messages()
	if err != nil {
		log.Fatalf("Failed to start consuming messages: %v", err)
	}

	log.Printf("Started JetStream Consumer on Subject '%s' with Queue Group '%s'", c.subject, c.queueGroup)

	for {
		msg, err := iter.Next()
		if err != nil {
			log.Println("fetch error:", err)
			continue
		}

		// Extract Trace Context
		ctx := observability.ExtractNATSHeaders(ctx, msg.Headers())

		tracer := otel.Tracer("nats-consumer")
		ctx, span := tracer.Start(ctx, "process_message", trace.WithAttributes(
			attribute.String("messaging.system", "nats"),
			attribute.String("messaging.destination", msg.Subject()),
		))

		var event Event
		if err := json.Unmarshal(msg.Data(), &event); err != nil {
			log.Println("invalid event:", err)
			span.RecordError(err)
			span.End()
			msg.Ack() // Ack to prevent redelivery of poisonous messages
			continue
		}

		if err := c.handler(ctx, event); err != nil {
			log.Println("handler failed:", err)
			span.RecordError(err)
			span.End()
			// Don't ack so it's retried
			msg.Nak()
			continue
		}

		// Acknowledge the message AFTER success
		if err := msg.Ack(); err != nil {
			log.Println("commit failed:", err)
		}
		span.End()
	}
}

func (d *Deduplicator) Seen(eventID string) (bool, error) {
	count, err := d.redis.Exists("processed:event:" + eventID).Result()
	if err != nil {
		return false, err
	}
	return count == 1, nil
}

func (d *Deduplicator) Mark(eventID string) error {
	seconds := 10
	return d.redis.Set(
		"processed:event:"+eventID,
		"1",
		time.Duration(seconds)*time.Second,
	).Err()
}
