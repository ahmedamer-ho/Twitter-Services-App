package domain

import "context"

type OutboxRepository interface {
	Insert(ctx context.Context, event OutboxEvent) error
	FetchPending(ctx context.Context, limit int) ([]OutboxEvent, error)
	MarkSent(ctx context.Context, eventID string) error
}
