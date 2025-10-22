package database

import (
	"context"
	"fmt"

	"github.com/hoggir/re-path/redirect-service/internal/config"
	"github.com/hoggir/re-path/redirect-service/internal/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoDB struct {
	Client   *mongo.Client
	Database *mongo.Database
	logger   logger.Logger
}

func NewMongoDB(cfg *config.Config, log logger.Logger) (*MongoDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.MongoDB.ConnTimeout)
	defer cancel()

	clientOptions := options.Client().
		ApplyURI(cfg.MongoDB.URI).
		SetMaxPoolSize(cfg.MongoDB.MaxPoolSize).
		SetMinPoolSize(cfg.MongoDB.MinPoolSize).
		SetTimeout(cfg.MongoDB.QueryTimeout)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mongodb: %w", err)
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, fmt.Errorf("failed to ping mongodb: %w", err)
	}

	database := client.Database(cfg.MongoDB.Database)

	log.Info("MongoDB connected successfully",
		"database", cfg.MongoDB.Database,
		"minPoolSize", cfg.MongoDB.MinPoolSize,
		"maxPoolSize", cfg.MongoDB.MaxPoolSize,
		"queryTimeout", cfg.MongoDB.QueryTimeout)

	return &MongoDB{
		Client:   client,
		Database: database,
		logger:   log,
	}, nil
}

func (m *MongoDB) Close(cfg *config.Config) error {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.MongoDB.DisconnTimeout)
	defer cancel()

	m.logger.Info("closing MongoDB connection")
	if err := m.Client.Disconnect(ctx); err != nil {
		return fmt.Errorf("failed to disconnect from mongodb: %w", err)
	}

	m.logger.Info("MongoDB connection closed successfully")
	return nil
}

func (m *MongoDB) Collection(name string) *mongo.Collection {
	return m.Database.Collection(name)
}
