package domain

// Repository represents a nexus repository
type Repository struct {
	ID   string `hcl:",key"`
	Name string
}
