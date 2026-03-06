package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/yourusername/twitter-services-app/services/search-indexer/internal/application"
	"github.com/yourusername/twitter-services-app/services/search-indexer/internal/infrastructure"
	"github.com/yourusername/twitter-services-app/shared/nats"
	"github.com/yourusername/twitter-services-app/shared/observability"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 1. Telemetry
	shutdown, err := observability.InitProvider("search-indexer", "1.0.0")
	if err != nil {
		log.Printf("Failed to initialize telemetry: %v", err)
	} else {
		defer shutdown(ctx)
	}

	// 2. Load Env
	elasticURL := os.Getenv("ELASTICSEARCH_URL")
	if elasticURL == "" {
		elasticURL = "http://localhost:9200"
	}
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = "nats://localhost:4222"
	}
	subject := os.Getenv("NATS_TWEET_SUBJECT")
	if subject == "" {
		subject = "tweets"
	}
	qGroup := os.Getenv("NATS_QUEUE_GROUP")
	if qGroup == "" {
		qGroup = "search-indexer-group"
	}

	// 3. Connect Elastic
	esClient, err := infrastructure.NewESClient(elasticURL)
	if err != nil {
		log.Fatalf("Failed to create Elasticsearch client: %v", err)
	}
	if err := esClient.IndexSetup(ctx); err != nil {
		log.Fatalf("Failed to initialize ES Index: %v", err)
	}
	log.Println("Connected to Elasticsearch successfully.")

	// 4. Application Handler
	handler := application.NewIndexTweetHandler(esClient)

	// 5. Connect NATS & Start Consumer
	nc, js, err := nats.InitNATS(natsURL)
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	defer nc.Close()

	consumer := nats.NewConsumer(js, subject, qGroup, handler.HandleTweetEvent)
	go consumer.Start(ctx)
	log.Printf("Search Indexer Service listening on NATS Subject: %s", subject)

	// Graceful Shutdown OS Hook
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	log.Println("Shutting down gracefully...")
}
