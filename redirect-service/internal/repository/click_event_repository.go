package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/hoggir/re-path/redirect-service/internal/database"
	"github.com/hoggir/re-path/redirect-service/internal/domain"
	"go.mongodb.org/mongo-driver/mongo"
)

type ClickEventRepository interface {
	Create(ctx context.Context, clickEvent *domain.ClickEvent) error
}

type clickEventRepository struct {
	db         *database.MongoDB
	collection *mongo.Collection
}

func NewClickEventRepository(db *database.MongoDB) ClickEventRepository {
	collection := db.Collection(domain.ClickEvent{}.CollectionName())
	return &clickEventRepository{
		db:         db,
		collection: collection,
	}
}

func (r *clickEventRepository) Create(ctx context.Context, clickEvent *domain.ClickEvent) error {
	clickEvent.ClickedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, clickEvent)
	if err != nil {
		return fmt.Errorf("failed to create click event: %w", err)
	}

	return nil
}
