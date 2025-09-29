package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	// Server configuration
	ServerHost string
	ServerPort string

	// Database configuration
	LevelDBPath string

	// Consul configuration
	ConsulAddress    string
	ConsulDatacenter string
	ConsulToken      string

	// Redis configuration
	RedisAddress  string
	RedisPassword string
	RedisDB       int

	// JWT configuration
	JWTSecret string
	JWTExpiry time.Duration

	// Logging configuration
	LogLevel  string
	LogFormat string

	// Security configuration
	EnableCORS        bool
	RateLimitRequests int
	RateLimitWindow   time.Duration
}

func Load() (*Config, error) {
	// Load .env file if it exists
	_ = godotenv.Load()

	cfg := &Config{
		ServerHost:        getEnvString("SERVER_HOST", "localhost"),
		ServerPort:        getEnvString("SERVER_PORT", "8080"),
		LevelDBPath:       getEnvString("LEVELDB_PATH", "./data/leveldb"),
		ConsulAddress:     getEnvString("CONSUL_ADDRESS", "localhost:8500"),
		ConsulDatacenter:  getEnvString("CONSUL_DATACENTER", "dc1"),
		ConsulToken:       getEnvString("CONSUL_TOKEN", ""),
		RedisAddress:      getEnvString("REDIS_ADDRESS", "localhost:6379"),
		RedisPassword:     getEnvString("REDIS_PASSWORD", ""),
		RedisDB:           getEnvInt("REDIS_DB", 0),
		JWTSecret:         getEnvString("JWT_SECRET", "your-secret-key-here"),
		LogLevel:          getEnvString("LOG_LEVEL", "info"),
		LogFormat:         getEnvString("LOG_FORMAT", "json"),
		EnableCORS:        getEnvBool("ENABLE_CORS", true),
		RateLimitRequests: getEnvInt("RATE_LIMIT_REQUESTS", 100),
	}

	// Parse JWT expiry
	jwtExpiryStr := getEnvString("JWT_EXPIRY", "24h")
	jwtExpiry, err := time.ParseDuration(jwtExpiryStr)
	if err != nil {
		return nil, err
	}
	cfg.JWTExpiry = jwtExpiry

	// Parse rate limit window
	rateLimitWindowStr := getEnvString("RATE_LIMIT_WINDOW", "1m")
	rateLimitWindow, err := time.ParseDuration(rateLimitWindowStr)
	if err != nil {
		return nil, err
	}
	cfg.RateLimitWindow = rateLimitWindow

	return cfg, nil
}

func getEnvString(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseBool(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}
