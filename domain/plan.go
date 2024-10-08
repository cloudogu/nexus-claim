package domain

import (
	"fmt"
)

// Plan is a set of actions which must be done to sync the model with nexus
type Plan struct {
	actions []Action
}

// GetActions returns the planned actions
func (plan *Plan) GetActions() []Action {
	return plan.actions
}

// AddAction adds an action to the plan
func (plan *Plan) AddAction(actionType ActionType, repository Repository) {
	switch actionType {
	case ActionCreate:
		plan.Create(repository)
	case ActionModify:
		plan.Modify(repository)
	case ActionRemove:
		plan.Remove(repository)
	}
}

// Create adds a create action to the plan
func (plan *Plan) Create(repository Repository) {
	plan.appendAction(&createAction{baseAction{Type: ActionCreate, Repository: repository}})
}

// Merge merges the repository which was fetched from the nexus api with the one from the model
// and creates a modify action with the merged repository.
func (plan *Plan) Merge(clientRepository Repository, modelRepository Repository) error {
	mergedRepository, err := clientRepository.Merge(modelRepository)
	if err != nil {
		return fmt.Errorf("failed to merge repository %s: %w", modelRepository.ID, err)
	}
	plan.Modify(mergedRepository)
	return nil
}

// Modify adds a modify action to the plan
func (plan *Plan) Modify(repository Repository) {
	plan.appendAction(&modifyAction{baseAction{Type: ActionModify, Repository: repository}})
}

// Remove adds a remove action to the plan
func (plan *Plan) Remove(repository Repository) {
	plan.appendAction(&removeAction{baseAction{Type: ActionRemove, Repository: repository}})
}

func (plan *Plan) appendAction(action Action) {
	if plan.actions == nil {
		plan.actions = []Action{}
	}
	plan.actions = append(plan.actions, action)
}

// Execute executes every action against the nexus api
func (plan *Plan) Execute(writer NexusAPIWriter) error {
	for _, action := range plan.GetActions() {
		err := action.Execute(writer)
		if err != nil {
			return err
		}
	}
	return nil
}

// CreatePlan compares the model with the nexus and creates a plan, which describes actions to get nexus in sync
// with the described model.
func CreatePlan(modelDAO ModelDAO, reader NexusAPIReader) (*Plan, error) {
	model, err := modelDAO.Get()
	if err != nil {
		return nil, fmt.Errorf("failed to read model: %w", err)
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
	clientRepository, err := creator.reader.Get(repository.Type, repository.ID)
	if err != nil {
		return fmt.Errorf("failed to read repository %s from client api: %w", repository.ID, err)
	}

	switch repository.State {
	case StatePresent:
		err := creator.handleStatePresent(repository.Repository, clientRepository)
		if err != nil {
			return err
		}
	case StateAbsent:
		creator.handleStateAbsent(repository.Repository, clientRepository)
	}

	return nil
}

func (creator *planCreator) handleStatePresent(repository Repository, clientRepository *Repository) error {
	if clientRepository != nil {
		if !repository.IsEqual(*clientRepository) {
			err := creator.plan.Merge(*clientRepository, repository)
			if err != nil {
				return err
			}
		}
	} else {
		creator.plan.Create(repository)
	}
	return nil
}

func (creator *planCreator) handleStateAbsent(repository Repository, clientRepository *Repository) {
	if clientRepository != nil {
		creator.plan.Remove(repository)
	}
}
