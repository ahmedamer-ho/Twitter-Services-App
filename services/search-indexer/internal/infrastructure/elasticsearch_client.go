package infrastructure

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/yourusername/twitter-services-app/services/search-indexer/internal/domain"
)

const indexName = "tweets_index"

type ESClient struct {
	client *elasticsearch.Client
}

func NewESClient(url string) (*ESClient, error) {
	cfg := elasticsearch.Config{
		Addresses: []string{url},
		// The elasticsearch library natively handles retries for 429 & 502/503/504 status codes
		MaxRetries: 5,
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("error creating elasticsearch client: %w", err)
	}

	return &ESClient{client: es}, nil
}

func (es *ESClient) IndexSetup(ctx context.Context) error {
	// Check if index exists
	res, err := es.client.Indices.Exists([]string{indexName})
	if err != nil {
		return fmt.Errorf("error checking if index exists: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		return nil // Index already exists
	}

	// Create index with mapping
	mapping := `{
		"mappings": {
			"properties": {
				"id": { "type": "keyword" },
				"userId": { "type": "keyword" },
				"text": { "type": "text" },
				"hashtags": { "type": "keyword" },
				"mentions": { "type": "keyword" },
				"createdAt": { "type": "date" }
			}
		}
	}`

	req := esapi.IndicesCreateRequest{
		Index: indexName,
		Body:  bytes.NewReader([]byte(mapping)),
	}

	createRes, err := req.Do(ctx, es.client)
	if err != nil {
		return fmt.Errorf("error creating index: %w", err)
	}
	defer createRes.Body.Close()

	if createRes.IsError() {
		return fmt.Errorf("failed to create index. Status: %s", createRes.String())
	}

	log.Printf("Successfully created Elasticsearch index '%s'", indexName)
	return nil
}

func (es *ESClient) IndexTweet(ctx context.Context, tweet domain.IndexedTweet) error {
	data, err := json.Marshal(tweet)
	if err != nil {
		return fmt.Errorf("error marshaling tweet: %w", err)
	}

	req := esapi.IndexRequest{
		Index:      indexName,
		DocumentID: tweet.ID,
		Body:       bytes.NewReader(data),
		Refresh:    "true", // Refresh true for immediate visibility in dev/search
	}

	res, err := req.Do(ctx, es.client)
	if err != nil {
		return fmt.Errorf("error indexing document: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error response from Elasticsearch: %s", res.String())
	}

	return nil
}

func (es *ESClient) DeleteTweet(ctx context.Context, tweetID string) error {
	req := esapi.DeleteRequest{
		Index:      indexName,
		DocumentID: tweetID,
		Refresh:    "true",
	}

	res, err := req.Do(ctx, es.client)
	if err != nil {
		return fmt.Errorf("error deleting document: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error response from Elasticsearch: %s", res.String())
	}

	return nil
}
