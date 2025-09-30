package api

import (
	"mini-zanzibar/internal/api/handlers"
	"mini-zanzibar/internal/api/middleware"
	"mini-zanzibar/internal/config"
	"mini-zanzibar/internal/database/consul"
	"mini-zanzibar/internal/database/leveldb"
	"mini-zanzibar/internal/database/redis"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// NewRouter creates and configures the API router
func NewRouter(leveldbClient *leveldb.Client, consulClient *consul.Client, redisClient *redis.Client, logger *zap.SugaredLogger, cfg *config.Config) *gin.Engine {

	// Set Gin mode based on environment
	if cfg.LogLevel != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Add middleware
	router.Use(middleware.Logger(logger))
	router.Use(middleware.Recovery(logger))
	router.Use(middleware.SecurityHeaders()) // Add security headers
	router.Use(middleware.AuthUser())
	router.Use(middleware.DebugMiddleware(logger))

	if cfg.EnableCORS {
		router.Use(middleware.CORS())
	}

	router.Use(middleware.RateLimit(cfg.RateLimitRequests, cfg.RateLimitWindow))

	// Initialize handlers and auth middleware
	aclHandler := handlers.NewACLHandler(leveldbClient, consulClient, redisClient, logger)
	namespaceHandler := handlers.NewNamespaceHandler(consulClient, logger)
	healthHandler := handlers.NewHealthHandler(logger)

	// Initialize API key auth middleware
	apiKeyAuth, err := middleware.NewApiKeyAuthMiddleware(logger)
	if err != nil {
		logger.Fatalw("Failed to initialize API key authentication", "error", err)
	}

	// Health check endpoint (public)
	router.GET("/health", healthHandler.HealthCheck)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// === PUBLIC ENDPOINTS ===
		// Anyone can check authorization
		v1.GET("/acl/check", aclHandler.CheckACL)

		// === CLIENT-PROTECTED ENDPOINTS ===
		// These require valid API key authentication
		clientRoutes := v1.Group("")
		clientRoutes.Use(apiKeyAuth.ApiKeyAuth())
		{
			// ACL Management
			v1.POST("/acl", aclHandler.CreateACL)
			v1.DELETE("/acl", aclHandler.DeleteACL)
			v1.GET("/acl/object/:object", aclHandler.ListACLsByObject)
			v1.GET("/acl/user/:user", aclHandler.ListACLsByUser)

			// Namespace Management
			v1.POST("/namespace", namespaceHandler.CreateNamespace)
			v1.DELETE("/namespace/:namespace", namespaceHandler.DeleteNamespace)
			v1.GET("/namespace/:namespace", namespaceHandler.GetNamespace)
			v1.GET("/namespace/:namespace/version/:version", namespaceHandler.GetNamespaceVersion)
			v1.GET("/namespaces", namespaceHandler.ListNamespaces)
		}

	}

	// Legacy endpoints for compatibility
	legacy := router.Group("")
	legacy.Use(apiKeyAuth.ApiKeyAuth()) // Protect legacy endpoints with API key auth
	{
		router.POST("/acl", aclHandler.CreateACL)
		router.POST("/namespace", namespaceHandler.CreateNamespace)
	}

	// Legacy public endpoints
	router.GET("/acl/check", aclHandler.CheckACL)

	return router
}
