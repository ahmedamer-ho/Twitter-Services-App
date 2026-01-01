package application

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/twitter-services-app/services/twitte-service/internal/domain"
)

type TweetService struct {
	repo domain.TweetRepository
}
// type TweetService interface {
// 	CreateTweet(
// 		ctx context.Context,
// 		authorID string,
// 		content string,
// 		idempotencyKey string,
// 		correlationID string,
// 	) (string, error)
// }


func NewTweetService(repo domain.TweetRepository) *TweetService {
	return &TweetService{repo: repo}
}
// func (s *TweetService) CreateTweet(ctx context.Context, authorID, content string) error {
// 	t := domain.Tweet{
// 		ID:        uuid.NewString(),
// 		AuthorID:  authorID,
// 		Content:   content,
// 		CreatedAt: time.Now().UTC(),
// 	}

// 	return s.repo.Insert(ctx, t)
// }


func (s *TweetService) CreateTweet(
	ctx context.Context,
	authorID, content, idempotencyKey, correlationID string,
) (string, error) {

	// 1️⃣ Idempotency check
	existing, err := s.repo.FindByIdempotencyKey(ctx, idempotencyKey)
	if err == nil {
		return existing.ID, nil
	}

	tweetID := uuid.NewString()

	tweet := domain.Tweet{
		ID:             tweetID,
		AuthorID:       authorID,
		Content:        content,
		CreatedAt:      time.Now(),
		IdempotencyKey: idempotencyKey,
	}

	eventPayload, _ := json.Marshal(tweet)

	event := domain.OutboxEvent{
		ID:            uuid.NewString(),
		EventType:     "TweetCreated",
		AggregateID:   tweetID,
		Payload:       eventPayload,
		CorrelationID: correlationID,
		CreatedAt:     time.Now(),
		Sent:          false,
	}

	// 2️⃣ MongoDB transaction
	return tweetID, s.repo.WithTransaction(ctx, func(tx domain.TweetRepository) error {
		if err := tx.Insert(ctx, tweet); err != nil {
			return err
		}
		return tx.InsertOutbox(ctx, event)
	})
}
