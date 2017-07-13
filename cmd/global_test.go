package cmd

import (
	"testing"

	"github.com/cloudogu/nexus-claim/domain"
	"github.com/stretchr/testify/assert"
)

func TestCreateOperatorFromActionType(t *testing.T) {
	assert.Equal(t, "+", createOperatorFromActionType(domain.ActionCreate))
	assert.Equal(t, "~", createOperatorFromActionType(domain.ActionModify))
	assert.Equal(t, "-", createOperatorFromActionType(domain.ActionRemove))
	assert.Equal(t, "#", createOperatorFromActionType(domain.ActionType(12)))
}
