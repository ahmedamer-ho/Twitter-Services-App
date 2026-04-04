package events

import "time"

// UserCreated defines the event payload for when a user is created
type UserCreated struct {
	UserID    string    `json:"userId"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
}

// UserFollowed defines the event payload when one user follows another
type UserFollowed struct {
	FollowerID  string    `json:"followerId"`
	FollowingID string    `json:"followingId"`
	FollowedAt  time.Time `json:"followedAt"`
}
