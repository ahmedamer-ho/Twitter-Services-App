package domain

import "time"

type Twite struct {
	ID        string     `bson:"_id" json:"id"`
	AuthorID  string     `bson:"authorId" json:"authorId"`
	Content   string     `bson:"content" json:"content"`
	CreatedAt time.Time  `bson:"createdAt" json:"createdAt"`
	DeletedAt *time.Time `bson:"deletedAt,omitempty" json:"deletedAt,omitempty"`
	ExpireAt  *time.Time `bson:"expireAt,omitempty" json:"-"`
}
