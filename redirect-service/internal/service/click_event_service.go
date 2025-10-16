package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hoggir/re-path/redirect-service/internal/domain"
	"github.com/hoggir/re-path/redirect-service/internal/repository"
	"github.com/mileusna/useragent"
)

type ClickEventService interface {
	TrackClick(ctx context.Context, ginCtx *gin.Context, shortCode string) error
}

type clickEventService struct {
	clickEventRepo  repository.ClickEventRepository
	geoIPService    GeoIPService
	rabbitmqService RabbitMQService
}

func NewClickEventService(
	clickEventRepo repository.ClickEventRepository,
	geoIPService GeoIPService,
	rabbitmqService RabbitMQService,
) ClickEventService {
	return &clickEventService{
		clickEventRepo:  clickEventRepo,
		geoIPService:    geoIPService,
		rabbitmqService: rabbitmqService,
	}
}

func (s *clickEventService) TrackClick(ctx context.Context, ginCtx *gin.Context, shortCode string) error {
	newDate := time.Now().UTC()
	ua := useragent.Parse(ginCtx.Request.UserAgent())
	clientIP := ginCtx.ClientIP()
	log.Printf("📍 Tracking click from IP: %s", clientIP)

	ipHash := hashIPAddress(clientIP)

	deviceType := getDeviceType(ua)

	referrer := ginCtx.Request.Referer()
	referrerDomain := extractDomain(referrer)

	// geoLocation, err := s.geoIPService.GetLocation(ctx, clientIP)
	geoLocation, err := s.geoIPService.GetLocation(ctx, "203.175.11.126")
	if err != nil {
		log.Printf("⚠️  Failed to get geolocation for IP %s: %v", clientIP, err)
	}

	clickEvent := &domain.ClickEvent{
		ClickedAt:      newDate,
		ShortCode:      shortCode,
		IPAddressHash:  ipHash,
		UserAgent:      ginCtx.Request.UserAgent(),
		ReferrerURL:    referrer,
		ReferrerDomain: referrerDomain,
		DeviceType:     deviceType,
		BrowserName:    ua.Name,
		BrowserVersion: ua.Version,
		OSName:         ua.OS,
		OSVersion:      ua.OSVersion,
		IsBot:          ua.Bot,
	}

	if geoLocation != nil {
		clickEvent.CountryCode = geoLocation.CountryCode
		clickEvent.City = geoLocation.City
		clickEvent.Region = geoLocation.RegionName
		clickEvent.Lat = geoLocation.Lat
		clickEvent.Lon = geoLocation.Lon
	}

	if err := s.clickEventRepo.Create(context.Background(), clickEvent); err != nil {
		log.Printf("⚠️  Failed to track click event for shortCode %s: %v", shortCode, err)
		return nil
	}

	payloadElastic := domain.PayloadElasticClick{
		IndexType: "click_events",
		Data: domain.ClickData{
			ShortCode: shortCode,
			Metadata: domain.ClickMetaData{
				ClickedAt: newDate,
				IsBot:     ua.Bot,
				Client: domain.ClientInfo{
					IPHash: ipHash,
					Geo: domain.GeoInfo{
						CountryISOCode: clickEvent.CountryCode,
						RegionName:     clickEvent.Region,
						City:           clickEvent.City,
						Location: domain.GeoLocationElastic{
							Lat: clickEvent.Lat,
							Lon: clickEvent.Lon,
						},
					},
				},
				HTTP: domain.HTTPInfo{
					Referrer:       referrer,
					ReferrerDomain: referrerDomain,
				},
				UserAgent: domain.UserAgentInfo{
					Original: ginCtx.Request.UserAgent(),
					Device: domain.DeviceInfo{
						Name: deviceType,
					},
					Browser: domain.BrowserInfo{
						Name:    ua.Name,
						Version: ua.Version,
					},
					OS: domain.OSInfo{
						Name:    ua.OS,
						Version: ua.OSVersion,
					},
				},
			},
		},
	}

	jsonPayload, err := json.Marshal(payloadElastic)
	if err != nil {
		log.Printf("⚠️  Failed to marshal click event payload for shortCode %s: %v", shortCode, err)
		return nil
	}

	go func() {
		if err := s.rabbitmqService.PublishClickEvent(context.Background(), jsonPayload); err != nil {
			log.Printf("⚠️  Failed to publish click event to RabbitMQ for shortCode %s: %v", shortCode, err)
		}
	}()

	return nil
}

func hashIPAddress(ip string) string {
	hash := sha256.Sum256([]byte(ip))
	return hex.EncodeToString(hash[:])
}

func getDeviceType(ua useragent.UserAgent) string {
	if ua.Mobile {
		return "mobile"
	}
	if ua.Tablet {
		return "tablet"
	}
	if ua.Desktop {
		return "desktop"
	}
	return "unknown"
}

func extractDomain(url string) string {
	if url == "" {
		return ""
	}

	// Remove protocol
	url = strings.TrimPrefix(url, "http://")
	url = strings.TrimPrefix(url, "https://")

	// Get domain (before first /)
	parts := strings.Split(url, "/")
	if len(parts) > 0 {
		return parts[0]
	}

	return url
}
