package models

// ACLRequest represents a request to create or update an ACL tuple
type ACLRequest struct {
	Object   string `json:"object" binding:"required"`
	Relation string `json:"relation" binding:"required"`
	User     string `json:"user" binding:"required"`
}

// ACLCheckRequest represents a request to check authorization
type ACLCheckRequest struct {
	Object   string `json:"object" form:"object" binding:"required"`
	Relation string `json:"relation" form:"relation" binding:"required"`
	User     string `json:"user" form:"user" binding:"required"`
}

// ACLCheckResponse represents the response for an authorization check
type ACLCheckResponse struct {
	Authorized bool `json:"authorized"`
}

// ACLTuple represents an ACL tuple stored in the database
type ACLTuple struct {
	Object   string `json:"object"`
	Relation string `json:"relation"`
	User     string `json:"user"`
}
