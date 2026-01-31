package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/Twitter-Services-App/api-gateway/internal/middlewares"
	"github.com/Twitter-Services-App/api-gateway/internal/proxy"
	"github.com/yourusername/twitter-services-app/shared/observability"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func main() {
	// Initialize Observability
	ctx := context.Background()
	shutdown, err := observability.InitProvider("api-gateway", "1.0.0")
	if err != nil {
		log.Fatal(err)
	}
	defer shutdown(ctx)

	usersUrl := os.Getenv("USERS_SERVICE_URL")
	if usersUrl == "" {
		usersUrl = "http://localhost:8081"
	}
	tweetsUrl := os.Getenv("TWITTE_SERVICE_URL")
	if tweetsUrl == "" {
		tweetsUrl = "http://localhost:8082"
	}
	timelineUrl := os.Getenv("TIMELINE_SERVICE_URL")
	if timelineUrl == "" {
		timelineUrl = "http://localhost:8083"
	}

	usersProxy, _ := proxy.NewReverseProxy(usersUrl)
	tweetsProxy, _ := proxy.NewReverseProxy(tweetsUrl)
	timelineProxy, _ := proxy.NewReverseProxy(timelineUrl)

	mux := http.NewServeMux()

	// USERS
	mux.Handle("/users/", usersProxy)
	mux.Handle("/auth/", usersProxy)

	// TWEETS
	mux.Handle("/tweets/", tweetsProxy)

	// TIMELINE
	mux.Handle("/timeline/", timelineProxy)

	log.Println("API Gateway running on :8090")
	handler := middlewares.CorrelationID(mux)

	// Wrap with OpenTelemetry
	otelHandler := otelhttp.NewHandler(handler, "api-gateway")

	http.ListenAndServe(":8090", otelHandler)
}
