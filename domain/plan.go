package domain

import "github.com/pkg/errors"

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
	plan.actions = append(plan.actions, plan.action(ActionCreate, repository))
}

// Modify adds a modify action to the plan
func (plan *Plan) Modify(repository Repository) {
	plan.actions = append(plan.actions, plan.action(ActionModify, repository))
}

// Remove adds a remove action to the plan
func (plan *Plan) Remove(repository Repository) {
	plan.actions = append(plan.actions, plan.action(ActionRemove, repository))
}

func (plan *Plan) action(actionType ActionType, repository Repository) Action {
	return Action{Type: actionType, Repository: repository}
}

// CreatePlan compares the model with the nexus and creates a plan, which describes action to get nexus in sync
// with the described model.
func CreatePlan(modelDAO ModelDAO, client NexusAPIClient) (*Plan, error) {
	model, err := modelDAO.Get()
	if err != nil {
		return nil, errors.Wrap(err, "failed to read model")
	}

	plan := newPlan()
	creator := &planCreator{
		model:  model,
		client: client,
		plan:   plan,
	}

	err = creator.createPlan()
	if err != nil {
		return nil, err
	}

	return plan, err
}

func newPlan() *Plan {
	plan := Plan{}
	plan.actions = []Action{}
	return &plan
}

type planCreator struct {
	model  Model
	client NexusAPIClient
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
	clientRepository, err := creator.client.Get(repository.ID)
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
