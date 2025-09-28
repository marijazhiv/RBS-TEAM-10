package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Logger returns a gin.HandlerFunc that logs requests using zap
func Logger(logger *zap.SugaredLogger) gin.HandlerFunc {
	return gin.LoggerWithWriter(gin.DefaultWriter)
	// TODO: Implement custom zap-based logging middleware
	// This should log request details, response status, and timing
}

// Recovery returns a gin.HandlerFunc that recovers from panics
func Recovery(logger *zap.SugaredLogger) gin.HandlerFunc {
	return gin.RecoveryWithWriter(gin.DefaultWriter)
	// TODO: Implement custom recovery middleware that logs panics using zap
}

// CORS returns a gin.HandlerFunc that handles CORS
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// RateLimit returns a gin.HandlerFunc that implements rate limiting
func RateLimit(requests int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement rate limiting middleware
		// Should track requests per IP/user and enforce limits
		c.Next()
	}
}
