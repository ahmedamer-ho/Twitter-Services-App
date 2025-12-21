package kafka

import (
	"context"
	"encoding/json"

	"github.com/segmentio/kafka-go"
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

	return p.writer.WriteMessages(ctx, msg)
}
