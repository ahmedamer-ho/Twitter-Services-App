func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "" {
			http.Error(w, "missing token", http.StatusUnauthorized)
			return
		}

		// validate JWT (Keycloak public key)
		claims, err := ValidateToken(auth)
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		// Forward claims
		r.Header.Set("X-User-ID", claims.UserID)
		r.Header.Set("X-Tenant-ID", claims.TenantID)
		r.Header.Set("X-Roles", strings.Join(claims.Roles, ","))

		next.ServeHTTP(w, r)
	})
}
