package mongodb

import (
	"context"
	"time"

	"github.com/yourusername/twitter-services-app/services/twitte-service/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)
type Repository struct {
	collection *mongo.Collection
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
			collection: db.Collection("twites"),
	}
}
func (r *Repository) Insert(ctx context.Context, t domain.Twite) error {
	_, err := r.collection.InsertOne(ctx, t)
	return err
}
func (r *Repository) FindByAuthor(ctx context.Context, authorID string, limit int) ([]domain.Twite, error) {
	cursor, err := r.collection.Find(ctx,
			bson.M{"authorId": authorID, "deletedAt": nil},
	)
	if err != nil {
			return nil, err
	}

	var twites []domain.Twite
	err = cursor.All(ctx, &twites)
	return twites, err
}
func (r *Repository) SoftDelete(ctx context.Context, id string) error {
	now := time.Now().UTC()
	expire := now.Add(30 * 24 * time.Hour)

	_, err := r.collection.UpdateOne(ctx,
			bson.M{"_id": id},
			bson.M{
					"$set": bson.M{
							"deletedAt": now,
							"expireAt":  expire,
					},
			},
	)

	return err
}