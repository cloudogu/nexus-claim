package infrastructure

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewFileModelDAO(t *testing.T) {
	modelDAO := NewFileModelDAO("../resources/001-empty.hcl")
	require.NotNil(t, modelDAO)
}

func TestGetNonExisting(t *testing.T) {
	modelDAO := NewFileModelDAO("../resources/666-from-hell.hcl")
	_, err := modelDAO.Get()
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "could not find model at ../resources/666-from-hell.hcl")
}

func TestGetEmpty(t *testing.T) {
	modelDAO := NewFileModelDAO("../resources/001-empty.hcl")
	model, err := modelDAO.Get()
	require.Nil(t, err)
	require.NotNil(t, model)

	require.Equal(t, 0, len(model.Repositories))
}

func TestGetInvalid(t *testing.T) {
	modelDAO := NewFileModelDAO("../resources/004-invalid.hcl")
	_, err := modelDAO.Get()
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "failed to unmarshal model ../resources/004-invalid.hcl")
}

func TestGet(t *testing.T) {
	modelDAO := NewFileModelDAO("../resources/002-simple.hcl")
	model, err := modelDAO.Get()
	require.Nil(t, err)
	require.NotNil(t, model)

	require.Equal(t, 1, len(model.Repositories))
}
