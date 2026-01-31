package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	httpadapter "github.com/yourusername/twitter-services-app/services/timeline-service/internal/adapters/http"
	"github.com/yourusername/twitter-services-app/services/timeline-service/internal/adapters/kafka"
	mongoadapter "github.com/yourusername/twitter-services-app/services/timeline-service/internal/adapters/mongodb"
	"github.com/yourusername/twitter-services-app/services/timeline-service/internal/application"
	config "github.com/yourusername/twitter-services-app/services/timeline-service/internal/configs"
	"github.com/yourusername/twitter-services-app/shared/observability"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize Observability
	shutdown, err := observability.InitProvider("timeline-service", "1.0.0")
	if err != nil {
		log.Fatal(err)
	}
	defer shutdown(ctx)

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Cannot load config:", err)
	}

	// MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoDB.URL))
	if err != nil {
		log.Fatal(err)
	}
	db := client.Database("TimelineDB")
	repo := mongoadapter.NewRepository(db)

	// Application Service
	timelineService := application.NewTimelineService(repo)

	// Kafka Consumer
	kafkaConsumer := kafka.NewConsumer(cfg.Kafka.Brokers, cfg.Kafka.Topic, cfg.Kafka.GroupID, timelineService)
	go kafkaConsumer.Start(ctx)
	defer kafkaConsumer.Close()

	// HTTP Handler
	handler := httpadapter.NewHandler(timelineService)
	mux := http.NewServeMux()
	mux.HandleFunc("/timeline/", handler.GetTimeline)

	// Graceful shutdown
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		<-ch
		cancel()
	}()

	port := 8083
	if cfg.Port != 0 {
		port = cfg.Port
	}
	addr := fmt.Sprintf(":%d", port)

	log.Printf("Timeline Service running on %s", addr)
	// Wrap with OpenTelemetry
	otelHandler := otelhttp.NewHandler(mux, "timeline-service")

	server := &http.Server{
		Addr:    addr,
		Handler: otelHandler,
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
