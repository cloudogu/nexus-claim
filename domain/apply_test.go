package domain_test

import (
	"fmt"
	"testing"

	"github.com/cloudogu/nexus-claim/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApplyPlanWithEmptyPlan(t *testing.T) {
	plan := &domain.Plan{}
	writer := newMockNexusAPIWriter()

	err := domain.ApplyPlan(writer, plan)
	require.Nil(t, err)
}

func TestApplyPlanWithErrorOnCreate(t *testing.T) {
	plan := &domain.Plan{}
	plan.Create(domain.Repository{ID: domain.RepositoryID("c1")})
	writer := newMockNexusAPIWriter()
	writer.err = fmt.Errorf("server down")
	err := domain.ApplyPlan(writer, plan)
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "server down")
	require.Contains(t, err.Error(), "c1")
}

func TestApplyPlanWithErrorOnRemove(t *testing.T) {
	plan := &domain.Plan{}
	plan.Remove(domain.Repository{ID: domain.RepositoryID("r1")})
	writer := newMockNexusAPIWriter()
	writer.err = fmt.Errorf("server down")
	err := domain.ApplyPlan(writer, plan)
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "server down")
	require.Contains(t, err.Error(), "r1")
}

func TestApplyPlan(t *testing.T) {
	plan := &domain.Plan{}
	plan.Create(domain.Repository{ID: domain.RepositoryID("c1")})
	plan.Create(domain.Repository{ID: domain.RepositoryID("c2")})
	plan.Modify(domain.Repository{ID: domain.RepositoryID("m1")})
	plan.Remove(domain.Repository{ID: domain.RepositoryID("r1")})
	writer := newMockNexusAPIWriter()

	err := domain.ApplyPlan(writer, plan)
	require.Nil(t, err)
	require.Equal(t, 2, len(writer.creations))
	assert.Equal(t, domain.RepositoryID("c1"), writer.creations[0].ID)
	require.Equal(t, 1, len(writer.modifications))
	require.Equal(t, 1, len(writer.removes))
	assert.Equal(t, domain.RepositoryID("r1"), writer.removes[0].ID)
}

func newMockNexusAPIWriter() *mockNexusAPIWriter {
	return &mockNexusAPIWriter{
		creations:     make([]domain.Repository, 0),
		modifications: make([]domain.Repository, 0),
		removes:       make([]domain.Repository, 0),
	}
}

type mockNexusAPIWriter struct {
	err           error
	creations     []domain.Repository
	modifications []domain.Repository
	removes       []domain.Repository
}

func (mock *mockNexusAPIWriter) Create(repository domain.Repository) error {
	if mock.err != nil {
		return mock.err
	}
	mock.creations = append(mock.creations, repository)
	return nil
}

func (mock *mockNexusAPIWriter) Modify(repository domain.Repository) error {
	if mock.err != nil {
		return mock.err
	}
	mock.modifications = append(mock.modifications, repository)
	return nil
}

func (mock *mockNexusAPIWriter) Remove(repository domain.Repository) error {
	if mock.err != nil {
		return mock.err
	}
	mock.removes = append(mock.removes, repository)
	return nil
}
