package infrastructure_test

import (
	"testing"

	"github.com/cloudogu/nexus-claim/domain"
	"github.com/cloudogu/nexus-claim/infrastructure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSerializePlanAndDeserializePlan(t *testing.T) {
	t123 := domain.RepositoryID("test-123")
	t456 := domain.RepositoryID("test-456")
	t789 := domain.RepositoryID("test-789")

	plan := &domain.Plan{}
	plan.Create(domain.Repository{ID: t123})
	plan.Modify(domain.Repository{ID: t456})
	plan.Remove(domain.Repository{ID: t789})

	serializedPlan, err := infrastructure.SerializePlan(plan)
	require.Nil(t, err)

	deserializedPlan, err := infrastructure.DeserializePlan(serializedPlan)
	require.Nil(t, err)

	actions := deserializedPlan.GetActions()
	assert.Equal(t, 3, len(actions))
	assert.True(t, containsAction(actions, domain.ActionCreate, t123))
	assert.True(t, containsAction(actions, domain.ActionModify, t456))
	assert.True(t, containsAction(actions, domain.ActionRemove, t789))
}

func containsAction(actions []domain.Action, actionType domain.ActionType, repositoryID domain.RepositoryID) bool {
	for _, action := range actions {
		if action.GetType() == actionType && action.GetRepository().ID == repositoryID {
			return true
		}
	}
	return false
}
