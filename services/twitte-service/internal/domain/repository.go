package domain

import "context"

type TweetRepository interface {
	Insert(ctx context.Context, t Tweet) error
	FindByID(ctx context.Context, id string) (*Tweet, error)
	FindByAuthor(ctx context.Context, authorID string, limit int) ([]Tweet, error)
	SoftDelete(ctx context.Context, id string) error
	FindByIdempotencyKey(ctx context.Context, idempotencyKey string) (*Tweet, error)
	WithTransaction(ctx context.Context, fn func(TweetRepository) error) error
	InsertOutbox(ctx context.Context, event OutboxEvent) error
}
