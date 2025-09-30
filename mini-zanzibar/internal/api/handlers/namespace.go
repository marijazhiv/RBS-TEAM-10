package handlers

import (
	"mini-zanzibar/internal/database/consul"
	"mini-zanzibar/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type NamespaceHandler struct {
	consulClient *consul.Client
	logger       *zap.SugaredLogger
}

// NewNamespaceHandler creates a new namespace handler
func NewNamespaceHandler(consulClient *consul.Client, logger *zap.SugaredLogger) *NamespaceHandler {
	return &NamespaceHandler{
		consulClient: consulClient,
		logger:       logger,
	}
}

// CreateNamespace handles POST /namespace - Create or update a namespace
func (h *NamespaceHandler) CreateNamespace(c *gin.Context) {
	var req models.NamespaceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Validate namespace configuration
	// if err := h.validateNamespaceConfig(req); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }
	// TODO: Check for circular dependencies in relations
	// TODO: Implement authorization check for namespace management

	config := consul.NamespaceConfig{
		Namespace: req.Namespace,
		Relations: convertRelationConfig(req.Relations),
	}

	if err := h.consulClient.StoreNamespace(req.Namespace, config); err != nil {
		h.logger.Errorw("Failed to store namespace", "error", err, "namespace", req.Namespace)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store namespace"})
		return
	}

	h.logger.Infow("Namespace created/updated", "namespace", req.Namespace, "version", config.Version)
	c.JSON(http.StatusCreated, gin.H{
		"message":   "Namespace created successfully",
		"namespace": req.Namespace,
		"version":   config.Version,
	})
}

// GetNamespace handles GET /namespace/:namespace - Get latest namespace configuration
func (h *NamespaceHandler) GetNamespace(c *gin.Context) {
	namespace := c.Param("namespace")
	if namespace == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "namespace parameter is required"})
		return
	}

	// TODO: Implement authorization check for namespace access

	config, err := h.consulClient.GetNamespace(namespace)
	if err != nil {
		h.logger.Errorw("Failed to get namespace", "error", err, "namespace", namespace)
		c.JSON(http.StatusNotFound, gin.H{"error": "Namespace not found"})
		return
	}

	response := models.NamespaceResponse{
		Namespace: config.Namespace,
		Relations: convertConsulRelationConfig(config.Relations),
		Version:   config.Version,
	}

	c.JSON(http.StatusOK, response)
}

// GetNamespaceVersion handles GET /namespace/:namespace/version/:version - Get specific namespace version
func (h *NamespaceHandler) GetNamespaceVersion(c *gin.Context) {
	namespace := c.Param("namespace")
	versionStr := c.Param("version")

	if namespace == "" || versionStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "namespace and version parameters are required"})
		return
	}

	version, err := strconv.Atoi(versionStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid version number"})
		return
	}

	// TODO: Implement authorization check for namespace access

	config, err := h.consulClient.GetNamespaceVersion(namespace, version)
	if err != nil {
		h.logger.Errorw("Failed to get namespace version", "error", err, "namespace", namespace, "version", version)
		c.JSON(http.StatusNotFound, gin.H{"error": "Namespace version not found"})
		return
	}

	response := models.NamespaceResponse{
		Namespace: config.Namespace,
		Relations: convertConsulRelationConfig(config.Relations),
		Version:   config.Version,
	}

	c.JSON(http.StatusOK, response)
}

// ListNamespaces handles GET /namespaces - List all namespaces
func (h *NamespaceHandler) ListNamespaces(c *gin.Context) {
	// TODO: Implement authorization check for namespace listing
	// TODO: Add pagination for large result sets

	namespaces, err := h.consulClient.ListNamespaces()
	if err != nil {
		h.logger.Errorw("Failed to list namespaces", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list namespaces"})
		return
	}

	response := models.NamespaceListResponse{
		Namespaces: namespaces,
	}

	c.JSON(http.StatusOK, response)
}

// DeleteNamespace handles DELETE /namespace/:namespace - Delete a namespace
func (h *NamespaceHandler) DeleteNamespace(c *gin.Context) {
	namespace := c.Param("namespace")
	if namespace == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "namespace parameter is required"})
		return
	}

	// TODO: Implement authorization check for namespace management
	// TODO: Check if namespace is in use before deletion

	if err := h.consulClient.DeleteNamespace(namespace); err != nil {
		h.logger.Errorw("Failed to delete namespace", "error", err, "namespace", namespace)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete namespace"})
		return
	}

	h.logger.Infow("Namespace deleted", "namespace", namespace)
	c.JSON(http.StatusOK, gin.H{"message": "Namespace deleted successfully"})
}

// convertRelationConfig converts models.RelationConfig to consul.RelationConfig
func convertRelationConfig(relations map[string]models.RelationConfig) map[string]consul.RelationConfig {
	result := make(map[string]consul.RelationConfig)
	for name, config := range relations {
		var unions []consul.UnionConfig
		for _, union := range config.Union {
			consulUnion := consul.UnionConfig{}
			if union.This != nil {
				consulUnion.This = &consul.ThisConfig{}
			}
			if union.ComputedUserset != nil {
				consulUnion.ComputedUserset = &consul.ComputedUsersetConfig{
					Relation: union.ComputedUserset.Relation,
				}
			}
			unions = append(unions, consulUnion)
		}
		result[name] = consul.RelationConfig{Union: unions}
	}
	return result
}

// convertConsulRelationConfig converts consul.RelationConfig to models.RelationConfig
func convertConsulRelationConfig(relations map[string]consul.RelationConfig) map[string]models.RelationConfig {
	result := make(map[string]models.RelationConfig)
	for name, config := range relations {
		var unions []models.UnionConfig
		for _, union := range config.Union {
			modelUnion := models.UnionConfig{}
			if union.This != nil {
				modelUnion.This = &models.ThisConfig{}
			}
			if union.ComputedUserset != nil {
				modelUnion.ComputedUserset = &models.ComputedUsersetConfig{
					Relation: union.ComputedUserset.Relation,
				}
			}
			unions = append(unions, modelUnion)
		}
		result[name] = models.RelationConfig{Union: unions}
	}
	return result
}
