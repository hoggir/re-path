package service

import (
	"fmt"

	"github.com/hoggir/re-path/redirect-service/internal/config"
)

type CacheKeyGenerator struct {
	prefix string
}

func NewCacheKeyGenerator(cfg *config.Config) *CacheKeyGenerator {
	prefix := cfg.App.Name
	if prefix == "" {
		prefix = "repath"
	}
	return &CacheKeyGenerator{
		prefix: prefix,
	}
}

func (g *CacheKeyGenerator) URL(shortCode string) string {
	return fmt.Sprintf("%s:url:%s", g.prefix, shortCode)
}

func (g *CacheKeyGenerator) Dashboard(userID int) string {
	return fmt.Sprintf("%s:dashboard:%d", g.prefix, userID)
}

func (g *CacheKeyGenerator) GeoIP(ip string) string {
	return fmt.Sprintf("%s:geoip:%s", g.prefix, ip)
}

func (g *CacheKeyGenerator) DashboardInvalidationFlag(userID int) string {
	return fmt.Sprintf("%s:dashboard_invalid:%d", g.prefix, userID)
}
