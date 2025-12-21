package kafka

import "time"

type Event struct {
	EventID       string      `json:"eventId"`
	EventType     string      `json:"eventType"`
	AggregateID   string      `json:"aggregateId"`
	Timestamp     time.Time   `json:"timestamp"`
	CorrelationID string      `json:"correlationId"`
	Payload       interface{} `json:"payload"`
}
