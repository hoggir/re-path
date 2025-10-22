package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/hoggir/re-path/redirect-service/internal/config"
	"github.com/hoggir/re-path/redirect-service/internal/domain"
	"github.com/hoggir/re-path/redirect-service/internal/logger"
)

type GeoIPService interface {
	GetLocation(ctx context.Context, ip string) (*domain.GeoLocation, error)
}

func NewGeoIPService(
	cacheService CacheService,
	cacheKeys *CacheKeyGenerator,
	cfg *config.Config,
	log logger.Logger,
) GeoIPService {
	return &geoIPService{
		client: &http.Client{
			Timeout: cfg.Service.GeoIPTimeout,
		},
		cacheService: cacheService,
		cacheKeys:    cacheKeys,
		config:       cfg,
		logger:       log,
	}
}

type geoIPService struct {
	client       *http.Client
	cacheService CacheService
	cacheKeys    *CacheKeyGenerator
	config       *config.Config
	logger       logger.Logger
}

func (s *geoIPService) GetLocation(ctx context.Context, ip string) (*domain.GeoLocation, error) {
	if isLocalOrPrivateIP(ip) {
		s.logger.DebugContext(ctx, "IP is localhost or private, returning default location", "ip", ip)
		return &domain.GeoLocation{
			Country:     "Local",
			CountryCode: "XX",
			City:        "Localhost",
		}, nil
	}

	cacheKey := s.cacheKeys.GeoIP(ip)
	var location domain.GeoLocation
	err := s.cacheService.Get(ctx, cacheKey, &location)
	if err == nil {
		s.cacheService.RefreshTTL(ctx, cacheKey, s.config.Redis.CacheTTL)
		return &location, nil
	}

	reqCtx, cancel := context.WithTimeout(ctx, s.config.Service.GeoIPTimeout)
	defer cancel()

	url := fmt.Sprintf("http://ip-api.com/json/%s?fields=status,message,country,countryCode,region,regionName,city,zip,lat,lon,timezone,isp,org,as,query", ip)

	req, err := http.NewRequestWithContext(reqCtx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch geolocation: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("geolocation API returned status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var apiResponse struct {
		Status      string  `json:"status"`
		Message     string  `json:"message,omitempty"`
		Country     string  `json:"country"`
		CountryCode string  `json:"countryCode"`
		Region      string  `json:"region"`
		RegionName  string  `json:"regionName"`
		City        string  `json:"city"`
		Zip         string  `json:"zip"`
		Lat         float64 `json:"lat"`
		Lon         float64 `json:"lon"`
		Timezone    string  `json:"timezone"`
		ISP         string  `json:"isp"`
		Org         string  `json:"org"`
		AS          string  `json:"as"`
		Query       string  `json:"query"`
	}

	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if apiResponse.Status != "success" {
		return nil, fmt.Errorf("geolocation API error: %s", apiResponse.Message)
	}

	geoLocation := &domain.GeoLocation{
		Country:     apiResponse.Country,
		CountryCode: apiResponse.CountryCode,
		Region:      apiResponse.Region,
		RegionName:  apiResponse.RegionName,
		City:        apiResponse.City,
		Zip:         apiResponse.Zip,
		Lat:         apiResponse.Lat,
		Lon:         apiResponse.Lon,
		Timezone:    apiResponse.Timezone,
		ISP:         apiResponse.ISP,
		Org:         apiResponse.Org,
		AS:          apiResponse.AS,
		Query:       apiResponse.Query,
	}

	if err := s.cacheService.Set(ctx, cacheKey, geoLocation, s.config.Redis.CacheTTL); err != nil {
		s.logger.WarnContext(ctx, "failed to cache location for IP", "ip", ip, "error", err)
	}

	return geoLocation, nil
}

func isLocalOrPrivateIP(ip string) bool {
	// Check for localhost
	if ip == "127.0.0.1" || ip == "::1" || ip == "localhost" {
		return true
	}

	// Check for private IP ranges (simplified check)
	// 10.0.0.0/8, 172.16.0.0/12, 192.168.0.0/16
	if len(ip) >= 3 {
		if ip[:3] == "10." || ip[:8] == "192.168." || ip[:4] == "172." {
			return true
		}
	}

	return false
}
