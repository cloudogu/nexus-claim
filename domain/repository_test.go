package domain_test

import (
	"testing"

	"github.com/cloudogu/nexus-claim/domain"
	"github.com/stretchr/testify/assert"
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

	repo2 = createSimpleRepository("other-key", "other")
	assert.False(t, repo1.IsEqual(repo2))
}

func TestRepository_Merge(t *testing.T) {
	propsA := make(domain.Properties)
	propsA["name"] = "a"
	propsA["description"] = "A"
	repoA := domain.Repository{domain.RepositoryID("a"), propsA, domain.RepositoryTypeRepository}

	propsB := make(domain.Properties)
	propsB["name"] = "b"
	propsB["description"] = "B"
	propsB["contact"] = "b@b.de"
	repoB := domain.Repository{domain.RepositoryID("b"), propsB, domain.RepositoryTypeRepository}

	assertion := assert.New(t)

	mergedRepo := repoB.Merge(repoA)
	assertion.Equal(domain.RepositoryID("b"), mergedRepo.ID)

	mergedProps := mergedRepo.Properties
	assertion.Equal("a", mergedProps["name"])
	assertion.Equal("A", mergedProps["description"])
	assertion.Equal("b@b.de", mergedProps["contact"])

	// be sure the original repository has not changed
	assertion.Equal("b", repoB.Properties["name"])
	assertion.Equal("B", repoB.Properties["description"])
}

func createSimpleRepository(key string, value interface{}) domain.Repository {
	properties := make(domain.Properties)
	properties[key] = value

	return domain.Repository{
		ID:         domain.RepositoryID("simple"),
		Properties: properties,
	}
}
