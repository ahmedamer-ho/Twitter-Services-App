package domain

import "context"

type TwiteRepository interface {
	Insert(ctx context.Context, t Twite) error
	FindByID(ctx context.Context, id string) (*Twite, error)
	FindByAuthor(ctx context.Context, authorID string, limit int) ([]Twite, error)
	SoftDelete(ctx context.Context, id string) error
}
