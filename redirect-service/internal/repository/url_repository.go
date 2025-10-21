package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hoggir/re-path/redirect-service/internal/database"
	"github.com/hoggir/re-path/redirect-service/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrURLNotFound = errors.New("url not found")
	ErrURLExpired  = errors.New("url has expired")
	ErrURLInactive = errors.New("url is inactive")
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
			return nil, ErrURLNotFound
		}
		return nil, fmt.Errorf("failed to find url: %w", err)
	}

	if !url.IsActive {
		return nil, ErrURLInactive
	}

	if url.ExpiresAt != nil && url.ExpiresAt.Before(time.Now()) {
		return nil, ErrURLExpired
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
		return fmt.Errorf("failed to increment click count: %w", err)
	}

	if result.MatchedCount == 0 {
		return ErrURLNotFound
	}

	return nil
}
