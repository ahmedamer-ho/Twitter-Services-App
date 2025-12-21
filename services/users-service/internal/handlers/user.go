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
// func (h *UserHandler) FollowUser(w http.ResponseWriter, r *http.Request) {
// 	followerID := chi.URLParam(r, "id")
// 	followedID := chi.URLParam(r, "targetId")

// 	if err := h.userService.FollowUser(r.Context(), followerID, followedID); err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	w.WriteHeader(http.StatusNoContent)
// }
// func (s *UserService) FollowUser(
// 	ctx context.Context,
// 	followerID string,
// 	followedID string,
// ) error {

// 	// 1️⃣ Business logic
// 	if err := s.repo.AddFollower(ctx, followerID, followedID); err != nil {
// 		return err
// 	}

// 	// 2️⃣ Create event AFTER success
// 	event := kafka.Event{
// 		EventID:       uuid.New().String(),
// 		EventType:     "UserFollowed",
// 		AggregateID:   followerID,
// 		Timestamp:     time.Now().UTC(),
// 		CorrelationID: observability.FromContext(ctx),
// 		Payload: events.UserFollowedPayload{
// 			FollowerID: followerID,
// 			FollowedID: followedID,
// 		},
// 	}

// 	// 3️⃣ Publish event
// 	if err := s.producer.Publish(ctx, "user-events", event); err != nil {
// 		// IMPORTANT: do NOT break user flow
// 		// Log + retry later (outbox later)
// 		return nil
// 	}

// 	return nil
// }
