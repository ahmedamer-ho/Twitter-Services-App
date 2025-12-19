package main

import (
	"github.com/Twitter-Services-App/user-service/internal/logger"
	"github.com/Twitter-Services-App/user-service/internal/services"
	//"github.com/Twitter-Services-App/user-service/internal/core"
	"log"
	"net/http"

	"github.com/Twitter-Services-App/user-service/internal/configs"
	"github.com/Twitter-Services-App/user-service/internal/handlers"
	"github.com/Twitter-Services-App/user-service/internal/middlewares"
)

func main() {
     // Load configuration
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Cannot load config")
	}
	//1. Initialize Keycloak client with config
	keycloakClient := auth.NewKeycloakClient(
		cfg.Keycloak.URL,
		cfg.Keycloak.Realm,
		cfg.Keycloak.ClientID,
		cfg.Keycloak.ClientSecret,
		cfg.Keycloak.AdminUsername,
		cfg.Keycloak.AdminPassword,
	)

	//// 2. Inject client into service
	// High-level modules (handlers) depending on abstractions (AuthService)

    // Low-level details (Keycloak implementation) defined separately

    // Composition root (main.go) wiring everything together
	authService := auth.NewKeycloakService(keycloakClient)

	//// 3. Inject service into handlers
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(authService)

	// Setup router
	// ServeMux as router
	mux := http.NewServeMux()

	//Public routes with CORS
	publicMux := http.NewServeMux()
	authHandler.RegisterRoutes(publicMux)
	mux.Handle("/", middlewares.CORS(publicMux))

	// Protected routes using middleware
	protectedMux := http.NewServeMux()
		
	userHandler.RegisterRoutes(protectedMux)


	// protectedMux.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
	// 	w.Header().Set("Content-Type", "application/json")
	// 	w.Write([]byte(`{"message": "This is a protected route"}`))
	// })

	// Middleware wraps protected routes
	mux.Handle("/api/",middlewares.CORS( middlewares.KeycloakMiddleware(
		keycloakClient.Client,
		keycloakClient.Realm,
		keycloakClient.ClientID,
		keycloakClient.ClientSecret,
	)(protectedMux)))

	logger.Log.Info().Msg("Server running on :8081")
	
	log.Println("Server running on :8081")
	handler := middlewares.CorrelationID(mux)

	http.ListenAndServe(":8081", handler)
}
