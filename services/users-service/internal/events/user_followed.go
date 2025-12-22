package events

type UserFollowedPayload struct {
	FollowerID string `json:"followerId"`
	FollowedID string `json:"followedId"`
}
