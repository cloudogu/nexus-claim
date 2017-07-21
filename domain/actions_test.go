package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateAction_ToString(t *testing.T) {
	action := &createAction{baseAction{Type: ActionCreate, Repository: Repository{ID: RepositoryID("test-123")}}}
	assert.Equal(t, "+ test-123", action.ToString())
}

func TestModifyAction_ToString(t *testing.T) {
	action := &modifyAction{baseAction{Type: ActionModify, Repository: Repository{ID: RepositoryID("test-123")}}}
	assert.Equal(t, "~ test-123", action.ToString())
}

func TestRemoveAction_ToString(t *testing.T) {
	action := &removeAction{baseAction{Type: ActionRemove, Repository: Repository{ID: RepositoryID("test-123")}}}
	assert.Equal(t, "- test-123", action.ToString())
}
