package models

// Object represents a resource in the system
type Object struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

// Subject represents a user or another object
type Subject struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

// Relation represents a relationship between objects
type Relation struct {
	Object   Object  `json:"object"`
	Relation string  `json:"relation"`
	Subject  Subject `json:"subject"`
}

// CheckRequest represents a permission check request
type CheckRequest struct {
	Object   Object  `json:"object"`
	Relation string  `json:"relation"`
	Subject  Subject `json:"subject"`
}

// CheckResponse represents the result of a permission check
type CheckResponse struct {
	Allowed bool `json:"allowed"`
}

// WriteRequest represents a request to write a relation
type WriteRequest struct {
	Relations []Relation `json:"relations"`
}
