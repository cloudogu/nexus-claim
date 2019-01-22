package domain_test

import (
	"testing"

	"github.com/cloudogu/nexus-claim/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRepository_IsEqualWithStrings(t *testing.T) {
	repo1 := createSimpleRepository("name", "hello")
	repo2 := createSimpleRepository("name", "hello")
	assert.True(t, repo1.IsEqual(repo2))

	repo2 = createSimpleRepository("name", "other")
	assert.False(t, repo1.IsEqual(repo2))
}

func TestRepository_IsEqualWithNumbers(t *testing.T) {
	repo1 := createSimpleRepository("amount", 1)
	repo2 := createSimpleRepository("amount", 1)
	assert.True(t, repo1.IsEqual(repo2))

	repo2 = createSimpleRepository("amount", 2)
	assert.False(t, repo1.IsEqual(repo2))
}

func TestRepository_IsEqualWithBool(t *testing.T) {
	repo1 := createSimpleRepository("exists", true)
	repo2 := createSimpleRepository("exists", true)
	assert.True(t, repo1.IsEqual(repo2))

	repo2 = createSimpleRepository("exists", false)
	assert.False(t, repo1.IsEqual(repo2))
}

func TestRepository_IsEqualWithMissing(t *testing.T) {
	repo1 := createSimpleRepository("name", "hello")
	repo2 := createSimpleRepository("name", "hello")
	assert.True(t, repo1.IsEqual(repo2))

	repo2.Properties["missing"] = "somemissing"
	assert.True(t, repo1.IsEqual(repo2))
}

func TestRepository_Merge(t *testing.T) {
	propsA := make(domain.Properties)
	propsA["name"] = "a"
	propsA["description"] = "A"
	repoA := domain.Repository{ID: domain.RepositoryID("a"), Properties: propsA, Type: domain.TypeRepository}

	propsB := make(domain.Properties)
	propsB["name"] = "b"
	propsB["description"] = "B"
	propsB["contact"] = "b@b.de"
	repoB := domain.Repository{ID: domain.RepositoryID("b"), Properties: propsB, Type: domain.TypeRepository}

	assertion := assert.New(t)

	mergedRepo, err := repoB.Merge(repoA)
	require.Nil(t, err)
	assertion.Equal(domain.RepositoryID("b"), mergedRepo.ID)

	mergedProps := mergedRepo.Properties
	assertion.Equal("a", mergedProps["name"])
	assertion.Equal("A", mergedProps["description"])
	assertion.Equal("b@b.de", mergedProps["contact"])

	// be sure the original repository has not changed
	assertion.Equal("b", repoB.Properties["name"])
	assertion.Equal("B", repoB.Properties["description"])
}

func TestRepository_GetRecipeName(t *testing.T) {
	props := make(domain.Properties)
	props["recipeName"] = "repo"
	repo := domain.Repository{Properties: props}

	actual, err := repo.GetRecipeName()

	assert.Nil(t, err)
	assert.Equal(t, "repo", actual)
}

func TestRepository_GetRecipeName_error(t *testing.T) {
	props := make(domain.Properties)
	repo := domain.Repository{Properties: props}

	actual, err := repo.GetRecipeName()

	assert.Empty(t, actual)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "could not find property 'recipeName' in repository")
}

func TestRepository_Clone(t *testing.T) {
	props := make(domain.Properties)
	props["name"] = "name1"
	props["description"] = "desc1"
	props["contact"] = "asdf@asdf.de"
	submap := make(map[string]interface{})
	submap["subentry"] = "hey"
	props["maven"] = submap
	repo := domain.Repository{ID: domain.RepositoryID("repo"), Properties: props, Type: domain.TypeRepository}

	actual := repo.Clone()

	props2 := make(domain.Properties)
	props2["name"] = "name1"
	props2["description"] = "desc1"
	props2["contact"] = "asdf@asdf.de"
	submap2 := make(map[string]interface{})
	submap2["subentry"] = "hey"
	props2["maven"] = submap2
	expected := domain.Repository{ID: domain.RepositoryID("repo"), Properties: props2, Type: domain.TypeRepository}
	assert.Equal(t, expected, actual)
	assert.Equal(t, expected.Properties, actual.Properties)
}

func createSimpleRepository(key string, value interface{}) domain.Repository {
	properties := make(domain.Properties)
	properties[key] = value

	return domain.Repository{
		ID:         domain.RepositoryID("simple"),
		Properties: properties,
	}
}
