package application

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/yourusername/twitter-services-app/services/search-indexer/internal/domain"
	"github.com/yourusername/twitter-services-app/services/search-indexer/internal/infrastructure"
	"github.com/yourusername/twitter-services-app/shared/nats"
)

type IndexTweetHandler struct {
	esClient *infrastructure.ESClient
}

func NewIndexTweetHandler(es *infrastructure.ESClient) *IndexTweetHandler {
	return &IndexTweetHandler{esClient: es}
}

func (h *IndexTweetHandler) HandleTweetEvent(ctx context.Context, event nats.Event) error {
	log.Printf("Received Event: %s (AggregateID: %s)", event.Type, event.AggregateID)

	switch event.Type {
	case "TweetCreated", "TweetUpdated":
		return h.handleIndexTweet(ctx, event)
	case "TweetDeleted":
		return h.handleDeleteTweet(ctx, event.AggregateID)
	default:
		log.Printf("Ignored unsupported event type: %s", event.Type)
		return nil
	}
}

func (h *IndexTweetHandler) handleIndexTweet(ctx context.Context, event nats.Event) error {
	var payload struct {
		ID        string `json:"tweetId"`
		AuthorID  string `json:"authorId"`
		Content   string `json:"content"`
		CreatedAt string `json:"createdAt"`
	}

	bytes, ok := event.Data.([]byte)
	if !ok {
		return fmt.Errorf("invalid payload type")
	}

	if err := json.Unmarshal(bytes, &payload); err != nil {
		return err
	}

	// Extract features
	hashtags, mentions := extractHashtagsAndMentions(payload.Content)

	// Parse createdAt
	createdAt, err := time.Parse(time.RFC3339, payload.CreatedAt)
	if err != nil {
		log.Printf("Failed to parse createdAt: %v, using now", err)
		createdAt = time.Now()
	}

	// Create domain model
	indexedTweet := domain.IndexedTweet{
		ID:        payload.ID,
		UserID:    payload.AuthorID,
		Text:      payload.Content,
		Hashtags:  hashtags,
		Mentions:  mentions,
		CreatedAt: createdAt,
	}

	return h.esClient.IndexTweet(ctx, indexedTweet)
}

func (h *IndexTweetHandler) handleDeleteTweet(ctx context.Context, tweetID string) error {
	return h.esClient.DeleteTweet(ctx, tweetID)
}

func extractHashtagsAndMentions(text string) (hashtags []string, mentions []string) {
	words := strings.Fields(text)

	for _, word := range words {
		if strings.HasPrefix(word, "#") && len(word) > 1 {
			hashtags = append(hashtags, strings.ToLower(word[1:]))
		} else if strings.HasPrefix(word, "@") && len(word) > 1 {
			mentions = append(mentions, word[1:])
		}
	}

	return hashtags, mentions
}
