package storage

import "mini-zanzibar/pkg/models"

// Store defines the interface for storing and retrieving relations
type Store interface {
	WriteRelation(relation models.Relation) error
	ReadRelations(object models.Object) ([]models.Relation, error)
	CheckRelation(req models.CheckRequest) (bool, error)
	DeleteRelation(relation models.Relation) error
}
