package handlers


import (
    "net/http"
	"encoding/json"
    "github.com/Twitter-Services-App/user-service/internal/core"
    "github.com/Twitter-Services-App/user-service/internal/models"
    "context"
    "strings" 
)

type AuthHandler struct {
    service core.AuthService
}

func NewAuthHandler(service core.AuthService) *AuthHandler {
    return &AuthHandler{service: service}
}


func (h *AuthHandler) RegisterRoutes(router *http.ServeMux) {
    router.HandleFunc("/", h.home)
    

    router.HandleFunc("/auth/register", h.register)
    router.HandleFunc("/auth/login", h.login)
    router.HandleFunc("/auth/logout", h.logout)
    
}

func (h *AuthHandler) home(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message":"users-service"}`))
}

func (h *AuthHandler) register(w http.ResponseWriter, r *http.Request) {
    var user models.User
    if err := json.NewDecoder(r.Body).Decode(&user);err != nil {
        http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
        return
    }

    if err := h.service.RegisterUser(&user); err != nil {
        http.Error(w, "Registration failed: "+err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message":"User registered successfully"}`))
}

func (h *AuthHandler) login(w http.ResponseWriter, r *http.Request) {
    var req models.LoginRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
        return
    }

    token, err := h.service.AuthenticateUser(req.Username, req.Password)
    if err != nil {
        http.Error(w, "Login failed: "+err.Error(), http.StatusUnauthorized)
        return
    }

    w.Header().Set("Content-Type", "application/json")
	resp := models.TokenResponse{Token: token}
	json.NewEncoder(w).Encode(resp)
}
func (h *AuthHandler) logout(w http.ResponseWriter, r *http.Request) {
    authHeader := r.Header.Get("Authorization")
    if authHeader == "" {
        http.Error(w, `{"error":"authorization header required"}`, http.StatusUnauthorized)
        return
    }

    // Extract Bearer token
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

    // Perform logout
    if err := h.service.Logout(context.Background(),token); err != nil {
        http.Error(w,
			`{"error":"logout failed","details":"`+err.Error()+`"}`,
			http.StatusInternalServerError,
		)
        return
    }

    // Clear client-side data
    w.Header().Set("Clear-Site-Data", `"cookies", "storage", "cache"`)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message":"logged out successfully"}`))
}