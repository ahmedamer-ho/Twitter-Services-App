package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/yourusername/twitter-services-app/services/twitte-service/internal/configs"
	httpadapter "github.com/yourusername/twitter-services-app/services/twitte-service/internal/adapters/http"
	"github.com/yourusername/twitter-services-app/services/twitte-service/internal/adapters/mongodb"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Cannot load config:",err)
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
	server := httpadapter.NewServer(router)

	log.Println("Tweet Service running on :8082")
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
