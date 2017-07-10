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

func TestUnmarshalDetailed(t *testing.T) {
	repository := unmarshalRepository(t, "../resources/005-detail.hcl")
	assert.Equal(t, "releases", repository.ID)
	assert.Equal(t, "Releases", repository.Name)
	assert.Equal(t, "maven2", repository.Format)
	assert.Equal(t, "hosted", repository.RepoType)
	assert.Equal(t, "RELEASE", repository.RepoPolicy)
	assert.Equal(t, "maven2", repository.Provider)
	assert.True(t, repository.Browseable)
	assert.True(t, repository.Indexable)
	assert.True(t, repository.Exposed)
	assert.False(t, repository.DownloadRemoteIndexes)
	assert.Equal(t, "ALLOW_WRITE_ONCE", repository.WritePolicy)
	assert.Equal(t, "org.sonatype.nexus.proxy.repository.Repository", repository.ProviderRole)
	assert.Equal(t, 1440, repository.NotFoundCacheTTL)

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
