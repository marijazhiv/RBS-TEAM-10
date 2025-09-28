package models

// NamespaceRequest represents a request to create or update a namespace
type NamespaceRequest struct {
	Namespace string                    `json:"namespace" binding:"required"`
	Relations map[string]RelationConfig `json:"relations" binding:"required"`
}

// NamespaceConfig represents the complete namespace configuration
type NamespaceConfig struct {
	Namespace string                    `json:"namespace"`
	Relations map[string]RelationConfig `json:"relations"`
	Version   int                       `json:"version"`
}

// RelationConfig represents the configuration for a specific relation
type RelationConfig struct {
	Union []UnionConfig `json:"union,omitempty"`
}

// UnionConfig represents a union operation in relation configuration
type UnionConfig struct {
	This            *ThisConfig            `json:"this,omitempty"`
	ComputedUserset *ComputedUsersetConfig `json:"computed_userset,omitempty"`
}

// ThisConfig represents a direct relation
type ThisConfig struct{}

// ComputedUsersetConfig represents a computed userset based on another relation
type ComputedUsersetConfig struct {
	Relation string `json:"relation"`
}

// NamespaceResponse represents the response when retrieving a namespace
type NamespaceResponse struct {
	Namespace string                    `json:"namespace"`
	Relations map[string]RelationConfig `json:"relations"`
	Version   int                       `json:"version"`
}

// NamespaceListResponse represents the response when listing namespaces
type NamespaceListResponse struct {
	Namespaces []string `json:"namespaces"`
}
