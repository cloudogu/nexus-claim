package infrastructure

import (
	"testing"

	"github.com/cloudogu/nexus-claim/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFileModelDAO(t *testing.T) {
	modelDAO := NewFileModelDAO(getResourcesDir() + "empty.hcl")
	require.NotNil(t, modelDAO)
}

func TestGetNonExisting(t *testing.T) {
	modelDAO := NewFileModelDAO(getResourcesDir() + "666-from-hell.hcl")
	_, err := modelDAO.Get()
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "could not find model at "+getResourcesDir()+"666-from-hell.hcl")
}

func TestGetEmpty(t *testing.T) {
	modelDAO := NewFileModelDAO(getResourcesDir() + "empty.hcl")
	model, err := modelDAO.Get()
	require.Nil(t, err)
	require.NotNil(t, model)

	require.Equal(t, 0, len(model.Repositories))
}

func TestGetInvalid(t *testing.T) {
	modelDAO := NewFileModelDAO(getResourcesDir() + "invalid.hcl")
	_, err := modelDAO.Get()
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "failed to parse file "+getResourcesDir()+"invalid.hcl")
}

func TestGetWithoutState(t *testing.T) {
	modelDAO := NewFileModelDAO(getResourcesDir() + "without-state.hcl")
	_, err := modelDAO.Get()
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "repository wos has no _state field")
}

func TestGetWithInvalidState(t *testing.T) {
	modelDAO := NewFileModelDAO(getResourcesDir() + "invalid-state.hcl")
	_, err := modelDAO.Get()
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "state invalid of repository irepo is not a valid state")
}

func TestGetWithInvalidStateType(t *testing.T) {
	modelDAO := NewFileModelDAO(getResourcesDir() + "invalid-state-type.hcl")
	_, err := modelDAO.Get()
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "state of repository number is not a string")
}

func TestGetWithEmptyID(t *testing.T) {
	modelDAO := NewFileModelDAO(getResourcesDir() + "empty-id.hcl")
	_, err := modelDAO.Get()
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "repository with empty id found")
}

func TestGetStateGetsRemovedFromProperties(t *testing.T) {
	repositories := get(t, "state.hcl").Repositories
	assert.Nil(t, repositories[0].Properties["_state"])
	assert.Nil(t, repositories[1].Properties["_state"])
}

func TestGetWithNestedProperties(t *testing.T) {
	repository := get(t, "nested-property.hcl").Repositories[0]
	storage, ok := repository.Properties["remoteStorage"].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "https://maven.scm-manager.org/releases/", storage["remoteStorageUrl"])
}

func TestGet(t *testing.T) {
	model := get(t, "simple.hcl")
	require.Equal(t, 1, len(model.Repositories))

	repository := model.Repositories[0]
	assert.Equal(t, domain.RepositoryID("simple"), repository.ID)
	assert.Equal(t, "Simple Repository", repository.Properties["Name"])
	assert.Equal(t, domain.StatePresent, repository.State)
	assert.Equal(t, domain.TypeRepository, repository.Type)
}

func TestGetMultiple(t *testing.T) {
	repositories := get(t, "state.hcl").Repositories
	assert.Equal(t, domain.StatePresent, repositories[0].State)
	assert.Equal(t, domain.StateAbsent, repositories[1].State)
}

func TestGetGroup(t *testing.T) {
	repositories := get(t, "group.hcl").Repositories
	require.Equal(t, 1, len(repositories))

	assert.Equal(t, domain.RepositoryID("public"), repositories[0].ID)
	assert.Equal(t, domain.TypeGroup, repositories[0].Type)
	assert.Equal(t, "hello", repositories[0].Properties["data"])
	assert.Equal(t, domain.StatePresent, repositories[0].State)
}

func get(t *testing.T, file string) domain.Model {
	modelDAO := NewFileModelDAO(getResourcesDir() + file)
	model, err := modelDAO.Get()
	require.Nil(t, err)
	require.NotNil(t, model)
	return model
}

func getResourcesDir() string {
	return "../resources/nexus2/"
}
