package main

import (
	"github.com/Twitter-Services-App/user-service/internal/services"
	"github.com/Twitter-Services-App/user-service/internal/core"
	"github.com/Twitter-Services-App/user-service/internal/handlers"
	"github.com/Twitter-Services-App/user-service/internal/middlewares"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/Twitter-Services-App/user-service/internal/configs"
	"log"
)

func main() {
     // Load configuration
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Cannot load config:", err)
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
	authService := auth.NewKeycloakService(keycloakClient).(core.AuthService)

	//// 3. Inject service into handlers
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(authService)

	// Setup router
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	

	// Register routes
	authHandler.RegisterRoutes(router)

	//// 4. Inject client into middleware
	protected := router.Group("/")
	protected.Use(middlewares.KeycloakMiddleware(
		keycloakClient.Client,
		keycloakClient.Realm,
		keycloakClient.ClientID,
		keycloakClient.ClientSecret))
	{
		userHandler.RegisterRoutes(protected)
		protected.GET("/protected", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "This is a protected route"})
		})
	}

	router.Run(":8081")
}
