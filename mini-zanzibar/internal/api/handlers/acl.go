package handlers

import (
	"mini-zanzibar/internal/database/consul"
	"mini-zanzibar/internal/database/leveldb"
	"mini-zanzibar/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ACLHandler struct {
	leveldbClient *leveldb.Client
	consulClient  *consul.Client
	logger        *zap.SugaredLogger
}

// NewACLHandler creates a new ACL handler
func NewACLHandler(leveldbClient *leveldb.Client, consulClient *consul.Client, logger *zap.SugaredLogger) *ACLHandler {
	return &ACLHandler{
		leveldbClient: leveldbClient,
		consulClient:  consulClient,
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

	// TODO: Validate request format (object:type, user:username, etc.)
	// TODO: Check if namespace exists and relation is valid
	// TODO: Implement authorization check for ACL management

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

	// TODO: Implement full authorization logic with namespace rules
	// TODO: Handle computed usersets and union operations
	// TODO: Implement caching for performance

	// For now, simple direct tuple check
	authorized, err := h.leveldbClient.CheckTuple(req.Object, req.Relation, req.User)
	if err != nil {
		h.logger.Errorw("Failed to check authorization", "error", err, "request", req)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check authorization"})
		return
	}

	response := models.ACLCheckResponse{
		Authorized: authorized,
	}

	h.logger.Infow("Authorization check", "request", req, "authorized", authorized)
	c.JSON(http.StatusOK, response)
}

// DeleteACL handles DELETE /acl - Delete an ACL tuple
func (h *ACLHandler) DeleteACL(c *gin.Context) {
	var req models.ACLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Implement authorization check for ACL management

	if err := h.leveldbClient.DeleteTuple(req.Object, req.Relation, req.User); err != nil {
		h.logger.Errorw("Failed to delete ACL tuple", "error", err, "request", req)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete ACL"})
		return
	}

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

	// TODO: Implement authorization check for ACL listing
	// TODO: Add pagination for large result sets

	tuples, err := h.leveldbClient.ListTuplesByObject(object)
	if err != nil {
		h.logger.Errorw("Failed to list ACLs by object", "error", err, "object", object)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list ACLs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tuples": tuples})
}

// ListACLsByUser handles GET /acl/user/:user - List ACLs for a user
func (h *ACLHandler) ListACLsByUser(c *gin.Context) {
	user := c.Param("user")
	if user == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user parameter is required"})
		return
	}

	// TODO: Implement authorization check for ACL listing
	// TODO: Add pagination for large result sets

	tuples, err := h.leveldbClient.ListTuplesByUser(user)
	if err != nil {
		h.logger.Errorw("Failed to list ACLs by user", "error", err, "user", user)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list ACLs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tuples": tuples})
}
