package infrastructure

import (
	"encoding/json"

	"github.com/cloudogu/nexus-claim/domain"
	"github.com/pkg/errors"
)

// SerializePlan is used to serialize a plan to a format which can be written to any kind of storage
func SerializePlan(plan *domain.Plan) ([]byte, error) {
	jp := &jsonPlan{}
	for _, action := range plan.GetActions() {
		jp.add(action)
	}
	serializedPlan, err := json.Marshal(jp)
	if err != nil {
		return nil, errors.Wrap(err, "failed to serialize plan")
	}
	return serializedPlan, nil
}

// DeserializePlan is used to deserialize a plan
func DeserializePlan(serializedPlan []byte) (*domain.Plan, error) {
	jp := &jsonPlan{}
	err := json.Unmarshal(serializedPlan, &jp)
	if err != nil {
		return nil, errors.Wrap(err, "failed to deserialize plan")
	}

	plan := &domain.Plan{}
	for _, action := range jp.Actions {
		plan.AddAction(action.Type, action.Repository)
	}
	return plan, nil
}

type jsonPlan struct {
	Actions []jsonAction
}

func (jp *jsonPlan) add(action domain.Action) {
	if jp.Actions == nil {
		jp.Actions = make([]jsonAction, 0)
	}
	jp.Actions = append(jp.Actions, jsonAction{
		Type:       action.GetType(),
		Repository: action.GetRepository(),
	})
}

type jsonAction struct {
	Type       domain.ActionType
	Repository domain.Repository
}
