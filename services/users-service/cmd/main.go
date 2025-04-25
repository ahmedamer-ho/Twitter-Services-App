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
		"realm1",
		"go-app",
		"FjFT8TCcYKMQRsXobqdcJQxWcL0qIMlA",
		"admin", // Admin username amer.user
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
			Email     string `json:"email"`
			FirstName string `json:"firstname"`
			LastName  string `json:"lastname"`
			Password  string `json:"password"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		if err := keycloak.RegisterUser(req.Username, req.Email, req.FirstName, req.LastName, req.Password); err != nil {
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
	// Add these new routes after your existing routes
	r.GET("/users", func(c *gin.Context) {
		users, err := keycloak.GetUsers()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Convert gocloak users to simpler JSON structure
		var result []map[string]interface{}
		for _, user := range users {
			result = append(result, map[string]interface{}{
				"id":        user.ID,
				"username":  user.Username,
				"email":     user.Email,
				"firstName": user.FirstName,
				"lastName":  user.LastName,
				"enabled":   user.Enabled,
			})
		}

		c.JSON(http.StatusOK, gin.H{"users": result})
	})

	authenticated.GET("/users/:id", func(c *gin.Context) {
		userID := c.Param("id")
		user, err := keycloak.GetUserByID(userID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"user": map[string]interface{}{
				"id":        user.ID,
				"username":  user.Username,
				"email":     user.Email,
				"firstName": user.FirstName,
				"lastName":  user.LastName,
				"enabled":   user.Enabled,
			},
		})
	})
	r.Run(":8081")
}
