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

type actionExecutor func(writer NexusAPIWriter, repository Repository) error

type defaultAction struct {
	Type       ActionType
	Repository Repository
	executor   actionExecutor
}

func (action *defaultAction) GetType() ActionType {
	return action.Type
}
func (action *defaultAction) GetRepository() Repository {
	return action.Repository
}

func (action *defaultAction) Execute(writer NexusAPIWriter) error {
	return action.executor(writer, action.Repository)
}

// Action describes a single unit of work to sync the model with nexus
type Action interface {
	GetType() ActionType
	GetRepository() Repository
	Execute(writer NexusAPIWriter) error
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
	plan.action(ActionCreate, repository, func(writer NexusAPIWriter, repo Repository) error {
		err := writer.Create(repo)
		if err != nil {
			return errors.Wrapf(err, "failed to create repository %s", repo.ID)
		}
		return nil
	})
}

// Modify adds a modify action to the plan
func (plan *Plan) Modify(repository Repository) {
	plan.action(ActionModify, repository, func(writer NexusAPIWriter, repo Repository) error {
		err := writer.Modify(repo)
		if err != nil {
			return errors.Wrapf(err, "failed to modify repository %s", repo.ID)
		}
		return nil
	})
}

// Remove adds a remove action to the plan
func (plan *Plan) Remove(repository Repository) {
	plan.action(ActionRemove, repository, func(writer NexusAPIWriter, repo Repository) error {
		err := writer.Remove(repo.ID)
		if err != nil {
			return errors.Wrapf(err, "failed to remove repository %s", repo.ID)
		}
		return nil
	})
}

// AddAction to a plan
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

func (plan *Plan) action(actionType ActionType, repository Repository, executor actionExecutor) {
	action := &defaultAction{Type: actionType, Repository: repository, executor: executor}
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
			mergedRepository := clientRepository.Merge(repository)
			creator.plan.Modify(mergedRepository)
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
