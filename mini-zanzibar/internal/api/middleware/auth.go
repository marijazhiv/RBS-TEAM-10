// package middleware

// import (
// 	"net/http"
// 	"os"
// 	"strings"

// 	"github.com/gin-gonic/gin"
// 	"go.uber.org/zap"
// )

// // ServiceAuthMiddleware handles service-to-service authentication
// type ServiceAuthMiddleware struct {
// 	logger       *zap.SugaredLogger
// 	serviceToken string
// }

// // NewServiceAuthMiddleware creates a new service auth middleware
// func NewServiceAuthMiddleware(logger *zap.SugaredLogger) *ServiceAuthMiddleware {
// 	serviceToken := os.Getenv("SERVICE_TOKEN")
// 	if serviceToken == "" {
// 		serviceToken = "mini-zanzibar-service-token" // default for development
// 		logger.Warn("Using default service token - set SERVICE_TOKEN environment variable in production")
// 	}

// 	return &ServiceAuthMiddleware{
// 		logger:       logger,
// 		serviceToken: serviceToken,
// 	}
// }

// // ServiceAuth protects ACL modification endpoints
// func (m *ServiceAuthMiddleware) ServiceAuth() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		token := m.extractServiceToken(c)
// 		if token != m.serviceToken {
// 			m.logger.Warnw("Invalid service token attempt",
// 				"remote_ip", c.ClientIP(),
// 				"path", c.Request.URL.Path,
// 				"method", c.Request.Method)
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid service token"})
// 			c.Abort()
// 			return
// 		}

// 		// Set service context for logging/auditing
// 		c.Set("service_authenticated", true)
// 		m.logger.Debugw("Service authenticated", "path", c.Request.URL.Path)
// 		c.Next()
// 	}
// }

// // OptionalServiceAuth for endpoints that can work with or without service auth
// func (m *ServiceAuthMiddleware) OptionalServiceAuth() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		token := m.extractServiceToken(c)
// 		if token == m.serviceToken {
// 			c.Set("service_authenticated", true)
// 			m.logger.Debugw("Service authenticated (optional)", "path", c.Request.URL.Path)
// 		}
// 		c.Next()
// 	}
// }

// func (m *ServiceAuthMiddleware) extractServiceToken(c *gin.Context) string {
// 	// 1. Check authorization header (Bearer token)
// 	authHeader := c.Request.Header.Get("Authorization")
// 	if len(authHeader) > 7 && strings.ToUpper(authHeader[0:7]) == "BEARER" {
// 		return authHeader[7:]
// 	}

// 	// 2. Check X-Service-Token header
// 	if token := c.Request.Header.Get("X-Service-Token"); token != "" {
// 		return token
// 	}

// 	// 3. Check query parameter
// 	if token := c.Query("service_token"); token != "" {
// 		return token
// 	}

// 	return ""
// }

// // IsServiceAuthenticated checks if the request is from an authenticated service
// func (m *ServiceAuthMiddleware) IsServiceAuthenticated(c *gin.Context) bool {
// 	authenticated, exists := c.Get("service_authenticated")
// 	return exists && authenticated.(bool)
// }

package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ApiKeyAuthMiddleware handles client-to-service authentication using API keys
type ApiKeyAuthMiddleware struct {
	logger    *zap.SugaredLogger
	apiKeys   map[string]string // Map to store client name to API key mapping
	keysMutex sync.RWMutex      // Protects concurrent access to apiKeys map
}

// ClientConfig represents the configuration for API clients
type ClientConfig struct {
	APIKey string `json:"apiKey"`
}

// ApiKeysConfig represents the overall API keys configuration structure
type ApiKeysConfig struct {
	Clients map[string]ClientConfig `json:"clients"`
}

// NewApiKeyAuthMiddleware creates a new API key authentication middleware
func NewApiKeyAuthMiddleware(logger *zap.SugaredLogger) (*ApiKeyAuthMiddleware, error) {
	mw := &ApiKeyAuthMiddleware{
		logger:  logger,
		apiKeys: make(map[string]string),
	}

	// Load API keys from file
	configFile := os.Getenv("KEYS_FILE_PATH")
	if configFile == "" {
		return nil, fmt.Errorf("KEYS_FILE_PATH environment variable is required")
	}

	if err := mw.LoadAPIKeys(configFile); err != nil {
		return nil, fmt.Errorf("failed to load API keys: %w", err)
	}

	logger.Info("API key authentication middleware initialized successfully",
		"clients_loaded", len(mw.apiKeys))
	return mw, nil
}

