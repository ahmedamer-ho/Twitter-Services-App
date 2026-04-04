package nats

import "time"

type Event struct {
	ID            string      `json:"id"`
	Type          string      `json:"type"`
	Source        string      `json:"source"`
	Timestamp     time.Time   `json:"timestamp"`
	Data          interface{} `json:"data"`
	CorrelationID string      `json:"correlationId,omitempty"`
	AggregateID   string      `json:"aggregateId,omitempty"`
	RetryCount    int         `json:"retryCount,omitempty"`
}
