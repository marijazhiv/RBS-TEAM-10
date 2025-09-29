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
	// key := c.formatTupleKey(tuple)
	// value, err := json.Marshal(tuple)
	// if err != nil {
	// 	return fmt.Errorf("failed to marshal tuple: %w", err)
	// }

	// return c.db.Put([]byte(key), value, nil)

	// Primary key: object#relation@user
	primaryKey := c.formatTupleKey(tuple)

	// Reverse index key: user@object#relation
	reverseKey := c.formatReverseKey(tuple)

	value, err := json.Marshal(tuple)
	if err != nil {
		return fmt.Errorf("failed to marshal tuple: %w", err)
	}

	// Use batch for atomic operation
	batch := new(leveldb.Batch)
	batch.Put([]byte(primaryKey), value)
	batch.Put([]byte(reverseKey), []byte{}) // Reverse index doesn't need value, just the key

	return c.db.Write(batch, nil)
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
	// tuple := ACLTuple{Object: object, Relation: relation, User: user}
	// key := c.formatTupleKey(tuple)

	// return c.db.Delete([]byte(key), nil)
	tuple := ACLTuple{Object: object, Relation: relation, User: user}
	primaryKey := c.formatTupleKey(tuple)
	reverseKey := c.formatReverseKey(tuple)

	// Use batch for atomic operation
	batch := new(leveldb.Batch)
	batch.Delete([]byte(primaryKey))
	batch.Delete([]byte(reverseKey))

	return c.db.Write(batch, nil)
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

