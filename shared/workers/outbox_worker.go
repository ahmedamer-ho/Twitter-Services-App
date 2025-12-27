package workers

import (
	"context"
	"encoding/json"
	"time"

	"github.com/yourusername/twitter-services-app/shared/adapters/postgres"
	"github.com/yourusername/twitter-services-app/shared/kafka"
)

type OutboxWorker struct {
	repo     *postgres.OutboxRepository
	producer kafka.Producer
}
func (w *OutboxWorker) Run(ctx context.Context) {
	ticker := time.NewTicker(2 * time.Second)

	for {
		select {
		case <-ticker.C:
			w.process(ctx)
		case <-ctx.Done():
			return
		}
	}
}
func (w *OutboxWorker) process(ctx context.Context) {
	events, err := w.repo.FetchPending(ctx, 10)
	if err != nil {
		return
	}

	for _, e := range events {
		var payload map[string]any
		_ = json.Unmarshal(e.Payload, &payload)

		event := kafka.Event{
			EventID:       e.ID,
			EventType:     e.EventType,
			AggregateID:   e.AggregateID,
			Timestamp:     e.CreatedAt,
			CorrelationID: e.CorrelationID,
			Payload:       payload,
		}

		if err := w.producer.Publish(ctx, "user-events", event); err != nil {
			continue // retry later
		}

		_ = w.repo.MarkSent(ctx, e.ID)
	}
}
