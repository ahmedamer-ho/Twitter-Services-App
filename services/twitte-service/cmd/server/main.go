package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/yourusername/twitter-services-app/services/twitte-service/internal/adapters/mongodb"
	
    "github.com/Twitter-Services-App/twite-service/internal/adapters/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		<-ch
		cancel()
	}()

	// MongoDB connection
	mongoClient, err := mongo.Connect(ctx, os.Getenv("MONGO_URI"))
	if err != nil {
		log.Fatal(err)
	}

	twiteRepo := mongodb.NewRepository(mongoClient.Database("twites"))

	server := http.NewServer(twiteRepo) // Assuming NewServer expects a twiteRepo

	log.Println("Twite Service running on :8082")
	if err := server.Run(":8082 "); err != nil {
		log.Fatal(err)
	}
}
