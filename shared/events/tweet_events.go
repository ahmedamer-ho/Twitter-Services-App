package events

import "time"

// TweetCreated defines the event payload when a new tweet is posted
type TweetCreated struct {
	TweetID   string    `json:"tweetId"`
	AuthorID  string    `json:"authorId"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
}

// TweetDeleted defines the event payload when a tweet is deleted
type TweetDeleted struct {
	TweetID   string    `json:"tweetId"`
	AuthorID  string    `json:"authorId"`
	DeletedAt time.Time `json:"deletedAt"`
}

// TweetLiked defines the event payload when a tweet is liked
type TweetLiked struct {
	TweetID   string    `json:"tweetId"`
	UserID    string    `json:"userId"`
	LikedAt   time.Time `json:"likedAt"`
}
