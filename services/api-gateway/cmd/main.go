package main

import (
	"log"
	"net/http"

	"github.com/Twitter-Services-App/api-gateway/internal/proxy"
)

func main() {
	usersProxy, _ := proxy.NewReverseProxy("http://localhost:8081")
	tweetsProxy, _ := proxy.NewReverseProxy("http://localhost:8082")

	mux := http.NewServeMux()

	// USERS
	mux.Handle("/users/", usersProxy)
	mux.Handle("/auth/", usersProxy)

	// TWEETS
	mux.Handle("/tweets/", tweetsProxy)

	log.Println("API Gateway running on :8090")
	handler := CorrelationID(mux)

	http.ListenAndServe(":8090", handler)
}
