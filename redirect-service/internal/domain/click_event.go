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
	ReferrerDomain string             `bson:"referrerDomain,omitempty" json:"referrerDomain,omitempty"`
	CountryCode    string             `bson:"countryCode,omitempty" json:"countryCode,omitempty"`
	City           string             `bson:"city,omitempty" json:"city,omitempty"`
	Region         string             `bson:"region,omitempty" json:"region,omitempty"`
	DeviceType     string             `bson:"deviceType,omitempty" json:"deviceType,omitempty"` // desktop, mobile, tablet
	BrowserName    string             `bson:"browserName,omitempty" json:"browserName,omitempty"`
	BrowserVersion string             `bson:"browserVersion,omitempty" json:"browserVersion,omitempty"`
	OSName         string             `bson:"osName,omitempty" json:"osName,omitempty"`
	OSVersion      string             `bson:"osVersion,omitempty" json:"osVersion,omitempty"`
	Lat            float64            `bson:"lat,omitempty" json:"lat,omitempty"`
	Lon            float64            `bson:"lon,omitempty" json:"lon,omitempty"`
	IsBot          bool               `bson:"isBot" json:"isBot"`
}

func (ClickEvent) CollectionName() string {
	return "click_events"
}

// PayloadElasticClick represents the structure for sending click event data to Elasticsearch
type PayloadElasticClick struct {
	IndexType string    `json:"index_type"`
	Data      ClickData `json:"data"`
}

type ClickData struct {
	ShortCode string        `json:"short_code"`
	Metadata  ClickMetaData `json:"metadata"`
}

type ClickMetaData struct {
	ClickedAt time.Time     `json:"clicked_at"`
	IsBot     bool          `json:"is_bot"`
	Client    ClientInfo    `json:"client"`
	HTTP      HTTPInfo      `json:"http"`
	UserAgent UserAgentInfo `json:"user_agent"`
}

type ClientInfo struct {
	IPHash string  `json:"ip_hash"`
	Geo    GeoInfo `json:"geo"`
}

type GeoInfo struct {
	CountryISOCode string             `json:"country_iso_code"`
	RegionName     string             `json:"region_name,omitempty"`
	City           string             `json:"city"`
	Location       GeoLocationElastic `json:"location"`
}

type GeoLocationElastic struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type HTTPInfo struct {
	Referrer       string `json:"referrer"`
	ReferrerDomain string `json:"referrer_domain"`
}

type UserAgentInfo struct {
	Original string      `json:"original"`
	Device   DeviceInfo  `json:"device"`
	Browser  BrowserInfo `json:"browser"`
	OS       OSInfo      `json:"os"`
}

type DeviceInfo struct {
	Name string `json:"name"`
}

type BrowserInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type OSInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// End PayloadElasticClick
