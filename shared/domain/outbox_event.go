package domain

import "time"

type OutboxEvent struct {
	ID            string
	EventType     string
	AggregateID   string
	Payload       []byte
	CorrelationID string
	CreatedAt     time.Time
	SentAt        *time.Time
}