// LoadAPIKeys loads API keys from a configuration file
func (m *ApiKeyAuthMiddleware) LoadAPIKeys(configFile string) error {
	m.keysMutex.Lock()
	defer m.keysMutex.Unlock()

	// Read the configuration file
	byteValue, err := os.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("failed to read config file %s: %w", configFile, err)
	}

	// Parse JSON data
	var config ApiKeysConfig
	if err := json.Unmarshal(byteValue, &config); err != nil {
		return fmt.Errorf("failed to parse API keys config: %w", err)
	}

	// Transform into a simpler map for easier access
	m.apiKeys = make(map[string]string)
	for clientName, clientConfig := range config.Clients {
		if clientConfig.APIKey == "" {
			m.logger.Warnw("Skipping client with empty API key", "client", clientName)
			continue
		}
		m.apiKeys[clientName] = clientConfig.APIKey
	}

	m.logger.Infow("API keys loaded successfully",
		"clients_count", len(m.apiKeys),
		"config_file", configFile)
	return nil
}

// ReloadAPIKeys reloads API keys from file (useful for hot reloading)
func (m *ApiKeyAuthMiddleware) ReloadAPIKeys(configFile string) error {
	return m.LoadAPIKeys(configFile)
}

// ApiKeyAuth protects endpoints requiring API key authentication
func (m *ApiKeyAuthMiddleware) ApiKeyAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientName, apiKey := m.extractCredentials(c)

		if clientName == "" || apiKey == "" {
			m.logger.Warnw("Missing authentication credentials",
				"client_ip", c.ClientIP(),
				"path", c.Request.URL.Path,
				"method", c.Request.Method)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Client name and API key required",
				"details": "Provide Client-Name and X-API-KEY headers",
			})
			c.Abort()
			return
		}

		if !m.validateCredentials(clientName, apiKey) {
			m.logger.Warnw("Invalid authentication credentials",
				"client_ip", c.ClientIP(),
				"client_name", clientName,
				"path", c.Request.URL.Path,
				"method", c.Request.Method)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid client name or API key",
			})
			c.Abort()
			return
		}

		// Set client context for downstream handlers
		c.Set("client_authenticated", true)
		c.Set("client_name", clientName)
		c.Set("user", "client:"+clientName) // ← ADD THIS LINE - CRITICAL FIX

		m.logger.Debugw("API key authentication successful",
			"client_name", clientName,
			"path", c.Request.URL.Path)
		c.Next()
	}
}

// OptionalApiKeyAuth for endpoints that can work with or without API key auth
func (m *ApiKeyAuthMiddleware) OptionalApiKeyAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientName, apiKey := m.extractCredentials(c)

		if clientName != "" && apiKey != "" && m.validateCredentials(clientName, apiKey) {
			c.Set("client_authenticated", true)
			c.Set("client_name", clientName)
			m.logger.Debugw("Optional API key authentication successful",
				"client_name", clientName,
				"path", c.Request.URL.Path)
		} else {
			c.Set("client_authenticated", false)
		}
		c.Next()
	}
}

// extractCredentials extracts client name and API key from request
func (m *ApiKeyAuthMiddleware) extractCredentials(c *gin.Context) (string, string) {
	clientName := c.Request.Header.Get("Client-Name")
	apiKey := c.Request.Header.Get("X-API-KEY")
	return clientName, apiKey
}

// validateCredentials checks if the provided credentials are valid
func (m *ApiKeyAuthMiddleware) validateCredentials(clientName, apiKey string) bool {
	m.keysMutex.RLock()
	defer m.keysMutex.RUnlock()

	expectedAPIKey, exists := m.apiKeys[clientName]
	return exists && apiKey == expectedAPIKey
}

// IsClientAuthenticated checks if the request is from an authenticated client
func (m *ApiKeyAuthMiddleware) IsClientAuthenticated(c *gin.Context) bool {
	authenticated, exists := c.Get("client_authenticated")
	return exists && authenticated.(bool)
}

// GetClientName returns the authenticated client name from context
func (m *ApiKeyAuthMiddleware) GetClientName(c *gin.Context) string {
	clientName, exists := c.Get("client_name")
	if !exists {
		return ""
	}
	return clientName.(string)
}

// GetClientCount returns the number of registered API clients
func (m *ApiKeyAuthMiddleware) GetClientCount() int {
	m.keysMutex.RLock()
	defer m.keysMutex.RUnlock()
	return len(m.apiKeys)
}

// HasClient checks if a client name exists in the API keys
func (m *ApiKeyAuthMiddleware) HasClient(clientName string) bool {
	m.keysMutex.RLock()
	defer m.keysMutex.RUnlock()
	_, exists := m.apiKeys[clientName]
	return exists
}
