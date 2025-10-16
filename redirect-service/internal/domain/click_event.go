package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ClickEvent struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ShortCode      string             `bson:"shortCode" json:"shortCode"`
	ClickedAt      time.Time          `bson:"clickedAt" json:"clickedAt"`
	IPAddressHash  string             `bson:"ipAddressHash" json:"ipAddressHash"`
	UserAgent      string             `bson:"userAgent" json:"userAgent"`
	ReferrerURL    string             `bson:"referrerUrl,omitempty" json:"referrerUrl,omitempty"`
	CountryCode    string             `bson:"countryCode,omitempty" json:"countryCode,omitempty"`
	City           string             `bson:"city,omitempty" json:"city,omitempty"`
	Region         string             `bson:"region,omitempty" json:"region,omitempty"`
	DeviceType     string             `bson:"deviceType,omitempty" json:"deviceType,omitempty"` // desktop, mobile, tablet
	BrowserName    string             `bson:"browserName,omitempty" json:"browserName,omitempty"`
	BrowserVersion string             `bson:"browserVersion,omitempty" json:"browserVersion,omitempty"`
	OSName         string             `bson:"osName,omitempty" json:"osName,omitempty"`
	OSVersion      string             `bson:"osVersion,omitempty" json:"osVersion,omitempty"`
	IsBot          bool               `bson:"isBot" json:"isBot"`
}

func (ClickEvent) CollectionName() string {
	return "click_events"
}
