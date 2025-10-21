package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type URL struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ShortCode   string             `bson:"shortCode" json:"shortCode"`
	OriginalURL string             `bson:"originalUrl" json:"originalUrl"`
	CustomAlias string             `bson:"customAlias,omitempty" json:"customAlias,omitempty"`
	UserID      int                `bson:"userId" json:"userId"`
	ClickCount  int                `bson:"clickCount" json:"clickCount"`
	IsActive    bool               `bson:"isActive" json:"isActive"`
	ExpiresAt   *time.Time         `bson:"expiresAt,omitempty" json:"expiresAt,omitempty"`
	Metadata    URLMetadata        `bson:"metadata" json:"metadata"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type URLMetadata struct {
	Title       string   `bson:"title,omitempty" json:"title,omitempty"`
	Description string   `bson:"description,omitempty" json:"description,omitempty"`
	Tags        []string `bson:"tags,omitempty" json:"tags,omitempty"`
}

type FindByShortCode struct {
	OriginalURL string     `bson:"originalUrl" json:"originalUrl"`
	IsActive    bool       `bson:"isActive" json:"isActive"`
	UserID      int        `bson:"userId" json:"userId"`
	ExpiresAt   *time.Time `bson:"expiresAt,omitempty" json:"expiresAt,omitempty"`
}

func (URL) CollectionName() string {
	return "urls"
}
