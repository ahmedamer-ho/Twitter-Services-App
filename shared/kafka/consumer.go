package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
)

type HandlerFunc func(ctx context.Context, event Event) error

type Consumer struct {
	reader  *kafka.Reader
	handler HandlerFunc
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
