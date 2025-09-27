package storage

import (
	"fmt"
	"mini-zanzibar/pkg/models"
	"sync"
)

// MemoryStore implements in-memory storage for relations
type MemoryStore struct {
	relations map[string][]models.Relation
	mutex     sync.RWMutex
}

// NewMemoryStore creates a new in-memory store
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		relations: make(map[string][]models.Relation),
	}
}

// WriteRelation stores a relation in memory
func (ms *MemoryStore) WriteRelation(relation models.Relation) error {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()

	key := fmt.Sprintf("%s:%s", relation.Object.Type, relation.Object.ID)
	fmt.Printf("[STORAGE] Writing relation with key: %s, relation: %+v\n", key, relation)

	// Check if relation already exists to avoid duplicates
	for _, existingRel := range ms.relations[key] {
		if existingRel.Relation == relation.Relation &&
			existingRel.Subject.Type == relation.Subject.Type &&
			existingRel.Subject.ID == relation.Subject.ID {
			fmt.Printf("[STORAGE] Relation already exists, skipping\n")
			return nil // Relation already exists
		}
	}

	ms.relations[key] = append(ms.relations[key], relation)
	fmt.Printf("[STORAGE] Total relations in store: %d\n", len(ms.relations))
	fmt.Printf("[STORAGE] Relations for key %s: %d\n", key, len(ms.relations[key]))
	return nil
}

// ReadRelations retrieves all relations for a given object
func (ms *MemoryStore) ReadRelations(object models.Object) ([]models.Relation, error) {
	ms.mutex.RLock()
	defer ms.mutex.RUnlock()

	key := fmt.Sprintf("%s:%s", object.Type, object.ID)
	relations := ms.relations[key]

	fmt.Printf("[STORAGE] Reading relations for key: %s\n", key)
	fmt.Printf("[STORAGE] Found %d relations\n", len(relations))
	fmt.Printf("[STORAGE] All keys in store: %v\n", ms.getKeys())

	// Return a copy to avoid race conditions
	result := make([]models.Relation, len(relations))
	copy(result, relations)

	return result, nil
}

// Helper function to get all keys (for debugging)
func (ms *MemoryStore) getKeys() []string {
	keys := make([]string, 0, len(ms.relations))
	for k := range ms.relations {
		keys = append(keys, k)
	}
	return keys
}

// CheckRelation checks if a specific relation exists
func (ms *MemoryStore) CheckRelation(req models.CheckRequest) (bool, error) {
	relations, err := ms.ReadRelations(req.Object)
	if err != nil {
		return false, err
	}

	for _, rel := range relations {
		if rel.Relation == req.Relation &&
			rel.Subject.Type == req.Subject.Type &&
			rel.Subject.ID == req.Subject.ID {
			return true, nil
		}
	}

	return false, nil
}

// DeleteRelation removes a specific relation
func (ms *MemoryStore) DeleteRelation(relation models.Relation) error {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()

	key := fmt.Sprintf("%s:%s", relation.Object.Type, relation.Object.ID)
	relations := ms.relations[key]

	for i, rel := range relations {
		if rel.Relation == relation.Relation &&
			rel.Subject.Type == relation.Subject.Type &&
			rel.Subject.ID == relation.Subject.ID {
			ms.relations[key] = append(relations[:i], relations[i+1:]...)
			break
		}
	}

	return nil
}
