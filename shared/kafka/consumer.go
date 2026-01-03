package kafka

import (
	"context"
	"encoding/json"
	"log"
    "time"
	"github.com/go-redis/redis"
	"github.com/segmentio/kafka-go"
)

type HandlerFunc func(ctx context.Context, event Event) error

type Consumer struct {
	reader  *kafka.Reader
	handler HandlerFunc
	deduplicator *Deduplicator
}
type Deduplicator struct {
	redis *redis.Client
}
func NewConsumer(brokers []string, groupID, topic string, handler HandlerFunc) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:  brokers,
			GroupID: groupID,
			Topic:   topic,
		}),
		handler: handler,
	}
}

func (c *Consumer) Start(ctx context.Context) {
	for {
		msg, err := c.reader.FetchMessage(ctx)
		if err != nil {
			log.Println("fetch error:", err)
			continue
		}

		var event Event
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Println("invalid event:", err)
			continue
		}

		if err := c.handler(ctx, event); err != nil {
			log.Println("handler failed:", err)
			continue
		}

		// ✅ manual commit AFTER success
		if err := c.reader.CommitMessages(ctx, msg); err != nil {
			log.Println("commit failed:", err)
		}
	}
}


func (d *Deduplicator) Seen(eventID string) (bool, error) {
	count, err := d.redis.Exists("processed:event:"+eventID).Result()
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
