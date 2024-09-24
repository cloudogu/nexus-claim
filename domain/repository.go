package domain

import (
	"fmt"
	"github.com/imdario/mergo"
	"github.com/kr/pretty"
	"strings"
)

const (
	// TypeRepository can be a hosted, proxy or a virtual nexus repository
	TypeRepository RepositoryType = iota
	// TypeGroup are able to group a set of repositories to single one
	TypeGroup
	// repositoryRecipeNameKey denotes the property name which points to a Nexus recipe name (f.i. maven2-hosted)
	repositoryRecipeNameKey = "recipeName"
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
	diff := pretty.Diff(repository.Properties, other.Properties)

	changes := make([]string, 0)
	for _, diffEntry := range diff {
		if !strings.Contains(diffEntry, "missing") {
			changes = append(changes, diffEntry)
		}
	}

	return len(changes) == 0
}

func (repository Repository) String() string {
	return fmt.Sprintf("ID: %s\nType: %d\nProperties: %s",
		string(repository.ID), repository.Type, repository.Properties)
}

// Merge copies all properties from the other repository, merges them with this repository and returns a new repository.
func (repository Repository) Merge(other Repository) (Repository, error) {
	properties := repository.cloneProperties()

	err := mergo.Merge(&properties, other.Properties, mergo.WithOverride)
	if err != nil {
		return repository, fmt.Errorf("failed to merge repository properties: %w", err)
	}

	return Repository{ID: repository.ID, Properties: properties, Type: repository.Type}, nil
}

func (repository Repository) cloneProperties() Properties {
	properties := make(Properties)
	for key, value := range repository.Properties {
		properties[key] = value
	}
	return properties
}

// Clone returns a copy of the original repository.
func (repository Repository) Clone() Repository {
	properties := repository.cloneProperties()

	return Repository{ID: repository.ID, Properties: properties, Type: repository.Type}
}

// GetRecipeName returns the value of a Nexus recipe name (f.i. maven2-hosted). If no recipe name is set an error will
// be returned.
func (repository Repository) GetRecipeName() (string, error) {
	recipeName := repository.Properties[repositoryRecipeNameKey]
	if recipeName == nil {
		return "", fmt.Errorf("could not find property 'recipeName' in repository " + string(repository.ID))
	}

	return recipeName.(string), nil
}
