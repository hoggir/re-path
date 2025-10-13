package database

import (
	"context"
	"fmt"
	"log"

	"github.com/hoggir/re-path/redirect-service/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoDB struct {
	Client   *mongo.Client
	Database *mongo.Database
}

func NewMongoDB(cfg *config.Config) (*MongoDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.MongoDB.Timeout)
	defer cancel()

	clientOptions := options.Client().
		ApplyURI(cfg.MongoDB.URI).
		SetMaxPoolSize(100).
		SetMinPoolSize(10)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mongodb: %w", err)
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, fmt.Errorf("failed to ping mongodb: %w", err)
	}

	database := client.Database(cfg.MongoDB.Database)

	log.Printf("✅ MongoDB connected successfully to database: %s", cfg.MongoDB.Database)

	return &MongoDB{
		Client:   client,
		Database: database,
	}, nil
}

func (m *MongoDB) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10)
	defer cancel()

	log.Println("🔌 Closing MongoDB connection...")
	if err := m.Client.Disconnect(ctx); err != nil {
		return fmt.Errorf("failed to disconnect from mongodb: %w", err)
	}

	log.Println("✅ MongoDB connection closed successfully")
	return nil
}

func (m *MongoDB) Collection(name string) *mongo.Collection {
	return m.Database.Collection(name)
}
