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

// IsEqual returns true if all properties are equal to the other repository.
// Note the function can only compare primitives, it is not able to compare complex types such as slices.
func (repository Repository) IsEqual(other Repository) bool {
	for key, value := range repository.Properties {
		if value != other.Properties[key] {
			return false
		}
	}
	return true
}
