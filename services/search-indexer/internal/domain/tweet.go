package domain

import "time"

type IndexedTweet struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	Text      string    `json:"text"`
	Hashtags  []string  `json:"hashtags"`
	Mentions  []string  `json:"mentions"`
	CreatedAt time.Time `json:"createdAt"`
}
