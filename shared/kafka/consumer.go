package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis"
	"github.com/segmentio/kafka-go"
	"github.com/yourusername/twitter-services-app/shared/observability"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

type HandlerFunc func(ctx context.Context, event Event) error

type Consumer struct {
	reader       *kafka.Reader
	handler      HandlerFunc
	deduplicator *Deduplicator
}
type Deduplicator struct {
	redis *redis.Client
}

func NewConsumer(brokers []string, groupID, topic string, handler HandlerFunc) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: brokers,
			GroupID: groupID,
			Topic:   topic,
		}),
		handler: handler,
	}
}

func (c *Consumer) Start(ctx context.Context) {
	meter := otel.GetMeterProvider().Meter("kafka-consumer")
	lagGauge, _ := meter.Int64ObservableGauge("kafka_consumer_lag", metric.WithUnit("1"), metric.WithDescription("Current lag of the consumer group"))

	// Register callback for lag
	_, _ = meter.RegisterCallback(func(ctx context.Context, o metric.Observer) error {
		stats := c.reader.Stats()
		o.ObserveInt64(lagGauge, stats.Lag, metric.WithAttributes(
			attribute.String("topic", c.reader.Config().Topic),
			attribute.String("group_id", c.reader.Config().GroupID),
		))
		return nil
	}, lagGauge)

	for {
		msg, err := c.reader.FetchMessage(ctx)
		if err != nil {
			log.Println("fetch error:", err)
			continue
		}

		// Extract Trace Context
		ctx := observability.ExtractKafkaHeaders(ctx, &msg)

		tracer := otel.Tracer("kafka-consumer")
		ctx, span := tracer.Start(ctx, "process_message", trace.WithAttributes(
			attribute.String("messaging.system", "kafka"),
			attribute.String("messaging.destination", msg.Topic),
			attribute.String("messaging.kafka.offset", fmt.Sprintf("%d", msg.Offset)),
		))

		var event Event
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Println("invalid event:", err)
			span.RecordError(err)
			span.End()
			continue
		}

		if err := c.handler(ctx, event); err != nil {
			log.Println("handler failed:", err)
			span.RecordError(err)
			span.End()
			continue
		}

		// ✅ manual commit AFTER success
		if err := c.reader.CommitMessages(ctx, msg); err != nil {
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

func (c *Consumer) HandleMessage(msg kafka.Message) {
	var event Event
	if err := json.Unmarshal(msg.Value, &event); err != nil {
		log.Println("invalid event:", err)
		c.reader.CommitMessages(context.Background(), msg)
		return
	}

	// 1. Deduplication
	seen, _ := c.deduplicator.Seen(event.EventID)
	if seen {
		c.reader.CommitMessages(context.Background(), msg)
		return
	}

	// 2. Process
	err := c.handler(context.Background(), event)
	if err == nil {
		c.deduplicator.Mark(event.EventID)
		c.reader.CommitMessages(context.Background(), msg)
		return
	}

	// 3. Retry or DLQ
	c.reader.CommitMessages(context.Background(), msg)
}
