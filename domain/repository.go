package domain

// Repository represents a nexus repository
type Repository struct {
	ID                    string `hcl:",key"`
	Name                  string
	Format                string
	RepoType              string
	RepoPolicy            string
	Provider              string
	Browseable            bool
	Indexable             bool
	Exposed               bool
	DownloadRemoteIndexes bool
	WritePolicy           string
	ProviderRole          string
	NotFoundCacheTTL      int
}
