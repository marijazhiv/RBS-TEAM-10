package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// rateLimiter implements a simple token bucket rate limiter
type rateLimiter struct {
	mu          sync.Mutex
	ips         map[string]*rateLimit
	requests    int
	window      time.Duration
	cleanupTime time.Duration
}

type rateLimit struct {
	count     int
	lastReset time.Time
}

// Logger returns a gin.HandlerFunc that logs requests using zap
func Logger(logger *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Log after request is processed
		timestamp := time.Now()
		latency := timestamp.Sub(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()

		// Log with appropriate level based on status code
		logFields := []interface{}{
			"status", statusCode,
			"method", method,
			"path", path,
			"query", query,
			"ip", clientIP,
			"latency", latency,
			"user-agent", c.Request.UserAgent(),
		}

		// Add error message if present
		if errorMessage != "" {
			logFields = append(logFields, "error", errorMessage)
		}

		switch {
		case statusCode >= 500:
			logger.Errorw("Server error", logFields...)
		case statusCode >= 400:
			logger.Warnw("Client error", logFields...)
		case statusCode >= 300:
			logger.Infow("Redirection", logFields...)
		default:
			logger.Infow("Request completed", logFields...)
		}
	}
}

// Recovery returns a gin.HandlerFunc that recovers from panics
func Recovery(logger *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for broken connection (can't write response)
				var brokenPipe bool
				if ne, ok := err.(*http.ProtocolError); ok {
					if strings.Contains(strings.ToLower(ne.Error()), "broken pipe") ||
						strings.Contains(strings.ToLower(ne.Error()), "connection reset by peer") {
						brokenPipe = true
					}
				}

				// Log the panic
				clientIP := c.ClientIP()
				path := c.Request.URL.Path
				method := c.Request.Method

				logger.Errorw("Panic recovered",
					"error", err,
					"path", path,
					"method", method,
					"ip", clientIP,
					"user-agent", c.Request.UserAgent(),
					"stack", getStacktrace(),
				)

				// If the connection is dead, we can't write a status
				if brokenPipe {
					c.Error(fmt.Errorf("connection broken: %v", err))
					c.Abort()
				} else {
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
						"error":   "Internal server error",
						"message": "Something went wrong. Please try again later.",
					})
				}
			}
		}()

		c.Next()
	}
}

// getStacktrace returns a formatted stacktrace
func getStacktrace() string {
	// This is a simplified version. In production, you might want to use
	// runtime.Stack() to get the actual stack trace
	return "stack trace unavailable in production mode"
}

// CORS returns a gin.HandlerFunc that handles CORS
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin == "" {
			origin = "*"
		}

		c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Accept, Origin, Cache-Control, X-Requested-With, Client-Name, X-API-KEY")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400") // 24 hours

		// Handle preflight requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// RateLimit returns a gin.HandlerFunc that implements rate limiting
func RateLimit(requests int, window time.Duration) gin.HandlerFunc {
	limiter := &rateLimiter{
		ips:         make(map[string]*rateLimit),
		requests:    requests,
		window:      window,
		cleanupTime: window * 2, // Clean up old entries after 2 windows
	}

	// Start background cleanup goroutine
	go limiter.cleanupOldEntries()

	return func(c *gin.Context) {
		clientIP := c.ClientIP()

		// Skip rate limiting for health checks
		if c.Request.URL.Path == "/health" {
			c.Next()
			return
		}

		if !limiter.allow(clientIP) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":       "Rate limit exceeded",
				"message":     fmt.Sprintf("Maximum %d requests per %v allowed", requests, window),
				"retry_after": window.String(),
			})
			c.Abort()
			return
		}

		// Add rate limit headers
		c.Writer.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", requests))
		c.Writer.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", limiter.remaining(clientIP)))
		c.Writer.Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(window).Unix()))

		c.Next()
	}
}

// allow checks if the IP is allowed to make a request
func (rl *rateLimiter) allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	limit, exists := rl.ips[ip]

	if !exists {
		rl.ips[ip] = &rateLimit{
			count:     1,
			lastReset: now,
		}
		return true
	}

	// Reset counter if window has passed
	if now.Sub(limit.lastReset) > rl.window {
		limit.count = 1
		limit.lastReset = now
		return true
	}

	// Check if under limit
	if limit.count < rl.requests {
		limit.count++
		return true
	}

	return false
}

// remaining returns the number of remaining requests for an IP
func (rl *rateLimiter) remaining(ip string) int {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limit, exists := rl.ips[ip]
	if !exists {
		return rl.requests
	}

	// Reset if window has passed
	if time.Now().Sub(limit.lastReset) > rl.window {
		return rl.requests
	}

	remaining := rl.requests - limit.count
	if remaining < 0 {
		return 0
	}
	return remaining
}

// cleanupOldEntries removes old rate limit entries to prevent memory leaks
func (rl *rateLimiter) cleanupOldEntries() {
	ticker := time.NewTicker(rl.cleanupTime)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for ip, limit := range rl.ips {
			if now.Sub(limit.lastReset) > rl.cleanupTime {
				delete(rl.ips, ip)
			}
		}
		rl.mu.Unlock()
	}
}

// AuthUser extracts and sets user information from authentication
func AuthUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract user from authentication context (set by your auth middleware)
		// This is a placeholder - adapt to your actual authentication system
		clientName, exists := c.Get("client_name")
		if exists {
			c.Set("user", fmt.Sprintf("client:%s", clientName))
		}

		c.Next()
	}
}

// SecurityHeaders adds security-related headers to responses
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Security headers
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
		c.Writer.Header().Set("X-Frame-Options", "DENY")
		c.Writer.Header().Set("X-XSS-Protection", "1; mode=block")
		c.Writer.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		c.Writer.Header().Set("Content-Security-Policy", "default-src 'self'")
		c.Writer.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

		c.Next()
	}
}

func DebugMiddleware(logger *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Log the request details and headers
		logger.Debugw("Incoming request",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"headers", c.Request.Header,
			"client_ip", c.ClientIP())

		// Log before processing the request
		logger.Debugw("Context before processing",
			"client_name", c.GetString("client_name"),
			"user", c.GetString("user"),
			"client_authenticated", c.GetBool("client_authenticated"))

		c.Next()

		// Log after processing the request
		logger.Debugw("Context after processing",
			"client_name", c.GetString("client_name"),
			"user", c.GetString("user"),
			"client_authenticated", c.GetBool("client_authenticated"),
			"status", c.Writer.Status())
	}
}
