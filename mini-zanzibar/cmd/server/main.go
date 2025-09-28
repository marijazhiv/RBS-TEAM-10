package main

import (
	"log"
	"mini-zanzibar/internal/api"
	"mini-zanzibar/internal/config"
	"mini-zanzibar/internal/database/consul"
	"mini-zanzibar/internal/database/leveldb"
	"mini-zanzibar/internal/utils"
)

func main() {
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

	// Initialize API router
	router := api.NewRouter(leveldbClient, consulClient, logger, cfg)

	// Start server
	logger.Info("Starting Mini-Zanzibar server", "host", cfg.ServerHost, "port", cfg.ServerPort)
	if err := router.Run(cfg.ServerHost + ":" + cfg.ServerPort); err != nil {
		logger.Fatal("Failed to start server", err)
	}
}
