package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	AuthServicePort = ":8081"
	ZanzibarURL     = "http://localhost:8080"
	SessionName     = "auth-session"
	JWTSecret       = "your-secret-key-change-in-production"
)

// User represents a user in our system
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password,omitempty"` // omit in responses
	Role     string `json:"role"`
}

// Demo users for testing
var demoUsers = map[string]User{
	"alice": {
		ID:       "user:alice",
		Username: "alice",
		Password: "$2a$10$N9qo8uLOickgx2ZMRZoMye.Uo0bSZ0.WllQD0J5XvYB8X1e0V0t9u", // bcrypt hash of "alice123"
		Role:     "admin",
	},
	"bob": {
		ID:       "user:bob",
		Username: "bob",
		Password: "$2a$10$N9qo8uLOickgx2ZMRZoMye.Uo0bSZ0.WllQD0J5XvYB8X1e0V0t9u", // bcrypt hash of "bob123"
		Role:     "editor",
	},
	"charlie": {
		ID:       "user:charlie",
		Username: "charlie",
		Password: "$2a$10$N9qo8uLOickgx2ZMRZoMye.Uo0bSZ0.WllQD0J5XvYB8X1e0V0t9u", // bcrypt hash of "charlie123"
		Role:     "viewer",
	},
	"david": {
		ID:       "user:david",
		Username: "david",
		Password: "$2a$10$N9qo8uLOickgx2ZMRZoMye.Uo0bSZ0.WllQD0J5XvYB8X1e0V0t9u", // bcrypt hash of "david123"
		Role:     "user",
	},
}

// LoginRequest represents login credentials
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents login response
type LoginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	User    *User  `json:"user,omitempty"`
	Token   string `json:"token,omitempty"`
}

// AuthClaims represents JWT claims
type AuthClaims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

func main() {
	// Create Gin router
	r := gin.Default()

	// CORS configuration
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000", "http://127.0.0.1:3000", "http://[::1]:3000"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"}
	config.AllowCredentials = true
	r.Use(cors.New(config))

	// Session middleware
	store := cookie.NewStore([]byte("secret-key-change-in-production"))
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   3600, // 1 hour
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
	})
	r.Use(sessions.Sessions(SessionName, store))

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy", "service": "auth-service"})
	})

	// Authentication routes
	r.POST("/auth/login", handleLogin)
	r.POST("/auth/logout", handleLogout)
	r.GET("/auth/me", authMiddleware(), handleMe)

	// Protected routes - proxy to Zanzibar with user context
	protected := r.Group("/api")
	protected.Use(authMiddleware())
	{
		// ACL management (admin only)
		protected.POST("/acl", adminOnly(), proxyToZanzibar)
		protected.DELETE("/acl", adminOnly(), proxyToZanzibar)

		// Authorization checks (all authenticated users)
		protected.GET("/acl/check", proxyToZanzibar)

		// Namespace management (admin only)
		protected.POST("/namespace", adminOnly(), proxyToZanzibar)
		protected.GET("/namespace/:id", proxyToZanzibar)
	}

	// Document operations
	r.GET("/documents", authMiddleware(), handleListDocuments)
	r.POST("/documents/:id/access", authMiddleware(), handleDocumentAccess)

	log.Printf("🔐 Auth Service starting on port %s", AuthServicePort)
	log.Printf("🔗 Proxying to Mini-Zanzibar at %s", ZanzibarURL)
	log.Fatal(r.Run(AuthServicePort))
}

