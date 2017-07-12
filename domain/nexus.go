package domain

// NexusAPIClient is able to execute crud operations against the nexus api
type NexusAPIClient interface {
	// Get return the repository with the given id
	Get(id RepositoryID) (*Repository, error)
}
