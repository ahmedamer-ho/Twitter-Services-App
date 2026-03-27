package main

import (
	"context"

	"github.com/Twitter-Services-App/user-service/internal/logger"
	auth "github.com/Twitter-Services-App/user-service/internal/services"

	//"github.com/Twitter-Services-App/user-service/internal/core"
	"log"
	"net/http"

	config "github.com/Twitter-Services-App/user-service/internal/configs"
	"github.com/Twitter-Services-App/user-service/internal/handlers"
	"github.com/Twitter-Services-App/user-service/internal/middlewares"
	"github.com/yourusername/twitter-services-app/shared/observability"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func main() {
	// Load configuration
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Cannot load config")
	}
	// Initialize Observability
	ctx := context.Background()
	shutdown, err := observability.InitProvider("users-service", "1.0.0")
	if err != nil {
		log.Fatal(err)
	}
	defer shutdown(ctx)

	if err := observability.InitLogger(); err != nil {
		log.Fatal(err)
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

	//Public routes
	publicMux := http.NewServeMux()
	authHandler.RegisterRoutes(publicMux)
	mux.Handle("/", publicMux)
	//Kubernetes / Load balancers need truth, not HTTP 200 lies.
	///health/live → process alive?
	///health/ready → dependencies ready?
	health := &handlers.HealthHandler{}
	mux.HandleFunc("/health/live", health.Live)
	mux.HandleFunc("/health/ready", health.Ready)

	// Protected routes using middleware
	protectedMux := http.NewServeMux()

	userHandler.RegisterRoutes(protectedMux)

	// protectedMux.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
	// 	w.Header().Set("Content-Type", "application/json")
	// 	w.Write([]byte(`{"message": "This is a protected route"}`))
	// })

	// Middleware wraps protected routes
	mux.Handle("/api/", middlewares.KeycloakMiddleware(
		keycloakClient.Client,
		keycloakClient.Realm,
		keycloakClient.ClientID,
		keycloakClient.ClientSecret,
	)(protectedMux))

	logger.Log.Info().Msg("Server running on :8081")

	log.Println("Server running on :8081")
	handler := middlewares.CorrelationID(mux)

	// Wrap with OpenTelemetry
	otelHandler := otelhttp.NewHandler(handler, "users-service")

	http.ListenAndServe(":8081", otelHandler)
}
