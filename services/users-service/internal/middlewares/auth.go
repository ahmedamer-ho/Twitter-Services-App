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
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		// Extract token from the header
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			return
		}

		// Validate the token with Keycloak
		ctx := c.Request.Context()
		rptResult, err := client.RetrospectToken(ctx, token, "user-service-client", "x78i9YgWY4B2ZyLntDBFIJKu5ocPUrPj", realm)
		if err != nil || !*rptResult.Active {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		c.Next()
	}
}
