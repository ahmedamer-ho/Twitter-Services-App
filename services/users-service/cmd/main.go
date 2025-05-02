package main

import (
	"github.com/Twitter-Services-App/user-service/internal/auth"
	"github.com/Twitter-Services-App/user-service/internal/core"
	"github.com/Twitter-Services-App/user-service/internal/handlers"
	"github.com/Twitter-Services-App/user-service/internal/middlewares"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize Keycloak client
	keycloakClient := auth.NewKeycloakClient(
		"http://localhost:8080",
		"realm1",
		"go-app",
		"FjFT8TCcYKMQRsXobqdcJQxWcL0qIMlA",
		"admin",
		"admin",
	)

	// Create services
	authService := auth.NewKeycloakService(keycloakClient).(core.AuthService)

	// Create handlers
	authHandler := auth.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(authService)

	// Setup router
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	

	// Register routes
	authHandler.RegisterRoutes(router)

	// Protected routes
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
