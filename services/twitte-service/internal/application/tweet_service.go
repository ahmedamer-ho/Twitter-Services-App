package application

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v8"

	"github.com/google/uuid"
	"github.com/yourusername/twitter-services-app/services/twitte-service/internal/domain"
)

type TweetService struct {
	repo     domain.TweetRepository
	esClient *elasticsearch.Client
}

func NewTweetService(repo domain.TweetRepository, es *elasticsearch.Client) *TweetService {
	return &TweetService{repo: repo, esClient: es}
}

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

func (s *TweetService) SearchTweets(
	ctx context.Context,
	query, hashtag, mention, user string,
	page, limit int, sort string,
) ([]domain.Tweet, error) {

	if s.esClient == nil {
		return nil, fmt.Errorf("search is currently unavailable")
	}

	// Build bool query
	var must []map[string]interface{}

	if query != "" {
		must = append(must, map[string]interface{}{
			"match": map[string]interface{}{
				"text": query,
			},
		})
	}
	if hashtag != "" {
		must = append(must, map[string]interface{}{
			"term": map[string]interface{}{
				"hashtags": hashtag,
			},
		})
	}
	if mention != "" {
		must = append(must, map[string]interface{}{
			"term": map[string]interface{}{
				"mentions": mention,
			},
		})
	}
	if user != "" {
		must = append(must, map[string]interface{}{
			"term": map[string]interface{}{
				"userId": user,
			},
		})
	}

	esQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": must,
			},
		},
		"from": (page - 1) * limit,
		"size": limit,
	}

	if sort == "createdAt" {
		esQuery["sort"] = []map[string]interface{}{
			{"createdAt": map[string]interface{}{"order": "desc"}}, // Usually we want new first by default
		}
	}

	var buf strings.Builder
	if err := json.NewEncoder(&buf).Encode(esQuery); err != nil {
		return nil, fmt.Errorf("error encoding query: %w", err)
	}

	res, err := s.esClient.Search(
		s.esClient.Search.WithContext(ctx),
		s.esClient.Search.WithIndex("tweets_index"),
		s.esClient.Search.WithBody(strings.NewReader(buf.String())),
		s.esClient.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, fmt.Errorf("error querying elasticsearch: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("elasticsearch error: %s", res.String())
	}

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	var results []domain.Tweet
	hits, ok := r["hits"].(map[string]interface{})
	if !ok {
		return results, nil // No hits payload
	}

	hitList, ok := hits["hits"].([]interface{})
	if !ok {
		return results, nil // Empty hits
	}

	for _, hit := range hitList {
		hitMap := hit.(map[string]interface{})
		source := hitMap["_source"].(map[string]interface{})

		// Reconstruct basic tweet. Missing idempotency/deleteAt since search index doesn't have it
		tweet := domain.Tweet{
			ID:       source["id"].(string),
			AuthorID: source["userId"].(string),
			Content:  source["text"].(string),
		}

		if caStr, ok := source["createdAt"].(string); ok {
			if t, err := time.Parse(time.RFC3339, caStr); err == nil {
				tweet.CreatedAt = t
			}
		}

		results = append(results, tweet)
	}

	return results, nil
}
