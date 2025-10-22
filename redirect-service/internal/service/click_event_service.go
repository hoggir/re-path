package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"time"

	"github.com/hoggir/re-path/redirect-service/internal/domain"
	"github.com/hoggir/re-path/redirect-service/internal/logger"
	"github.com/hoggir/re-path/redirect-service/internal/repository"
	"github.com/mileusna/useragent"
)

type ClickEventService interface {
	TrackClick(ctx context.Context, metadata domain.ClickMetadata, shortCode string) error
}

type clickEventService struct {
	clickEventRepo  repository.ClickEventRepository
	geoIPService    GeoIPService
	redirectService RedirectService
	logger          logger.Logger
}

func NewClickEventService(
	clickEventRepo repository.ClickEventRepository,
	geoIPService GeoIPService,
	redirectService RedirectService,
	log logger.Logger,
) ClickEventService {
	return &clickEventService{
		clickEventRepo:  clickEventRepo,
		geoIPService:    geoIPService,
		redirectService: redirectService,
		logger:          log,
	}
}

func (s *clickEventService) TrackClick(ctx context.Context, metadata domain.ClickMetadata, shortCode string) error {
	newDate := time.Now().UTC()
	ua := useragent.Parse(metadata.UserAgent)
	s.logger.DebugContext(ctx, "tracking click", "ip", metadata.ClientIP, "shortCode", shortCode)

	if err := s.redirectService.IncrementClickCount(ctx, shortCode); err != nil {
		s.logger.WarnContext(ctx, "failed to increment click count", "shortCode", shortCode, "error", err)
	}

	ipHash := hashIPAddress(metadata.ClientIP)

	deviceType := getDeviceType(ua)

	referrerDomain := extractDomain(metadata.Referrer)

	geoLocation, err := s.geoIPService.GetLocation(ctx, metadata.ClientIP)
	if err != nil {
		s.logger.WarnContext(ctx, "failed to get geolocation", "ip", metadata.ClientIP, "error", err)
	}

	clickEvent := &domain.ClickEvent{
		ClickedAt:      newDate,
		ShortCode:      shortCode,
		IPAddressHash:  ipHash,
		UserAgent:      metadata.UserAgent,
		ReferrerURL:    metadata.Referrer,
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

	if err := s.clickEventRepo.Create(ctx, clickEvent); err != nil {
		s.logger.WarnContext(ctx, "failed to track click event", "shortCode", shortCode, "error", err)
		return nil
	}

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
