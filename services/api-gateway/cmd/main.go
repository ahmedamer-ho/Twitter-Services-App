package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Twitter-Services-App/api-gateway/internal/middlewares"
	"github.com/Twitter-Services-App/api-gateway/internal/proxy"
)

func main() {
	usersUrl := os.Getenv("USERS_SERVICE_URL")
	if usersUrl == "" {
		usersUrl = "http://localhost:8081"
	}
	tweetsUrl := os.Getenv("TWITTE_SERVICE_URL")
	if tweetsUrl == "" {
		tweetsUrl = "http://localhost:8082"
	}

	usersProxy, _ := proxy.NewReverseProxy(usersUrl)
	tweetsProxy, _ := proxy.NewReverseProxy(tweetsUrl)

	mux := http.NewServeMux()

	// USERS
	mux.Handle("/users/", usersProxy)
	mux.Handle("/auth/", usersProxy)

	// TWEETS
	mux.Handle("/tweets/", tweetsProxy)

	log.Println("API Gateway running on :8090")
	handler := middlewares.CorrelationID(mux)

	http.ListenAndServe(":8090", handler)
}