// handleLogin authenticates users
func handleLogin(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, LoginResponse{
			Success: false,
			Message: "Invalid request format",
		})
		return
	}

	// Find user
	user, exists := demoUsers[req.Username]
	if !exists {
		c.JSON(http.StatusUnauthorized, LoginResponse{
			Success: false,
			Message: "Invalid credentials",
		})
		return
	}

	// For demo purposes, accept simple passwords
	validPasswords := map[string]string{
		"alice":   "alice123",
		"bob":     "bob123",
		"charlie": "charlie123",
		"david":   "david123",
	}

	if validPassword, ok := validPasswords[req.Username]; !ok || req.Password != validPassword {
		c.JSON(http.StatusUnauthorized, LoginResponse{
			Success: false,
			Message: "Invalid credentials",
		})
		return
	}

	// Create JWT token
	token, err := createJWTToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, LoginResponse{
			Success: false,
			Message: "Failed to create token",
		})
		return
	}

	// Save session
	session := sessions.Default(c)
	session.Set("user_id", user.ID)
	session.Set("username", user.Username)
	session.Set("role", user.Role)
	session.Save()

	// Remove password from response
	responseUser := user
	responseUser.Password = ""

	c.JSON(http.StatusOK, LoginResponse{
		Success: true,
		Message: "Login successful",
		User:    &responseUser,
		Token:   token,
	})
}

// handleLogout logs out the user
func handleLogout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Logout successful",
	})
}

// handleMe returns current user info
func handleMe(c *gin.Context) {
	userID := c.GetString("user_id")
	username := c.GetString("username")
	role := c.GetString("role")

	c.JSON(http.StatusOK, gin.H{
		"user_id":  userID,
		"username": username,
		"role":     role,
	})
}

// handleListDocuments returns available documents based on user role
func handleListDocuments(c *gin.Context) {
	username := c.GetString("username")
	role := c.GetString("role")

	// Demo documents based on role
	var documents []string
	switch role {
	case "admin":
		documents = []string{"report1", "manual2", "config", "logs"}
	case "editor":
		documents = []string{"report1", "manual2"}
	case "viewer":
		documents = []string{"manual2"}
	default:
		documents = []string{}
	}

	c.JSON(http.StatusOK, gin.H{
		"documents": documents,
		"user":      username,
		"role":      role,
	})
}

// handleDocumentAccess checks if user can access a document
func handleDocumentAccess(c *gin.Context) {
	docID := c.Param("id")
	userID := c.GetString("user_id")
	permission := c.DefaultQuery("permission", "viewer")

	// Forward to Zanzibar for authorization check
	checkURL := fmt.Sprintf("%s/acl/check?object=doc:%s&relation=%s&user=%s",
		ZanzibarURL, docID, permission, userID)

	resp, err := http.Get(checkURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to check authorization",
		})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to read authorization response",
		})
		return
	}

	var result map[string]interface{}
	json.Unmarshal(body, &result)

	c.JSON(resp.StatusCode, result)
}

// authMiddleware validates authentication
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userID := session.Get("user_id")
		username := session.Get("username")
		role := session.Get("role")

		if userID == nil || username == nil || role == nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authentication required",
			})
			c.Abort()
			return
		}

		// Set user context for handlers
		c.Set("user_id", userID.(string))
		c.Set("username", username.(string))
		c.Set("role", role.(string))

		c.Next()
	}
}

// adminOnly middleware restricts access to admin users
func adminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("role")
		if role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Admin access required",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// proxyToZanzibar forwards requests to Mini-Zanzibar with user context
func proxyToZanzibar(c *gin.Context) {
	userID := c.GetString("user_id")

	// Build target URL
	targetURL := ZanzibarURL + c.Request.URL.Path
	if c.Request.URL.RawQuery != "" {
		targetURL += "?" + c.Request.URL.RawQuery
	}

	// Read request body
	var body []byte
	if c.Request.Body != nil {
		var err error
		body, err = io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read request body"})
			return
		}
	}

	// Create new request
	req, err := http.NewRequest(c.Request.Method, targetURL, bytes.NewReader(body))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	// Copy headers
	for key, values := range c.Request.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	// Add user context header
	req.Header.Set("X-User-ID", userID)

	// Make request to Zanzibar
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Failed to reach authorization service"})
		return
	}
	defer resp.Body.Close()

	// Read response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
		return
	}

	// Copy response headers
	for key, values := range resp.Header {
		for _, value := range values {
			c.Header(key, value)
		}
	}

	// Return response
	c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), respBody)
}

// createJWTToken creates a JWT token for the user
func createJWTToken(user User) (string, error) {
	claims := AuthClaims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "auth-service",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(JWTSecret))
}
