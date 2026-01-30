package mongodb

import (
	"context"
	"time"

	"github.com/yourusername/twitter-services-app/services/twitte-service/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository struct {
	collection *mongo.Collection
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		collection: db.Collection("Tweets"),
	}
}
func (r *Repository) Insert(ctx context.Context, t domain.Tweet) error {
	_, err := r.collection.InsertOne(ctx, t)
	return err
}

func (r *Repository) FindByID(ctx context.Context, id string) (*domain.Tweet, error) {
	var t domain.Tweet
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&t)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *Repository) FindByAuthor(ctx context.Context, authorID string, limit int) ([]domain.Tweet, error) {
	cursor, err := r.collection.Find(ctx,
		bson.M{"authorId": authorID, "deletedAt": nil},
	)
	if err != nil {
		return nil, err
	}

	var tweets []domain.Tweet
	err = cursor.All(ctx, &tweets)
	return tweets, err
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

func (r *Repository) FindByIdempotencyKey(ctx context.Context, idempotencyKey string) (*domain.Tweet, error) {
	var t domain.Tweet
	err := r.collection.FindOne(ctx, bson.M{"idempotencyKey": idempotencyKey}).Decode(&t)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *Repository) InsertOutbox(ctx context.Context, event domain.OutboxEvent) error {
	collection := r.collection.Database().Collection("OutboxEvents")
	_, err := collection.InsertOne(ctx, event)
	return err
}

func (r *Repository) WithTransaction(ctx context.Context, fn func(domain.TweetRepository) error) error {
	session, err := r.collection.Database().Client().StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		return nil, fn(r)
	})
	return err
}

func (r *Repository) GetByIDs(ctx context.Context, ids []string) ([]domain.Tweet, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		return nil, err
	}
	var tweets []domain.Tweet
	err = cursor.All(ctx, &tweets)
	return tweets, err
}

func (r *Repository) GetLatestByAuthors(ctx context.Context, authorIDs []string) ([]domain.Tweet, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"authorId": bson.M{"$in": authorIDs}}, options.Find().SetSort(bson.M{"createdAt": -1}))
	if err != nil {
		return nil, err
	}
	var tweets []domain.Tweet
	err = cursor.All(ctx, &tweets)
	return tweets, err
}

func (r *Repository) FetchUnsentEvents(ctx context.Context, limit int) ([]domain.OutboxEvent, error) {
	collection := r.collection.Database().Collection("OutboxEvents")
	cursor, err := collection.Find(ctx, bson.M{"sent": false}, options.Find().SetLimit(int64(limit)))
	if err != nil {
		return nil, err
	}
	var events []domain.OutboxEvent
	err = cursor.All(ctx, &events)
	return events, err
}

func (r *Repository) MarkEventAsSent(ctx context.Context, id string) error {
	collection := r.collection.Database().Collection("OutboxEvents")
	_, err := collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"sent": true}})
	return err
}
