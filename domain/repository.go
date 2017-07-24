package domain

import "log"

const (
	// TypeRepository can be a hosted, proxy or a virtual nexus repository
	TypeRepository RepositoryType = iota
	// TypeGroup are able to group a set of repositories to single one
	TypeGroup
)

// RepositoryID is the identifier of a nexus repository
type RepositoryID string

// RepositoryType defines the type of a repository
type RepositoryType uint8

// Properties represents each field except id of repository
type Properties map[string]interface{}

// Repository represents a nexus repository
type Repository struct {
	ID         RepositoryID
	Properties Properties
	Type       RepositoryType
}

// IsEqual returns true if all properties are equal to the other repository.
func (repository Repository) IsEqual(other Repository) bool {
	for key, value := range repository.Properties {
		if !IsEqual(value, other.Properties[key]) {
			log.Printf("property %s of repository %s has changed (%v != %v)", key, repository.ID, value, other.Properties[key])
			return false
		}
	}
	return true
}

// Merge copies all properties from the other repository, merges them with this repository and returns a new repository.
func (repository Repository) Merge(other Repository) Repository {
	properties := make(Properties)
	for key, value := range repository.Properties {
		properties[key] = value
	}
	for key, value := range other.Properties {
		properties[key] = value
	}
	return Repository{ID: repository.ID, Properties: properties, Type: repository.Type}
}
