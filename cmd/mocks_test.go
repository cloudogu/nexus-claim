package cmd

import "github.com/cloudogu/nexus-claim/domain"

type mockNexusAPIClient struct {
	repositories map[domain.RepositoryID]*domain.Repository
	Created      []domain.Repository
	Modified     []domain.Repository
	Removed      []domain.Repository
}

func (mock *mockNexusAPIClient) add(repo domain.Repository) {
	mock.repositories[repo.ID] = &repo
}

func (mock *mockNexusAPIClient) init() {
	if mock.repositories == nil {
		mock.repositories = make(map[domain.RepositoryID]*domain.Repository)
		mock.add(domain.Repository{
			ID: domain.RepositoryID("apache-snapshots"),
		})
		mock.add(domain.Repository{
			ID: domain.RepositoryID("central-m1"),
		})
		properies := make(domain.Properties)
		properies["name"] = "3rd Party"
		mock.add(domain.Repository{
			ID:         domain.RepositoryID("thirdparty"),
			Properties: properies,
		})
	}
}

func (mock *mockNexusAPIClient) Get(id domain.RepositoryID) (*domain.Repository, error) {
	mock.init()
	return mock.repositories[id], nil
}

func (client *mockNexusAPIClient) Create(repository domain.Repository) error {
	if client.Created == nil {
		client.Created = make([]domain.Repository, 0)
	}
	client.Created = append(client.Created, repository)
	return nil
}

func (client *mockNexusAPIClient) Modify(repository domain.Repository) error {
	if client.Modified == nil {
		client.Modified = make([]domain.Repository, 0)
	}
	client.Modified = append(client.Modified, repository)
	return nil
}

func (client *mockNexusAPIClient) Remove(repository domain.Repository) error {
	if client.Removed == nil {
		client.Removed = make([]domain.Repository, 0)
	}
	client.Removed = append(client.Removed, repository)
	return nil
}
