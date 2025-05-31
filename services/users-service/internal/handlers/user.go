package handlers

import (
    "net/http"
    "strings"
    "github.com/Twitter-Services-App/user-service/internal/core"
    "encoding/json"
)

type UserHandler struct {
    service core.AuthService
}

func NewUserHandler(service core.AuthService) *UserHandler {
    return &UserHandler{service: service}
}


func (h *UserHandler) RegisterRoutes(router *http.ServeMux) {
    router.HandleFunc("/api/users", h.getUsers)
    router.HandleFunc("/api/users/", h.getUserByID)
}

func (h *UserHandler) getUsers(w http.ResponseWriter, r *http.Request) {
    users, err := h.service.GetUsers()
    if err != nil {
        http.Error(w, "Failed to fetch users: "+err.Error(), http.StatusInternalServerError)
        return
    }

   
	w.Header().Set("Content-Type", "application/json")
	// Wrap results in a JSON object with key "users"
	resp := map[string]interface{}{"users": users}
	json.NewEncoder(w).Encode(resp)
}

func (h *UserHandler) getUserByID(w http.ResponseWriter, r *http.Request) {

    parts := strings.Split(r.URL.Path, "/")
    if len(parts) < 4 {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }
    userID := parts[3] 
	if userID == "" {
		http.Error(w, "User ID missing in URL", http.StatusBadRequest)
		return
	}
    user, err := h.service.GetUserByID(userID)
    if err != nil {
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }

   w.Header().Set("Content-Type", "application/json")
	resp := map[string]interface{}{"user": user}
	json.NewEncoder(w).Encode(resp)
}