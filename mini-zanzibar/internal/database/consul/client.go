package consul

import (
	"encoding/json"
	"fmt"
	"path"
	"strings"

	"github.com/hashicorp/consul/api"
)

type Client struct {
	client *api.Client
}

type NamespaceConfig struct {
	Namespace string                    `json:"namespace"`
	Relations map[string]RelationConfig `json:"relations"`
	Version   int                       `json:"version"`
}

type RelationConfig struct {
	Union []UnionConfig `json:"union,omitempty"`
}

type UnionConfig struct {
	This            *ThisConfig            `json:"this,omitempty"`
	ComputedUserset *ComputedUsersetConfig `json:"computed_userset,omitempty"`
}

type ThisConfig struct{}

type ComputedUsersetConfig struct {
	Relation string `json:"relation"`
}

// NewClient creates a new Consul client
func NewClient(address, datacenter, token string) (*Client, error) {
	config := api.DefaultConfig()
	config.Address = address
	config.Datacenter = datacenter
	if token != "" {
		config.Token = token
	}

	client, err := api.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Consul client: %w", err)
	}

	return &Client{
		client: client,
	}, nil
}

// StoreNamespace stores a namespace configuration with versioning
func (c *Client) StoreNamespace(namespace string, config NamespaceConfig) error {
	// Get current version
	currentVersion, err := c.getLatestVersion(namespace)
	if err != nil {
		return fmt.Errorf("failed to get current version: %w", err)
	}

	// Increment version
	config.Version = currentVersion + 1

	// Marshal configuration
	data, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal namespace config: %w", err)
	}

	// Store versioned configuration
	key := c.getVersionedKey(namespace, config.Version)
	kv := c.client.KV()

	pair := &api.KVPair{
		Key:   key,
		Value: data,
	}

	_, err = kv.Put(pair, nil)
	if err != nil {
		return fmt.Errorf("failed to store namespace config: %w", err)
	}

	// Update latest pointer
	latestKey := c.getLatestKey(namespace)
	latestPair := &api.KVPair{
		Key:   latestKey,
		Value: []byte(fmt.Sprintf("%d", config.Version)),
	}

	_, err = kv.Put(latestPair, nil)
	if err != nil {
		return fmt.Errorf("failed to update latest version pointer: %w", err)
	}

	return nil
}

// GetNamespace retrieves the latest namespace configuration
func (c *Client) GetNamespace(namespace string) (*NamespaceConfig, error) {
	// Get latest version
	version, err := c.getLatestVersion(namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest version: %w", err)
	}

	if version == 0 {
		return nil, fmt.Errorf("namespace not found: %s", namespace)
	}

	return c.GetNamespaceVersion(namespace, version)
}

// GetNamespaceVersion retrieves a specific version of namespace configuration
func (c *Client) GetNamespaceVersion(namespace string, version int) (*NamespaceConfig, error) {
	key := c.getVersionedKey(namespace, version)
	kv := c.client.KV()

	pair, _, err := kv.Get(key, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get namespace config: %w", err)
	}

	if pair == nil {
		return nil, fmt.Errorf("namespace version not found: %s v%d", namespace, version)
	}

	var config NamespaceConfig
	if err := json.Unmarshal(pair.Value, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal namespace config: %w", err)
	}

	return &config, nil
}

// ListNamespaces returns all available namespaces
func (c *Client) ListNamespaces() ([]string, error) {
	prefix := "zanzibar/namespaces/"
	kv := c.client.KV()

	pairs, _, err := kv.List(prefix, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list namespaces: %w", err)
	}

	namespaceSet := make(map[string]bool)
	for _, pair := range pairs {
		// Extract namespace from key path
		relativePath := pair.Key[len(prefix):]
		parts := strings.Split(relativePath, "/")
		if len(parts) > 0 && parts[0] != "" {
			namespaceSet[parts[0]] = true
		}
	}

	var namespaces []string
	for namespace := range namespaceSet {
		namespaces = append(namespaces, namespace)
	}

	return namespaces, nil
}

// DeleteNamespace removes all versions of a namespace
func (c *Client) DeleteNamespace(namespace string) error {
	prefix := path.Join("zanzibar/namespaces", namespace)
	kv := c.client.KV()

	_, err := kv.DeleteTree(prefix, nil)
	if err != nil {
		return fmt.Errorf("failed to delete namespace: %w", err)
	}

	return nil
}

// getLatestVersion retrieves the latest version number for a namespace
func (c *Client) getLatestVersion(namespace string) (int, error) {
	key := c.getLatestKey(namespace)
	kv := c.client.KV()

	pair, _, err := kv.Get(key, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to get latest version: %w", err)
	}

	if pair == nil {
		return 0, nil // No versions exist yet
	}

	var version int
	if _, err := fmt.Sscanf(string(pair.Value), "%d", &version); err != nil {
		return 0, fmt.Errorf("failed to parse version: %w", err)
	}

	return version, nil
}

// getVersionedKey returns the Consul key for a specific namespace version
func (c *Client) getVersionedKey(namespace string, version int) string {
	return path.Join("zanzibar/namespaces", namespace, "versions", fmt.Sprintf("%d", version))
}

// getLatestKey returns the Consul key for the latest version pointer
func (c *Client) getLatestKey(namespace string) string {
	return path.Join("zanzibar/namespaces", namespace, "latest")
}
