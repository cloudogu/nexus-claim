package domain

import "github.com/pkg/errors"

// ApplyPlan executes a plan against the nexus api
func ApplyPlan(writer NexusAPIWriter, plan *Plan) error {
	for _, action := range plan.actions {
		err := applyAction(writer, action)
		if err != nil {
			return err
		}
	}
	return nil
}

func applyAction(writer NexusAPIWriter, action Action) error {
	repository := action.GetRepository()
	switch action.GetType() {
	case ActionCreate:
		return applyActionFunc(repository, writer.Create)
	case ActionModify:
		return applyActionFunc(repository, writer.Modify)
	case ActionRemove:
		return applyRemove(writer, repository.ID)
	}
	return nil
}

type actionFunc func(repository Repository) error

func applyActionFunc(repository Repository, actionFunc actionFunc) error {
	err := actionFunc(repository)
	if err != nil {
		return errors.Wrapf(err, "repository action failed %s", repository.ID)
	}
	return nil
}

func applyRemove(writer NexusAPIWriter, id RepositoryID) error {
	err := writer.Remove(id)
	if err != nil {
		return errors.Wrapf(err, "failed to remove repository %s", id)
	}
	return nil
}
