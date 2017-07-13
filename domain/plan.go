package domain

import (
	"encoding/json"

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
type Action struct {
	Type       ActionType
	Repository Repository
}

// Plan is a set of actions which must be done to sync the model with nexus
type Plan struct {
	actions []Action
}

// GetActions returns the planned actions
func (plan *Plan) GetActions() []Action {
	return plan.actions
}

// Create adds a create action to the plan
func (plan *Plan) Create(repository Repository) {
	plan.action(ActionCreate, repository)
}

// Modify adds a modify action to the plan
func (plan *Plan) Modify(repository Repository) {
	plan.action(ActionModify, repository)
}

// Remove adds a remove action to the plan
func (plan *Plan) Remove(repository Repository) {
	plan.action(ActionRemove, repository)
}

func (plan *Plan) action(actionType ActionType, repository Repository) {
	action := Action{Type: actionType, Repository: repository}
	if plan.actions == nil {
		plan.actions = []Action{}
	}
	plan.actions = append(plan.actions, action)
}

// MarshalJSON is used to marshal the private actions, which could not be written using reflection.
func (plan *Plan) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Actions []Action
	}{
		Actions: plan.actions,
	})
}

// UnmarshalJSON is used to unmarshal the private actions from a marshaled plan
func (plan *Plan) UnmarshalJSON(bytes []byte) error {
	unmarshaledJSON := struct {
		Actions []Action
	}{}
	err := json.Unmarshal(bytes, &unmarshaledJSON)
	if err != nil {
		return err
	}

	plan.actions = unmarshaledJSON.Actions
	return nil
}

// CreatePlan compares the model with the nexus and creates a plan, which describes action to get nexus in sync
// with the described model.
func CreatePlan(modelDAO ModelDAO, reader NexusAPIReader) (*Plan, error) {
	model, err := modelDAO.Get()
	if err != nil {
		return nil, errors.Wrap(err, "failed to read model")
	}

	plan := &Plan{}
	creator := &planCreator{
		model:  model,
		reader: reader,
		plan:   plan,
	}

	err = creator.createPlan()
	if err != nil {
		return nil, err
	}

	return plan, err
}

type planCreator struct {
	model  Model
	reader NexusAPIReader
	plan   *Plan
}

func (creator *planCreator) createPlan() error {
	for _, repository := range creator.model.Repositories {
		err := creator.createActionFor(repository)
		if err != nil {
			return err
		}
	}
	return nil
}

func (creator *planCreator) createActionFor(repository ModelRepository) error {
	clientRepository, err := creator.reader.Get(repository.ID)
	if err != nil {
		return errors.Wrapf(err, "failed to read repository %s from client api", repository.ID)
	}

	switch repository.State {
	case StatePresent:
		creator.handleStatePresent(repository.Repository, clientRepository)
	case StateAbsent:
		creator.handleStateAbsent(repository.Repository, clientRepository)
	}

	return nil
}

func (creator *planCreator) handleStatePresent(repository Repository, clientRepository *Repository) {
	if clientRepository != nil {
		if !repository.IsEqual(*clientRepository) {
			creator.plan.Modify(repository)
		}
	} else {
		creator.plan.Create(repository)
	}
}

func (creator *planCreator) handleStateAbsent(repository Repository, clientRepository *Repository) {
	if clientRepository != nil {
		creator.plan.Remove(repository)
	}
}
