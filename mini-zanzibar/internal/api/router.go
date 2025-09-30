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
	router.Use(middleware.UserContext()) // Extract user context from headers

	if cfg.EnableCORS {
		router.Use(middleware.CORS())
	}

	router.Use(middleware.RateLimit(cfg.RateLimitRequests, cfg.RateLimitWindow))

	// Initialize handlers
	aclHandler := handlers.NewACLHandler(leveldbClient, consulClient, redisClient, logger)
	namespaceHandler := handlers.NewNamespaceHandler(consulClient, logger)
	healthHandler := handlers.NewHealthHandler(logger)

	// Health check endpoint
	router.GET("/health", healthHandler.HealthCheck)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// ACL endpoints
		v1.POST("/acl", aclHandler.CreateACL)
		v1.GET("/acl/check", aclHandler.CheckACL)
		v1.DELETE("/acl", aclHandler.DeleteACL)
		v1.GET("/acl/object/:object", aclHandler.ListACLsByObject)
		v1.GET("/acl/user/:user", aclHandler.ListACLsByUser)

		// Namespace endpoints
		v1.POST("/namespace", namespaceHandler.CreateNamespace)
		v1.GET("/namespace/:namespace", namespaceHandler.GetNamespace)
		v1.GET("/namespace/:namespace/version/:version", namespaceHandler.GetNamespaceVersion)
		v1.GET("/namespaces", namespaceHandler.ListNamespaces)
		v1.DELETE("/namespace/:namespace", namespaceHandler.DeleteNamespace)
	}

	// Legacy endpoints for compatibility
	router.POST("/acl", aclHandler.CreateACL)
	router.GET("/acl/check", aclHandler.CheckACL)
	router.POST("/namespace", namespaceHandler.CreateNamespace)

	return router
}
