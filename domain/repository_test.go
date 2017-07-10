package domain

import (
	"io/ioutil"
	"testing"

	"github.com/hashicorp/hcl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnmarshal(t *testing.T) {
	repository := unmarshalRepository(t, "../resources/002-simple.hcl")
	assert.Equal(t, "simple", repository.ID)
	assert.Equal(t, "Simple Repository", repository.Name)
}

func unmarshalRepository(t *testing.T, path string) *Repository {
	bytes, err := ioutil.ReadFile(path)
	require.Nil(t, err)

	model := struct {
		Repository []*Repository `hcl:"repository"`
	}{}

	err = hcl.Unmarshal(bytes, &model)
	require.Nil(t, err)

	return model.Repository[0]
}
