package middlewares

import (
	"net/http"
	"strings"

	"github.com/Nerzal/gocloak/v12"
	"github.com/gin-gonic/gin"
)

func KeycloakMiddleware(client *gocloak.GoCloak, realm string, clientID string, clientSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
        if authHeader == "" {
            authHeader = c.Request.Header.Get("authorization")
        }
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":       "authorization_header_required",
				"message":     "Please include a valid Bearer token in the Authorization header",
			})
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":       "invalid_authorization_header",
				"message":     "Format should be: 'Bearer <token>'",
			})
			return
		}

		token := tokenParts[1]
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":       "empty_token",
				"message":     "Token cannot be empty",
			})
			return
		}

		rptResult, err := client.RetrospectToken(c.Request.Context(), token, clientID, clientSecret, realm)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":       "token_validation_failed",
				"message":     "Failed to validate token",
			})
			return
		}

		if !*rptResult.Active {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":       "inactive_token",
				"message":     "Token is not active",
			})
			return
		}

		c.Set("token", token)
		c.Set("token_claims", rptResult)
		c.Next()
	}
}