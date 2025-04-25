package middlewares

import (
	"net/http"
	"strings"

	"github.com/Nerzal/gocloak/v12"
	"github.com/gin-gonic/gin"
)

// KeycloakMiddleware enforces authentication via Keycloak.
func KeycloakMiddleware(client gocloak.GoCloak, realm string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip middleware for certain routes if needed
		if c.Request.URL.Path == "/" || c.Request.URL.Path == "/login" || c.Request.URL.Path == "/register" {
			c.Next()
			return
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":       "Authorization header required",
				"description": "Please include a valid Bearer token in the Authorization header",
			})
			return
		}

		// Extract token from the header
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":       "Invalid authorization header format",
				"description": "Format should be: 'Bearer <token>'",
			})
			return
		}

		token := tokenParts[1]
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Empty token provided"})
			return
		}

		// Validate the token with Keycloak
		ctx := c.Request.Context()
		rptResult, err := client.RetrospectToken(ctx, token, "client1", "FjFT8TCcYKMQRsXobqdcJQxWcL0qIMlA", realm)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":       "Token validation failed",
				"description": err.Error(),
			})
			return
		}

		if !*rptResult.Active {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token is not active"})
			return
		}

		// Store token claims in context for later use if needed
		c.Set("token", token)
		c.Set("token_claims", rptResult)

		c.Next()
	}
}