// ListTuplesByUser returns all tuples for a specific user (USING REVERSE INDEX)
func (c *Client) ListTuplesByUser(user string) ([]ACLTuple, error) {
	// TODO: Implement efficient user-based querying
	// For now, we'll scan all keys and filter by user
	// iter := c.db.NewIterator(nil, nil)
	// defer iter.Release()

	// var tuples []ACLTuple
	// for iter.Next() {
	// 	var tuple ACLTuple
	// 	if err := json.Unmarshal(iter.Value(), &tuple); err != nil {
	// 		continue // Skip malformed entries
	// 	}
	// 	if tuple.User == user {
	// 		tuples = append(tuples, tuple)
	// 	}
	// }

	// return tuples, iter.Error()

	// Use reverse index for efficient querying
	prefix := user + "@"
	iter := c.db.NewIterator(util.BytesPrefix([]byte(prefix)), nil)
	defer iter.Release()

	var tuples []ACLTuple
	for iter.Next() {
		// Parse the reverse key to get object and relation
		object, relation, parsedUser, err := c.parseReverseKey(string(iter.Key()))
		if err != nil {
			continue // Skip malformed entries
		}

		// Get the actual tuple using primary key
		tuple, err := c.GetTuple(object, relation, parsedUser)
		if err != nil || tuple == nil {
			continue // Skip if tuple not found
		}

		tuples = append(tuples, *tuple)
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

// formatReverseKey formats the reverse index key in the format: user@object#relation
func (c *Client) formatReverseKey(tuple ACLTuple) string {
	return fmt.Sprintf("%s@%s#%s", tuple.User, tuple.Object, tuple.Relation)
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

// parseReverseKey parses a reverse index key back into components
func (c *Client) parseReverseKey(key string) (object, relation, user string, err error) {
	// Split by '@' to separate user
	parts := strings.Split(key, "@")
	if len(parts) != 2 {
		return "", "", "", fmt.Errorf("invalid reverse key format: missing @")
	}
	user = parts[0]

	// Split by '#' to separate object and relation
	objRelParts := strings.Split(parts[1], "#")
	if len(objRelParts) != 2 {
		return "", "", "", fmt.Errorf("invalid reverse key format: missing #")
	}

	object = objRelParts[0]
	relation = objRelParts[1]

	return object, relation, user, nil

}

// ListTuplesByObjectAndRelation returns all tuples for a specific object and relation
func (c *Client) ListTuplesByObjectAndRelation(object, relation string) ([]ACLTuple, error) {
	prefix := fmt.Sprintf("%s#%s@", object, relation)
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

// ListTuplesByObjectPagination returns paginated tuples for a specific object
func (c *Client) ListTuplesByObjectPagination(object string, page, pageSize int) ([]ACLTuple, int, error) {
	prefix := object + "#"
	iter := c.db.NewIterator(util.BytesPrefix([]byte(prefix)), nil)
	defer iter.Release()

	var tuples []ACLTuple
	total := 0
	skip := (page - 1) * pageSize
	count := 0

	for iter.Next() {
		total++
		if count < skip {
			count++
			continue
		}

		if len(tuples) >= pageSize {
			continue
		}

		var tuple ACLTuple
		if err := json.Unmarshal(iter.Value(), &tuple); err != nil {
			continue // Skip malformed entries
		}
		tuples = append(tuples, tuple)
		count++
	}

	return tuples, total, iter.Error()
}

// ListTuplesByUserPagination returns paginated tuples for a specific user (USING REVERSE INDEX)
func (c *Client) ListTuplesByUserPagination(user string, page, pageSize int) ([]ACLTuple, int, error) {
	// Use reverse index for efficient pagination
	prefix := user + "@"
	iter := c.db.NewIterator(util.BytesPrefix([]byte(prefix)), nil)
	defer iter.Release()

	var tuples []ACLTuple
	total := 0
	skip := (page - 1) * pageSize
	count := 0

	for iter.Next() {
		total++
		if count < skip {
			count++
			continue
		}

		if len(tuples) >= pageSize {
			continue
		}

		// Parse the reverse key to get object and relation
		object, relation, parsedUser, err := c.parseReverseKey(string(iter.Key()))
		if err != nil {
			continue // Skip malformed entries
		}

		// Get the actual tuple using primary key
		tuple, err := c.GetTuple(object, relation, parsedUser)
		if err != nil || tuple == nil {
			continue // Skip if tuple not found
		}

		tuples = append(tuples, *tuple)
		count++
	}

	return tuples, total, iter.Error()
}

// ListTuplesByUserAndRelation returns all tuples for a specific user and relation (USING REVERSE INDEX)
func (c *Client) ListTuplesByUserAndRelation(user, relation string) ([]ACLTuple, error) {
	// This is much more efficient with reverse index
	prefix := user + "@"
	iter := c.db.NewIterator(util.BytesPrefix([]byte(prefix)), nil)
	defer iter.Release()

	var tuples []ACLTuple
	for iter.Next() {
		// Parse the reverse key to get object and relation
		object, rel, parsedUser, err := c.parseReverseKey(string(iter.Key()))
		if err != nil || rel != relation {
			continue // Skip if relation doesn't match or malformed
		}

		// Get the actual tuple using primary key
		tuple, err := c.GetTuple(object, relation, parsedUser)
		if err != nil || tuple == nil {
			continue // Skip if tuple not found
		}

		tuples = append(tuples, *tuple)
	}

	return tuples, iter.Error()
}

// MigrateExistingData adds reverse indexes for existing data (one-time operation)
func (c *Client) MigrateExistingData() error {
	iter := c.db.NewIterator(util.BytesPrefix([]byte("")), nil)
	defer iter.Release()

	batch := new(leveldb.Batch)
	batchSize := 0
	const maxBatchSize = 1000

	for iter.Next() {
		key := string(iter.Key())

		// Skip if it's already a reverse index key
		if strings.Contains(key, "@") && strings.Index(key, "@") < strings.Index(key, "#") {
			continue
		}

		// Parse the primary key
		object, relation, user, err := c.parseTupleKey(key)
		if err != nil {
			continue // Skip malformed keys
		}

		// Create reverse index
		reverseKey := c.formatReverseKey(ACLTuple{Object: object, Relation: relation, User: user})
		batch.Put([]byte(reverseKey), []byte{})

		batchSize++
		if batchSize >= maxBatchSize {
			if err := c.db.Write(batch, nil); err != nil {
				return fmt.Errorf("failed to write batch: %w", err)
			}
			batch = new(leveldb.Batch)
			batchSize = 0
		}
	}

	// Write remaining batch
	if batchSize > 0 {
		if err := c.db.Write(batch, nil); err != nil {
			return fmt.Errorf("failed to write final batch: %w", err)
		}
	}

	return iter.Error()
}
