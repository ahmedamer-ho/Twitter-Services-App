package domain

import (

	"time"
)
type OutboxEvent struct {
	ID            string    `bson:"_id"`
	EventType     string    `bson:"eventType"`
	AggregateID   string    `bson:"aggregateId"`
	Payload       []byte    `bson:"payload"`
	CorrelationID string    `bson:"correlationId"`
	CreatedAt     time.Time `bson:"createdAt"`
	Sent          bool      `bson:"sent"`
}
