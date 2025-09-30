package main

import (
	"log"
	"mini-zanzibar/internal/api"
	"mini-zanzibar/internal/config"
	"mini-zanzibar/internal/database/consul"
	"mini-zanzibar/internal/database/leveldb"
	"mini-zanzibar/internal/database/redis"
	"mini-zanzibar/internal/utils"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func main() {
	// Get the absolute path to the project root
	projectRoot, err := filepath.Abs("../..")
	if err != nil {
		panic(err)
	}

	// Load .env.example file from project root
	envPath := filepath.Join(projectRoot, ".env.example")
	if err := godotenv.Load(envPath); err != nil {
		// If .env.example file doesn't exist, try .env as fallback
		envPath = filepath.Join(projectRoot, ".env")
		if err2 := godotenv.Load(envPath); err2 != nil {
			println("No .env file found, using system environment variables")
		}
	}

	// Override KEYS_FILE_PATH with absolute path
	keysPath := filepath.Join(projectRoot, "internal", "config", "api-keys.json")
	os.Setenv("KEYS_FILE_PATH", keysPath)

	// Debug: Check paths
	println("Project root:", projectRoot)
	println("KEYS_FILE_PATH:", keysPath)
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	logger := utils.InitLogger(cfg.LogLevel, cfg.LogFormat)
	defer logger.Sync()

	// Initialize LevelDB for ACL tuples
	leveldbClient, err := leveldb.NewClient(cfg.LevelDBPath)
	if err != nil {
		logger.Fatal("Failed to initialize LevelDB", err)
	}
	defer leveldbClient.Close()

	// Initialize Consul for namespace configuration
	consulClient, err := consul.NewClient(cfg.ConsulAddress, cfg.ConsulDatacenter, cfg.ConsulToken)
	if err != nil {
		logger.Fatal("Failed to initialize Consul", err)
	}

	// Initialize Redis for caching
	redisClient, err := redis.NewClient(cfg.RedisAddress, cfg.RedisPassword, cfg.RedisDB)
	if err != nil {
		logger.Fatal("failed to initialize Redis", err)
	}

	defer redisClient.Close()

	// Initialize API router
	router := api.NewRouter(leveldbClient, consulClient, redisClient, logger, cfg)

	// Start server
	logger.Info("Starting Mini-Zanzibar server", "host", cfg.ServerHost, "port", cfg.ServerPort)
	if err := router.Run(cfg.ServerHost + ":" + cfg.ServerPort); err != nil {
		logger.Fatal("Failed to start server", err)
	}
}
