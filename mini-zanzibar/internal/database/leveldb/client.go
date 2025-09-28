package leveldb

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type Client struct {
	db *leveldb.DB
}

type ACLTuple struct {
	Object   string `json:"object"`
	Relation string `json:"relation"`
	User     string `json:"user"`
}

// NewClient creates a new LevelDB client
func NewClient(dbPath string) (*Client, error) {
	db, err := leveldb.OpenFile(dbPath, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to open LevelDB: %w", err)
	}

	return &Client{
		db: db,
	}, nil
}

// Close closes the LevelDB connection
func (c *Client) Close() error {
	return c.db.Close()
}

// StoreTuple stores an ACL tuple in the format: object#relation@user
func (c *Client) StoreTuple(tuple ACLTuple) error {
	key := c.formatTupleKey(tuple)
	value, err := json.Marshal(tuple)
	if err != nil {
		return fmt.Errorf("failed to marshal tuple: %w", err)
	}

	return c.db.Put([]byte(key), value, nil)
}

// GetTuple retrieves a specific ACL tuple
func (c *Client) GetTuple(object, relation, user string) (*ACLTuple, error) {
	tuple := ACLTuple{Object: object, Relation: relation, User: user}
	key := c.formatTupleKey(tuple)

	value, err := c.db.Get([]byte(key), nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get tuple: %w", err)
	}

	var result ACLTuple
	if err := json.Unmarshal(value, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal tuple: %w", err)
	}

	return &result, nil
}

// DeleteTuple removes an ACL tuple
func (c *Client) DeleteTuple(object, relation, user string) error {
	tuple := ACLTuple{Object: object, Relation: relation, User: user}
	key := c.formatTupleKey(tuple)

	return c.db.Delete([]byte(key), nil)
}

// ListTuplesByObject returns all tuples for a specific object
func (c *Client) ListTuplesByObject(object string) ([]ACLTuple, error) {
	prefix := object + "#"
	iter := c.db.NewIterator(util.BytesPrefix([]byte(prefix)), nil)
	defer iter.Release()

	var tuples []ACLTuple
	for iter.Next() {
		var tuple ACLTuple
		if err := json.Unmarshal(iter.Value(), &tuple); err != nil {
			continue // Skip malformed entries
		}
		tuples = append(tuples, tuple)
	}

	return tuples, iter.Error()
}

// ListTuplesByUser returns all tuples for a specific user
func (c *Client) ListTuplesByUser(user string) ([]ACLTuple, error) {
	// TODO: Implement efficient user-based querying
	// For now, we'll scan all keys and filter by user
	iter := c.db.NewIterator(nil, nil)
	defer iter.Release()

	var tuples []ACLTuple
	for iter.Next() {
		var tuple ACLTuple
		if err := json.Unmarshal(iter.Value(), &tuple); err != nil {
			continue // Skip malformed entries
		}
		if tuple.User == user {
			tuples = append(tuples, tuple)
		}
	}

	return tuples, iter.Error()
}

// CheckTuple checks if a specific tuple exists
func (c *Client) CheckTuple(object, relation, user string) (bool, error) {
	tuple, err := c.GetTuple(object, relation, user)
	if err != nil {
		return false, err
	}
	return tuple != nil, nil
}

// formatTupleKey formats the tuple key in the format: object#relation@user
func (c *Client) formatTupleKey(tuple ACLTuple) string {
	return fmt.Sprintf("%s#%s@%s", tuple.Object, tuple.Relation, tuple.User)
}

// parseTupleKey parses a tuple key back into components
func (c *Client) parseTupleKey(key string) (object, relation, user string, err error) {
	// Split by '@' to separate user
	parts := strings.Split(key, "@")
	if len(parts) != 2 {
		return "", "", "", fmt.Errorf("invalid tuple key format: missing @")
	}
	user = parts[1]

	// Split by '#' to separate object and relation
	objRelParts := strings.Split(parts[0], "#")
	if len(objRelParts) != 2 {
		return "", "", "", fmt.Errorf("invalid tuple key format: missing #")
	}
	object = objRelParts[0]
	relation = objRelParts[1]

	return object, relation, user, nil
}
