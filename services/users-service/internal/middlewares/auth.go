package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/Nerzal/gocloak/v12"
)

func KeycloakMiddleware(client *gocloak.GoCloak, realm string, clientID string, clientSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, `{"error":"authorization header required"}`, http.StatusUnauthorized)
				return
			}

			tokenParts := strings.Split(authHeader, " ")
			if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
				http.Error(w, `{"error":"invalid authorization format"}`, http.StatusUnauthorized)
				return
			}

			token := tokenParts[1]
			if token == "" {
				http.Error(w, `{"error":"empty token"}`, http.StatusUnauthorized)
				return
			}

			// Validate token
			rptResult, err := client.RetrospectToken(r.Context(), token, clientID, clientSecret, realm)
			if err != nil || !*rptResult.Active {
				http.Error(w, `{"error":"invalid or inactive token"}`, http.StatusUnauthorized)
				return
			}

			// Optionally store token or claims in context for later use
			ctx := context.WithValue(r.Context(), "token", token)
			ctx = context.WithValue(ctx, "token_claims", rptResult)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
