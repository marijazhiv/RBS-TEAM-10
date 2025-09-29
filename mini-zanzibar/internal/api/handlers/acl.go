package handlers

import (
	"fmt"
	"mini-zanzibar/internal/database/consul"
	"mini-zanzibar/internal/database/leveldb"
	"mini-zanzibar/internal/database/redis"
	"mini-zanzibar/internal/models"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ACLHandler struct {
	leveldbClient *leveldb.Client
	consulClient  *consul.Client
	redisClient   *redis.Client
	logger        *zap.SugaredLogger
}

// NewACLHandler creates a new ACL handler
func NewACLHandler(leveldbClient *leveldb.Client, consulClient *consul.Client, redisClient *redis.Client, logger *zap.SugaredLogger) *ACLHandler {
	return &ACLHandler{
		leveldbClient: leveldbClient,
		consulClient:  consulClient,
		redisClient:   redisClient,
		logger:        logger,
	}
}

// CreateACL handles POST /acl - Create or update an ACL tuple
func (h *ACLHandler) CreateACL(c *gin.Context) {
	var req models.ACLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate request format (object:type, user:username, etc.)
	if err := h.validateACLRequest(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if namespace exists and relation is valid
	if err := h.validateNamespaceAndRelation(req.Object, req.Relation); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	// Implement authorization check for ACL management
	if !h.isAuthorizedForACLManagement(c, req.Object) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unauthorized to manage ACLs for this object"})
		return
	}

	tuple := leveldb.ACLTuple{
		Object:   req.Object,
		Relation: req.Relation,
		User:     req.User,
	}

	if err := h.leveldbClient.StoreTuple(tuple); err != nil {
		h.logger.Errorw("Failed to store ACL tuple", "error", err, "tuple", tuple)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store ACL"})
		return
	}

	// Invalidate cache for this object and user
	h.invalidateAuthorizationCache(req.Object, req.Relation, req.User)

	h.logger.Infow("ACL tuple created", "tuple", tuple)
	c.JSON(http.StatusCreated, gin.H{"message": "ACL created successfully"})
}

// CheckACL handles GET /acl/check - Check authorization
func (h *ACLHandler) CheckACL(c *gin.Context) {
	var req models.ACLCheckRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate required parameters
	if req.Object == "" || req.Relation == "" || req.User == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "object, relation and user parameters are required"})
		return
	}

	// Caching for performance
	cacheKey := fmt.Sprintf("auth:%s:%s:%s", req.Object, req.Relation, req.User)

	// Try to get from cache first
	if cached, found := h.redisClient.Get(cacheKey); found {
		authorized, ok := cached.(bool)
		if ok {
			response := models.ACLCheckResponse{
				Authorized: authorized,
			}
			h.logger.Debugw("authorization cache hit", "request", req, "authorized", authorized)
			c.JSON(http.StatusOK, response)
			return
		}
	}

	// Implement full authorization logic with namespace rules
	// Handle computed usersets and union operations
	authorized, err := h.performAuthorizationCheck(req.Object, req.Relation, req.User)
	if err != nil {
		h.logger.Errorw("failed to check authorization", "error", err, "request", req)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check authorization"})
		return
	}

	// Cache the result with 5-minute expiration
	h.redisClient.Set(cacheKey, authorized, 5*time.Minute)

	response := models.ACLCheckResponse{
		Authorized: authorized,
	}

	h.logger.Infow("Authorization check", "request", req, "authorized", authorized, "cache_miss", true)
	c.JSON(http.StatusOK, response)
}

// DeleteACL handles DELETE /acl - Delete an ACL tuple
func (h *ACLHandler) DeleteACL(c *gin.Context) {
	var req models.ACLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.validateACLRequest(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	// Implement authorization check for ACL management
	if !h.isAuthorizedForACLManagement(c, req.Object) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized to manage ACLs for this object"})
		return
	}

	if err := h.leveldbClient.DeleteTuple(req.Object, req.Relation, req.User); err != nil {
		h.logger.Errorw("Failed to delete ACL tuple", "error", err, "request", req)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete ACL"})
		return
	}

	// Invalidate cache for this object and user
	h.invalidateAuthorizationCache(req.Object, req.Relation, req.User)

	h.logger.Infow("ACL tuple deleted", "request", req)
	c.JSON(http.StatusOK, gin.H{"message": "ACL deleted successfully"})
}

