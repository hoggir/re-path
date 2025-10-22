package repository

import (
	"context"
	"time"

	"github.com/hoggir/re-path/redirect-service/internal/database"
	"github.com/hoggir/re-path/redirect-service/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type URLRepository interface {
	FindByShortCode(ctx context.Context, shortCode string) (*domain.FindByShortCode, error)
	IncrementClickCount(ctx context.Context, shortCode string) error
}

type urlRepository struct {
	db         *database.MongoDB
	collection *mongo.Collection
}

func NewURLRepository(db *database.MongoDB) URLRepository {
	return &urlRepository{
		db:         db,
		collection: db.Collection(domain.URL{}.CollectionName()),
	}
}

func (r *urlRepository) FindByShortCode(ctx context.Context, shortCode string) (*domain.FindByShortCode, error) {
	var url domain.FindByShortCode

	filter := bson.M{
		"shortCode": shortCode,
	}

	err := r.collection.FindOne(ctx, filter, options.FindOne().SetProjection(bson.M{"userId": 1, "originalUrl": 1, "isActive": 1, "expiresAt": 1, "_id": 0})).Decode(&url)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrURLNotFound.WithContext("shortCode", shortCode)
		}
		return nil, domain.ErrDatabaseError.
			WithContext("shortCode", shortCode).
			WithContext("operation", "FindByShortCode").
			Wrap(err)
	}

	if !url.IsActive {
		return nil, domain.ErrURLInactive.
			WithContext("shortCode", shortCode).
			WithContext("userId", url.UserID)
	}

	if url.ExpiresAt != nil && url.ExpiresAt.Before(time.Now()) {
		return nil, domain.ErrURLExpired.
			WithContext("shortCode", shortCode).
			WithContext("expiresAt", url.ExpiresAt)
	}

	return &url, nil
}

func (r *urlRepository) IncrementClickCount(ctx context.Context, shortCode string) error {
	filter := bson.M{
		"shortCode": shortCode,
	}

	update := bson.M{
		"$inc": bson.M{
			"clickCount": 1,
		},
		"$set": bson.M{
			"updatedAt": time.Now(),
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return domain.ErrDatabaseError.
			WithContext("shortCode", shortCode).
			WithContext("operation", "IncrementClickCount").
			Wrap(err)
	}

	if result.MatchedCount == 0 {
		return domain.ErrURLNotFound.WithContext("shortCode", shortCode)
	}

	return nil
}
