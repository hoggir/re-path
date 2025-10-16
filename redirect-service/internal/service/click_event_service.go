package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/hoggir/re-path/redirect-service/internal/domain"
	"github.com/hoggir/re-path/redirect-service/internal/repository"
	"github.com/mileusna/useragent"
)

type ClickEventService interface {
	TrackClick(ctx context.Context, ginCtx *gin.Context, shortCode string) error
}

type clickEventService struct {
	clickEventRepo repository.ClickEventRepository
	geoIPService   GeoIPService
}

func NewClickEventService(
	clickEventRepo repository.ClickEventRepository,
	geoIPService GeoIPService,
) ClickEventService {
	return &clickEventService{
		clickEventRepo: clickEventRepo,
		geoIPService:   geoIPService,
	}
}

func (s *clickEventService) TrackClick(ctx context.Context, ginCtx *gin.Context, shortCode string) error {
	ua := useragent.Parse(ginCtx.Request.UserAgent())
	clientIP := ginCtx.ClientIP()
	log.Printf("📍 Tracking click from IP: %s", clientIP)

	ipHash := hashIPAddress(clientIP)

	deviceType := getDeviceType(ua)

	referrer := ginCtx.Request.Referer()

	// geoLocation, err := s.geoIPService.GetLocation(ctx, clientIP)
	geoLocation, err := s.geoIPService.GetLocation(ctx, "203.175.11.126")
	if err != nil {
		log.Printf("⚠️  Failed to get geolocation for IP %s: %v", clientIP, err)
		// Continue without geolocation data
	}

	clickEvent := &domain.ClickEvent{
		ShortCode:      shortCode,
		IPAddressHash:  ipHash,
		UserAgent:      ginCtx.Request.UserAgent(),
		ReferrerURL:    referrer,
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
	}

	if err := s.clickEventRepo.Create(context.Background(), clickEvent); err != nil {
		log.Printf("⚠️  Failed to track click event for shortCode %s: %v", shortCode, err)
	} else {
		locationStr := "Unknown"
		if geoLocation != nil {
			locationStr = fmt.Sprintf("%s, %s (%s)", geoLocation.City, geoLocation.Country, geoLocation.CountryCode)
		}
		log.Printf("📊 Click event tracked for shortCode: %s | Location: %s | Device: %s | Browser: %s",
			shortCode, locationStr, deviceType, ua.Name)
	}

	return nil
}

// hashIPAddress hashes IP address using SHA256 for privacy
func hashIPAddress(ip string) string {
	hash := sha256.Sum256([]byte(ip))
	return hex.EncodeToString(hash[:])
}

// getDeviceType determines device type from user agent
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