// ListACLsByObject handles GET /acl/object/:object - List ACLs for an object
func (h *ACLHandler) ListACLsByObject(c *gin.Context) {
	object := c.Param("object")
	if object == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "object parameter is required"})
		return
	}

	// Authorization check for ACL listing
	if !h.isAuthorizedForACLListing(c, object) {
		c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized to list ACLs for this object"})
		return
	}
	// Add pagination for large result sets
	page, pageSize := h.getPaginationParams(c)

	// Use paginated version
	tuples, total, err := h.leveldbClient.ListTuplesByObjectPagination(object, page, pageSize)
	if err != nil {
		h.logger.Errorw("failed to list ACLs by object", "error", err, "object", object)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list ACLs"})
		return
	}

	response := gin.H{
		"tuples": tuples,
		"pagination": gin.H{
			"page":      page,
			"page_size": pageSize,
			"total":     total,
			"has_more":  len(tuples) == pageSize,
		},
	}

	c.JSON(http.StatusOK, response)
}

// ListACLsByUser handles GET /acl/user/:user - List ACLs for a user
func (h *ACLHandler) ListACLsByUser(c *gin.Context) {
	user := c.Param("user")
	if user == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user parameter is required"})
		return
	}

	// Authorization check for ACL listing
	if !h.isAuthorizedForACLListing(c, user) {
		c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized to list ACLs for this user"})
		return
	}
	// Pagination for large result sets
	page, pageSize := h.getPaginationParams(c)

	// Use paginated version
	tuples, total, err := h.leveldbClient.ListTuplesByUserPagination(user, page, pageSize)

	if err != nil {
		h.logger.Errorw("failed to list ACLs by user", "error", err, "user", user)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list ACLs"})
		return
	}

	response := gin.H{
		"tuples": tuples,
		"pagination": gin.H{
			"page":      page,
			"page_size": pageSize,
			"total":     total,
			"has_more":  len(tuples) == pageSize,
		},
	}

	c.JSON(http.StatusOK, response)
}

func (h *ACLHandler) validateACLRequest(req models.ACLRequest) error {
	if req.Object == "" {
		return fmt.Errorf("object is required")
	}

	if req.Relation == "" {
		return fmt.Errorf("relation is required")
	}

	if req.User == "" {
		return fmt.Errorf("user is required")
	}

	parts := strings.Split(req.Object, ":")
	if len(parts) != 2 {
		return fmt.Errorf("object must be in format 'namespace:object_id'")
	}

	userParts := strings.Split(req.User, ":")
	if len(userParts) < 2 {
		return fmt.Errorf("user must be in format 'user_type:user_id' or 'userset:namespace:relation'")
	}

	return nil
}

func (h *ACLHandler) validateNamespaceAndRelation(object string, relation string) error {
	parts := strings.Split(object, ":")
	if len(parts) != 2 {
		return fmt.Errorf("invalid object format")
	}
	namespace := parts[0]

	// Check if namespace exists in Consul
	exists, err := h.consulClient.NamespaceExists(namespace)
	if err != nil {
		return fmt.Errorf("failed to check namespace: %v", err)
	}

	if !exists {
		return fmt.Errorf("namespace '%s' does not exist", namespace)
	}

	// Check if relation is valid for this namespace
	valid, err := h.consulClient.RelationExists(namespace, relation)
	if err != nil {
		return fmt.Errorf("failed to check relation: %v", err)
	}

	if !valid {
		return fmt.Errorf("relation '%s' is not valid for namespace '%s'", relation, namespace)
	}

	return nil
}

// Authorization check for ACL management
func (h *ACLHandler) isAuthorizedForACLManagement(c *gin.Context, object string) bool {
	// Extract user from context (set by authentication middleware)
	user, exists := c.Get("user")
	if !exists {
		h.logger.Warnw("No user in context for ACL management", "object", object)
		return false
	}

	// Check if user has admin rights for this object's namespace
	parts := strings.Split(object, ":")
	if len(parts) != 2 {
		return false
	}
	namespace := parts[0]

	// Check if user is namespace admin
	authorized, err := h.performAuthorizationCheck(
		fmt.Sprintf("namespace:%s", namespace),
		"admin",
		user.(string),
	)

	if err != nil {
		h.logger.Errorw("Failed to check admin authorization", "error", err, "namespace", namespace, "user", user)
		return false
	}

	return authorized
}

