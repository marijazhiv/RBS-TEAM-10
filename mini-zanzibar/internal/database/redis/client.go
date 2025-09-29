package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type Client struct {
	client *redis.Client
	ctx    context.Context
}

type Cache interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{}, expiration time.Duration)
	Delete(key string) error
	DeletePattern(pattern string) error
}

// NewClient creates a new Redis client
func NewClient(addr, password string, db int) (*Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	ctx := context.Background()

	// Test connection
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &Client{
		client: client,
		ctx:    ctx,
	}, nil
}

// Get retrieves a value from cache
func (c *Client) Get(key string) (interface{}, bool) {
	val, err := c.client.Get(c.ctx, key).Result()
	if err == redis.Nil {
		return nil, false
	} else if err != nil {
		return nil, false
	}

	// Try to unmarshal as bool first (for authorization results)
	var result interface{}
	if err := json.Unmarshal([]byte(val), &result); err != nil {
		// If unmarshal fails, return the string value
		return val, true
	}

	return result, true
}

// Set stores a value in cache with expiration
func (c *Client) Set(key string, value interface{}, expiration time.Duration) {
	// Marshal the value to JSON
	jsonData, err := json.Marshal(value)
	if err != nil {
		// If marshaling fails, store as string
		jsonData = []byte(fmt.Sprintf("%v", value))
	}

	err = c.client.Set(c.ctx, key, jsonData, expiration).Err()
	if err != nil {
		// Log error but don't fail the request
		fmt.Printf("failed to set cache key %s: %v\n", key, err)
	}
}

// Delete removes a key from cache
func (c *Client) Delete(key string) error {
	return c.client.Del(c.ctx, key).Err()
}

// DeletePattern removes keys matching a pattern
func (c *Client) DeletePattern(pattern string) error {
	iter := c.client.Scan(c.ctx, 0, pattern, 0).Iterator()
	var keys []string

	for iter.Next(c.ctx) {
		keys = append(keys, iter.Val())
	}

	if err := iter.Err(); err != nil {
		return err
	}

	if len(keys) > 0 {
		return c.client.Del(c.ctx, keys...).Err()
	}

	return nil
}

// Close closes the Redis connection
func (c *Client) Close() error {
	return c.client.Close()
}
