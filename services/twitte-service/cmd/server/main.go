package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	httpadapter "github.com/yourusername/twitter-services-app/services/twitte-service/internal/adapters/http"
	"github.com/yourusername/twitter-services-app/services/twitte-service/internal/adapters/mongodb"
	config "github.com/yourusername/twitter-services-app/services/twitte-service/internal/configs"
	"github.com/yourusername/twitter-services-app/shared/observability"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

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

	_ = mongodb.NewRepository(client.Database("Tweets")) // used later

	// HTTP
	router := httpadapter.NewRouter()

	// Add OpenTelemetry Middleware
	handler := otelhttp.NewHandler(router, "twitte-service")
	server := httpadapter.NewServer(handler)

	log.Println("Tweet Service running on :8082")
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