// Full authorization logic with namespace rules
// Handle computed usersets and union operations
func (h *ACLHandler) performAuthorizationCheck(object, relation, user string) (bool, error) {
	// 1. Check direct tuple first (common case)
	directAuthorized, err := h.leveldbClient.CheckTuple(object, relation, user)

	if err != nil {
		return false, fmt.Errorf("failed to check direct tuple: %v", err)
	}

	if directAuthorized {
		return true, nil
	}

	// 2. Check namespace rules for computed usersets and union operations
	parts := strings.Split(object, ":")
	if len(parts) != 2 {
		return false, nil
	}
	namespace := parts[0]

	// Get namespace configuration to check for computed usersets
	config, err := h.consulClient.GetNamespace(namespace)
	if err != nil {
		return false, fmt.Errorf("failed to get namespace config: %v", err)
	}

	// Check if this relation has computed usersets defined
	if relationConfig, exists := config.Relations[relation]; exists {
		for _, union := range relationConfig.Union {
			// Handle computed userset
			if union.ComputedUserset != nil {
				computedRelation := union.ComputedUserset.Relation

				authorized, err := h.checkComputedUserset(object, computedRelation, user)

				if err != nil {
					return false, err
				}

				if authorized {
					return true, nil
				}
			}

			// Handle union of multiple relations (this is a simplified version)
			// In a full implementation, you'd handle union, intersection, and exclusion
			if union.This != nil {
				// This represents direct membership - already checked above
				continue
			}
		}
	}

	return false, nil
}

// checkComputedUserset handles computed userset logic
func (h *ACLHandler) checkComputedUserset(object, computedRelation, user string) (bool, error) {
	// For computed userset, we need to find all objects that this object has the computed relation to
	// and check if the user has the base relation to those objects

	// This is a simplified implementation - in a real Zanzibar-like system,
	// you'd have a more complex graph traversal

	// Example: If document:1 has editor → group:eng and user is member of group:eng
	// Then user is effectively an editor of document:1

	// Get all tuples where this object has the computed relation
	relatedTuples, err := h.leveldbClient.ListTuplesByObjectAndRelation(object, computedRelation)
	if err != nil {
		return false, err
	}

	// Check if user has the base relation to any of the related objects
	for _, tuple := range relatedTuples {
		// tuple.User is the target object in the computed relation
		// Check if our user has the base relation to that target object
		authorized, err := h.performAuthorizationCheck(tuple.User, computedRelation, user)

		if err != nil {
			return false, err
		}

		if authorized {
			return true, nil
		}
	}

	return false, nil
}

// Invalidate authorization cache when ACLs change
func (h *ACLHandler) invalidateAuthorizationCache(object, relation, user string) {

	// Invalidate specific check
	h.redisClient.Delete(fmt.Sprintf("auth:%s:%s:%s", object, relation, user))

	// Invalidate pattern-based cache using Redis SCAN
	h.redisClient.DeletePattern(fmt.Sprintf("auth:%s:%s:*", object, relation))

	h.redisClient.DeletePattern(fmt.Sprintf("auth:%s:*", object))

	h.logger.Debugw("Invalidated authorization cache", "object", object, "relation", relation, "user", user)
}

// Authorization check for ACL listing
func (h *ACLHandler) isAuthorizedForACLListing(c *gin.Context, resource string) bool {
	// Extract user from context (set by authentication middleware)
	user, exists := c.Get("user")
	if !exists {
		h.logger.Warnw("no user in context for ACL listing", "resource", resource)
		return false
	}

	// For objects, check if user has view permissions
	if strings.Contains(resource, ":") {
		parts := strings.Split(resource, ":")
		if len(parts) == 2 {
			namespace := parts[0]

			// Check if user can view this namespace's ACLs
			authorized, err := h.performAuthorizationCheck(
				fmt.Sprintf("namespace:%s", namespace),
				"view_acls",
				user.(string),
			)
			if err == nil && authorized {
				return true
			}
		}
	}

	// Allow users to view their own ACLs
	return resource == user.(string)
}

// getPaginationParams extracts pagination parameters from query
func (h *ACLHandler) getPaginationParams(c *gin.Context) (int, int) {
	page := 1
	pageSize := 50

	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 {
			pageSize = ps
		}
	}

	// Limit page size to prevent excessive memory usage
	if pageSize > 1000 {
		pageSize = 1000
	}

	if pageSize < 1 {
		pageSize = 1
	}

	return page, pageSize
}
