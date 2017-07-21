package domain

import (
	"fmt"

	"github.com/pkg/errors"
)

// ActionType describes the type of the action create, modify or remove
type ActionType uint8

const (
	// ActionCreate creates the missing repository on nexus
	ActionCreate ActionType = iota
	// ActionModify modifies an existing nexus repository, because some properties have changed
	ActionModify
	// ActionRemove removes an existing nexus repository, because the model describe it with _state = absent
	ActionRemove
)

// Action describes a single unit of work to sync the model with nexus
type Action interface {
	// GetType returns the type of the action
	GetType() ActionType
	// GetRepository returns the repository which is the target of the action
	GetRepository() Repository
	// Execute executes the action against the nexus api
	Execute(writer NexusAPIWriter) error
	// ToString returns a human readable representation of the action
	ToString() string
}

type baseAction struct {
	Type       ActionType
	Repository Repository
}

func (action *baseAction) GetType() ActionType {
	return action.Type
}

func (action *baseAction) GetRepository() Repository {
	return action.Repository
}

func (action *baseAction) subject(operator string) string {
	return fmt.Sprintf("%s %s", operator, action.Repository.ID)
}

// -- Create

type createAction struct {
	baseAction
}

func (action *createAction) ToString() string {
	return action.subject("+")
}

func (action *createAction) Execute(writer NexusAPIWriter) error {
	err := writer.Create(action.Repository)
	if err != nil {
		return errors.Wrapf(err, "failed to create repository %s", action.Repository.ID)
	}
	return nil
}

// -- Modify

type modifyAction struct {
	baseAction
}

func (action *modifyAction) ToString() string {
	return action.subject("~")
}

func (action *modifyAction) Execute(writer NexusAPIWriter) error {
	err := writer.Modify(action.Repository)
	if err != nil {
		return errors.Wrapf(err, "failed to modify repository %s", action.Repository.ID)
	}
	return nil
}

// -- Remove

type removeAction struct {
	baseAction
}

func (action *removeAction) ToString() string {
	return action.subject("-")
}

func (action *removeAction) Execute(writer NexusAPIWriter) error {
	err := writer.Remove(action.Repository.ID)
	if err != nil {
		return errors.Wrapf(err, "failed to remove repository %s", action.Repository.ID)
	}
	return nil
}
