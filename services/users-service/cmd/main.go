package main

import (
	"net/http"

	"github.com/Twitter-Services-App/user-service/internal/auth"
	"github.com/Twitter-Services-App/user-service/internal/middlewares"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	keycloak := auth.NewKeycloakClient(
		"http://localhost:8080",
		"user-service-realm",
		"user-service-client",
		"FjFT8TCcYKMQRsXobqdcJQxWcL0qIMlA",
		"amer.user", // Admin username amer.user
		"admin", // Admin password admin
	)

	r := gin.Default()
	r.Use(cors.Default())
	// Apply Keycloak middleware (pass client as interface)
	// Public route
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "users-services"})
	})

	// Protected route with middleware
	authenticated := r.Group("/")
	authenticated.Use(middlewares.KeycloakMiddleware(*keycloak.Client, keycloak.Realm))
	{
		authenticated.GET("/protected", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "This is a protected route"})
		})
	}

	r.POST("/register", func(c *gin.Context) {
		var req struct {
			Username  string `json:"username"`
			FirstName string `json:"firstname"`
			LastName  string `json:"lastname"`
			Password  string `json:"password"`
			Email     string `json:"email"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		if err := keycloak.RegisterUser(req.Username,req.FirstName,req.LastName, req.Password, req.Email); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(201, gin.H{"message": "User registered successfully"})
	})
	r.POST("/login", func(c *gin.Context) {
		var req struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		token, err := keycloak.AuthenticateUser(req.Username, req.Password)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Login failed"})
			return
		}

		c.JSON(200, gin.H{"token": token})
	})

	r.Run(":8081")
}
