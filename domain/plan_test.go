package domain_test

import (
	"testing"

	"github.com/cloudogu/nexus-claim/domain"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreatePlanWithEmptyModel(t *testing.T) {
	dao := &mockModelDAO{domain.Model{}, nil}
	client := &mockNexusAPIClient{}

	plan, err := domain.CreatePlan(dao, client)
	require.Nil(t, err)
	require.NotNil(t, plan)
	assert.Equal(t, 0, len(plan.GetActions()))
}

func TestCreatePlanFailedToReadModel(t *testing.T) {
	dao := &mockModelDAO{domain.Model{}, errors.New("-- no --")}
	client := &mockNexusAPIClient{}

	_, err := domain.CreatePlan(dao, client)
	require.NotNil(t, err)
	assert.Contains(t, err.Error(), "-- no --")
}

func TestCreatePlanFailedToReadFromNexusAPI(t *testing.T) {
	model := createTestModel()

	dao := &mockModelDAO{model: model}
	client := &mockNexusAPIClient{err: errors.New("api down")}

	_, err := domain.CreatePlan(dao, client)
	require.NotNil(t, err)
	assert.Contains(t, err.Error(), "api down")
}

func TestCreatePlanCreateMissingRepository(t *testing.T) {
	model := createTestModel()

	dao := &mockModelDAO{model: model}
	client := &mockNexusAPIClient{}

	plan, err := domain.CreatePlan(dao, client)
	require.Nil(t, err)
	require.NotNil(t, plan)

	require.Equal(t, 1, len(plan.GetActions()))
	action := plan.GetActions()[0]

	assert.Equal(t, domain.ActionCreate, action.GetType())
	assert.Equal(t, domain.RepositoryID("missing-repo"), action.GetRepository().ID)
}

func TestCreatePlanWithUnknownState(t *testing.T) {
	model := domain.Model{
		Repositories: []domain.ModelRepository{
			{
				Repository: domain.Repository{
					ID:         domain.RepositoryID("absent-repo"),
					Properties: make(domain.Properties),
				},
				State: domain.State("unknown"),
			},
		},
	}

	dao := &mockModelDAO{model, nil}
	client := &mockNexusAPIClient{}

	plan, err := domain.CreatePlan(dao, client)
	require.Nil(t, err)
	require.NotNil(t, plan)
	require.Equal(t, 0, len(plan.GetActions()))
}

func TestCreatePlanWithAbsentWhichIsAlreadyAbsent(t *testing.T) {
	model := domain.Model{
		Repositories: []domain.ModelRepository{
			{
				Repository: domain.Repository{
					ID:         domain.RepositoryID("absent-repo"),
					Properties: make(domain.Properties),
				},
				State: domain.StateAbsent,
			},
		},
	}

	dao := &mockModelDAO{model, nil}
	client := &mockNexusAPIClient{}

	plan, err := domain.CreatePlan(dao, client)
	require.Nil(t, err)
	require.NotNil(t, plan)
	require.Equal(t, 0, len(plan.GetActions()))
}

func TestCreatePlanRemoveAbsentRepository(t *testing.T) {
	repository := domain.Repository{
		ID:         domain.RepositoryID("absent-repo"),
		Properties: make(domain.Properties),
	}

	model := domain.Model{
		Repositories: []domain.ModelRepository{
			{
				Repository: repository,
				State:      domain.StateAbsent,
			},
		},
	}

	dao := &mockModelDAO{model, nil}
	client := &mockNexusAPIClient{}
	client.Add(repository)

	plan, err := domain.CreatePlan(dao, client)
	require.Nil(t, err)
	require.NotNil(t, plan)

	require.Equal(t, 1, len(plan.GetActions()))

	action := plan.GetActions()[0]
	assert.Equal(t, domain.ActionRemove, action.GetType())
	assert.Equal(t, repository, action.GetRepository())
}

func TestCreatePlanWithNonChanged(t *testing.T) {
	repository := domain.Repository{
		ID:         domain.RepositoryID("repo"),
		Properties: make(domain.Properties),
	}

	model := domain.Model{
		Repositories: []domain.ModelRepository{
			{
				Repository: repository,
				State:      domain.StatePresent,
			},
		},
	}

	dao := &mockModelDAO{model, nil}
	client := &mockNexusAPIClient{}
	client.Add(repository)

	plan, err := domain.CreatePlan(dao, client)
	require.Nil(t, err)
	require.NotNil(t, plan)

	assert.Equal(t, 0, len(plan.GetActions()))
}

func TestCreatePlanWithChangedProperty(t *testing.T) {
	id := domain.RepositoryID("repo")

	modelProperties := make(domain.Properties)
	modelProperties["name"] = "super simple"
	modelProperties["type"] = "maven2"
	modelProperties["mp"] = "model"

	modelRepository := domain.Repository{
		ID:         id,
		Properties: modelProperties,
	}

	model := domain.Model{
		Repositories: []domain.ModelRepository{
			{
				Repository: modelRepository,
				State:      domain.StatePresent,
			},
		},
	}
	dao := &mockModelDAO{model, nil}

	// ---

	clientProperties := make(domain.Properties)
	clientProperties["name"] = "simple"
	clientProperties["type"] = "maven2"
	clientProperties["cp"] = "client"

	clientRepository := domain.Repository{
		ID:         id,
		Properties: clientProperties,
	}

	client := &mockNexusAPIClient{}
	client.Add(clientRepository)

	plan, err := domain.CreatePlan(dao, client)
	require.Nil(t, err)
	require.NotNil(t, plan)

	require.Equal(t, 1, len(plan.GetActions()))

	action := plan.GetActions()[0]
	assert.Equal(t, domain.ActionModify, action.GetType())
	assert.Equal(t, id, action.GetRepository().ID)
	assert.Equal(t, "super simple", action.GetRepository().Properties["name"])
	assert.Equal(t, "maven2", action.GetRepository().Properties["type"])
	assert.Equal(t, "model", action.GetRepository().Properties["mp"])
	assert.Equal(t, "client", action.GetRepository().Properties["cp"])
	assert.Equal(t, 4, len(action.GetRepository().Properties))
}

func createTestModel() domain.Model {
	return domain.Model{
		Repositories: []domain.ModelRepository{
			{
				Repository: domain.Repository{
					ID:         domain.RepositoryID("missing-repo"),
					Properties: make(domain.Properties),
				},
				State: domain.StatePresent,
			},
		},
	}
}

type mockNexusAPIClient struct {
	repositories map[domain.RepositoryID]*domain.Repository
	err          error
}

func (mock *mockNexusAPIClient) Add(repository domain.Repository) {
	if mock.repositories == nil {
		mock.repositories = make(map[domain.RepositoryID]*domain.Repository)
	}
	mock.repositories[repository.ID] = &repository
}

func (mock *mockNexusAPIClient) Get(id domain.RepositoryID) (*domain.Repository, error) {
	return mock.repositories[id], mock.err
}

type mockModelDAO struct {
	model domain.Model
	err   error
}

func (mock *mockModelDAO) Get() (domain.Model, error) {
	return mock.model, mock.err
}
