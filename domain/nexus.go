package domain

// NexusAPIReader is able to read repositories from the nexus api
type NexusAPIReader interface {
	// Get return the repository with the given id
	Get(id RepositoryID) (*Repository, error)
}

// NexusAPIWriter is able to create, modify and remove nexus repositories
type NexusAPIWriter interface {
	// Create creates a new nexus repository
	Create(repository Repository) error

	// Modify modifies an existing nexus repository
	Modify(repository Repository) error

	// Remove removes an existing repository from nexus server
	Remove(repository Repository) error
}

// NexusAPIClient is able to execute crud operations against the nexus api
type NexusAPIClient interface {
	NexusAPIReader
	NexusAPIWriter
}
