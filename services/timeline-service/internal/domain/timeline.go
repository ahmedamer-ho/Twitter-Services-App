package domain

import "time"

type TimelineTweet struct {
	ID        string    `json:"id" bson:"_id"`
	AuthorID  string    `json:"authorId" bson:"authorId"`
	Content   string    `json:"content" bson:"content"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
}

type Timeline struct {
	UserID    string          `json:"userId" bson:"_id"`
	Tweets    []TimelineTweet `json:"tweets" bson:"tweets"`
	UpdatedAt time.Time       `json:"updatedAt" bson:"updatedAt"`
}
