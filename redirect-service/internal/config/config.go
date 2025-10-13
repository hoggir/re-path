package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	App     AppConfig
	MongoDB MongoDBConfig
	Redis   RedisConfig
	Server  ServerConfig
	CORS    CORSConfig
}

type AppConfig struct {
	Env  string
	Port string
	Name string
}

type MongoDBConfig struct {
	URI      string
	Database string
	Timeout  time.Duration
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
	CacheTTL time.Duration
}

type ServerConfig struct {
	GinMode        string
	TrustedProxies []string
}

type CORSConfig struct {
	AllowOrigins string
	AllowMethods string
	AllowHeaders string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	return &Config{
		App: AppConfig{
			Env:  getEnv("APP_ENV", "development"),
			Port: getEnv("APP_PORT", "3011"),
			Name: getEnv("APP_NAME", "redirect-service"),
		},
		MongoDB: MongoDBConfig{
			URI:      getEnv("MONGODB_URI", "mongodb://localhost:27017"),
			Database: getEnv("MONGODB_DATABASE", "repath"),
			Timeout:  time.Duration(getEnvAsInt("MONGODB_TIMEOUT", 10)) * time.Second,
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
			CacheTTL: time.Duration(getEnvAsInt("REDIS_CACHE_TTL", 3600)) * time.Second,
		},
		Server: ServerConfig{
			GinMode:        getEnv("GIN_MODE", "debug"),
			TrustedProxies: []string{getEnv("TRUSTED_PROXIES", "127.0.0.1")},
		},
		CORS: CORSConfig{
			AllowOrigins: getEnv("CORS_ALLOW_ORIGINS", "*"),
			AllowMethods: getEnv("CORS_ALLOW_METHODS", "GET,POST,PUT,DELETE,OPTIONS"),
			AllowHeaders: getEnv("CORS_ALLOW_HEADERS", "Origin,Content-Type,Accept,Authorization"),
		},
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return fallback
}
