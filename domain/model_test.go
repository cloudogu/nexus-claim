package domain

import (
	"io/ioutil"
	"testing"

	"github.com/hashicorp/hcl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnmarshalModel(t *testing.T) {
	model := unmarshalModel(t, "../resources/002-simple.hcl")
	assert.Equal(t, 1, len(model.Repositories))
}

func TestUnmarshalOfEmbeddedRepository(t *testing.T) {
	repository := unmarshalModel(t, "../resources/002-simple.hcl").Repositories[0]
	assert.Equal(t, "simple", repository.ID)
	assert.Equal(t, "Simple Repository", repository.Name)
}

func TestUnmarshalOfState(t *testing.T) {
	repositories := unmarshalModel(t, "../resources/003-state.hcl").Repositories
	assert.Equal(t, StatePresent, repositories[0].State)
	assert.Equal(t, StateAbsent, repositories[1].State)
}

func unmarshalModel(t *testing.T, path string) Model {
	bytes, err := ioutil.ReadFile(path)
	require.Nil(t, err)

	model := Model{}
	err = hcl.Unmarshal(bytes, &model)
	require.Nil(t, err)

	return model
}
