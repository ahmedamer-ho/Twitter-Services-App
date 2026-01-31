package mongodb

import (
	"context"
	"time"

	"github.com/yourusername/twitter-services-app/services/timeline-service/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository struct {
	collection *mongo.Collection
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		collection: db.Collection("Timelines"),
	}
}

func (r *Repository) SaveTweet(ctx context.Context, userID string, tweet domain.TimelineTweet) error {
	filter := bson.M{"_id": userID}
	update := bson.M{
		"$push": bson.M{
			"tweets": bson.M{
				"$each":  []domain.TimelineTweet{tweet},
				"$sort":  bson.M{"createdAt": -1},
				"$slice": 50, // Keep only top 50
			},
		},
		"$set": bson.M{"updatedAt": time.Now().UTC()},
	}

	opts := options.Update().SetUpsert(true)
	_, err := r.collection.UpdateOne(ctx, filter, update, opts)
	return err
}

func (r *Repository) GetTimeline(ctx context.Context, userID string) (*domain.Timeline, error) {
	var timeline domain.Timeline
	err := r.collection.FindOne(ctx, bson.M{"_id": userID}).Decode(&timeline)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &domain.Timeline{UserID: userID, Tweets: []domain.TimelineTweet{}}, nil
		}
		return nil, err
	}
	return &timeline, nil
}
