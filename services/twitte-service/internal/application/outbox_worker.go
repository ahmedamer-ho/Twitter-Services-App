package application

import (
	"context"
	"log"
	"time"

	"github.com/yourusername/twitter-services-app/services/twitte-service/internal/domain"
)

type EventPublisher interface {
	Publish(ctx context.Context, key string, value []byte) error
}

type OutboxWorker struct {
	repo      domain.TweetRepository
	publisher EventPublisher
	interval  time.Duration
}

func NewOutboxWorker(repo domain.TweetRepository, publisher EventPublisher, interval time.Duration) *OutboxWorker {
	return &OutboxWorker{
		repo:      repo,
		publisher: publisher,
		interval:  interval,
	}
}

func (w *OutboxWorker) Start(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.processEvents(ctx)
		}
	}
}

func (w *OutboxWorker) processEvents(ctx context.Context) {
	events, err := w.repo.FetchUnsentEvents(ctx, 10)
	if err != nil {
		log.Printf("Failed to fetch unsent events: %v", err)
		return
	}

	for _, event := range events {
		err := w.publisher.Publish(ctx, event.AggregateID, event.Payload)
		if err != nil {
			log.Printf("Failed to publish event %s: %v", event.ID, err)
			continue
		}

		err = w.repo.MarkEventAsSent(ctx, event.ID)
		if err != nil {
			log.Printf("Failed to mark event %s as sent: %v", event.ID, err)
		}
	}
}
