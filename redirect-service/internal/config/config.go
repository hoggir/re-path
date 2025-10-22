package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	App      AppConfig
	MongoDB  MongoDBConfig
	Redis    RedisConfig
	RabbitMQ RabbitMQConfig
	Server   ServerConfig
	CORS     CORSConfig
	JWT      JWTConfig
	Service  ServiceConfig
}

type AppConfig struct {
	Env  string
	Port string
	Name string
}

type MongoDBConfig struct {
	URI            string
	Database       string
	ConnTimeout    time.Duration
	QueryTimeout   time.Duration
	MaxPoolSize    uint64
	MinPoolSize    uint64
	DisconnTimeout time.Duration
}

type RedisConfig struct {
	Host                    string
	Port                    string
	Password                string
	DB                      int
	CacheTTL                time.Duration
	InvalidationFlagTTL     time.Duration
	ConnTimeout             time.Duration
	MaxRetries              int
	PoolSize                int
	MinIdleConns            int
}

type RabbitMQConfig struct {
	URL        string
	Queues     QueueConfig
	RPCTimeout time.Duration
}

type QueueConfig struct {
	ClickEvents      string
	DashboardRequest string
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

type JWTConfig struct {
	Secret     string
	Expiration time.Duration
	Issuer     string
}

type ServiceConfig struct {
	ClickTrackingTimeout time.Duration
	GeoIPTimeout         time.Duration
	ExternalAPITimeout   time.Duration
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	secret := getEnv("JWT_SECRET", "")
	if secret == "" {
		log.Fatal("JWT_SECRET must be set in environment variables")
	}

	return &Config{
		App: AppConfig{
			Env:  getEnv("APP_ENV", "development"),
			Port: getEnv("APP_PORT", "3011"),
			Name: getEnv("APP_NAME", "redirect-service"),
		},
		MongoDB: MongoDBConfig{
			URI:            getEnv("MONGODB_URI", "mongodb://localhost:27017"),
			Database:       getEnv("MONGODB_DATABASE", "repath"),
			ConnTimeout:    time.Duration(getEnvAsInt("MONGODB_CONN_TIMEOUT", 10)) * time.Second,
			QueryTimeout:   time.Duration(getEnvAsInt("MONGODB_QUERY_TIMEOUT", 5)) * time.Second,
			MaxPoolSize:    uint64(getEnvAsInt("MONGODB_MAX_POOL_SIZE", 100)),
			MinPoolSize:    uint64(getEnvAsInt("MONGODB_MIN_POOL_SIZE", 10)),
			DisconnTimeout: time.Duration(getEnvAsInt("MONGODB_DISCONN_TIMEOUT", 10)) * time.Second,
		},
		Redis: RedisConfig{
			Host:                getEnv("REDIS_HOST", "localhost"),
			Port:                getEnv("REDIS_PORT", "6379"),
			Password:            getEnv("REDIS_PASSWORD", ""),
			DB:                  getEnvAsInt("REDIS_DB", 0),
			CacheTTL:            time.Duration(getEnvAsInt("REDIS_CACHE_TTL", 300)) * time.Second,
			InvalidationFlagTTL: time.Duration(getEnvAsInt("REDIS_INVALIDATION_FLAG_TTL", 30)) * time.Second,
			ConnTimeout:         time.Duration(getEnvAsInt("REDIS_CONN_TIMEOUT", 5)) * time.Second,
			MaxRetries:          getEnvAsInt("REDIS_MAX_RETRIES", 3),
			PoolSize:            getEnvAsInt("REDIS_POOL_SIZE", 10),
			MinIdleConns:        getEnvAsInt("REDIS_MIN_IDLE_CONNS", 5),
		},
		RabbitMQ: RabbitMQConfig{
			URL:        getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
			RPCTimeout: time.Duration(getEnvAsInt("RABBITMQ_RPC_TIMEOUT", 5)) * time.Second,
			Queues: QueueConfig{
				ClickEvents:      getEnv("QUEUE_CLICK_EVENTS", "click_events"),
				DashboardRequest: getEnv("QUEUE_DASHBOARD_REQUEST", "dashboard_request"),
			},
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
		JWT: JWTConfig{
			Secret:     getEnv("JWT_SECRET", "your-256-bit-secret-change-this-in-production"),
			Expiration: time.Duration(getEnvAsInt("JWT_EXPIRATION_HOURS", 24)) * time.Hour,
			Issuer:     getEnv("JWT_ISSUER", "re-path-redirect-service"),
		},
		Service: ServiceConfig{
			ClickTrackingTimeout: time.Duration(getEnvAsInt("SERVICE_CLICK_TRACKING_TIMEOUT", 5)) * time.Second,
			GeoIPTimeout:         time.Duration(getEnvAsInt("SERVICE_GEOIP_TIMEOUT", 3)) * time.Second,
			ExternalAPITimeout:   time.Duration(getEnvAsInt("SERVICE_EXTERNAL_API_TIMEOUT", 10)) * time.Second,
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
