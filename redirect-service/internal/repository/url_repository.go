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
)

var (
	ErrURLNotFound = errors.New("url not found")
	ErrURLExpired  = errors.New("url has expired")
	ErrURLInactive = errors.New("url is inactive")
)

type URLRepository interface {
	FindByShortCode(ctx context.Context, shortCode string) (*domain.URL, error)
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

func (r *urlRepository) FindByShortCode(ctx context.Context, shortCode string) (*domain.URL, error) {
	var url domain.URL

	// Simple query - hanya filter by shortCode
	// TIDAK filter expiresAt di query
	filter := bson.M{
		"shortCode": shortCode,
	}

	err := r.collection.FindOne(ctx, filter).Decode(&url)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrURLNotFound
		}
		return nil, fmt.Errorf("failed to find url: %w", err)
	}

	// Check kondisi SETELAH fetch dari database
	// 1. Check if inactive
	if !url.IsActive {
		return nil, ErrURLInactive
	}

	// 2. Check if expired
	if url.ExpiresAt != nil && url.ExpiresAt.Before(time.Now()) {
		return nil, ErrURLExpired
	}

	return &url, nil
}
