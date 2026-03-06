package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	httpadapter "github.com/yourusername/twitter-services-app/services/twitte-service/internal/adapters/http"
	mongoadapter "github.com/yourusername/twitter-services-app/services/twitte-service/internal/adapters/mongodb"
	"github.com/yourusername/twitter-services-app/services/twitte-service/internal/application"
	config "github.com/yourusername/twitter-services-app/services/twitte-service/internal/configs"
	"github.com/yourusername/twitter-services-app/shared/nats"
	"github.com/yourusername/twitter-services-app/shared/observability"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type natsPublisherAdapter struct {
	producer nats.Producer
	subject  string
}

func (a *natsPublisherAdapter) Publish(ctx context.Context, key string, value []byte) error {
	natsEv := nats.Event{
		EventID:     key,
		AggregateID: key,
		Payload:     json.RawMessage(value),
	}
	return a.producer.Publish(ctx, a.subject, natsEv)
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize Observability (Tracing & Metrics)
	shutdown, err := observability.InitProvider("twitte-service", "1.0.0")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			log.Fatal("failed to shutdown TracerProvider: %w", err)
		}
	}()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Cannot load config:", err)
	}
	// Graceful shutdown
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		<-ch
		cancel()
	}()

	// MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoDB.URL))
	if err != nil {
		log.Fatal(err)
	}

	// NATS
	if cfg.NATS.URL == "" {
		cfg.NATS.URL = "nats://nats:4222"
	}
	if cfg.NATS.Subject == "" {
		cfg.NATS.Subject = "tweets"
	}

	// Create NATS connection & JetStream Context
	nc, js, err := nats.InitNATS(cfg.NATS.URL)
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	defer nc.Close()

	producer := nats.NewProducer(js)

	// Repositories
	repo := mongoadapter.NewRepository(client.Database("Tweets"))

	// Application Services
	tweetService := application.NewTweetService(repo)

	// Outbox Worker
	adapter := &natsPublisherAdapter{producer: producer, subject: cfg.NATS.Subject}
	outboxWorker := application.NewOutboxWorker(repo, adapter, 5*time.Second)
	go outboxWorker.Start(ctx)

	// HTTP Handler
	tweetHandler := httpadapter.NewHandler(tweetService)

	// HTTP
	router := httpadapter.NewRouter(tweetHandler)

	// Add OpenTelemetry Middleware
	handler := otelhttp.NewHandler(router, "twitte-service")
	server := httpadapter.NewServer(handler)

	log.Println("Tweet Service running on :8082")
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
