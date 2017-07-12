package domain

// RepositoryID is the identifier of a nexus repository
type RepositoryID string

// Properties represents each field except id of repository
type Properties map[string]interface{}

// Repository represents a nexus repository
type Repository struct {
	ID         RepositoryID
	Properties Properties
}
