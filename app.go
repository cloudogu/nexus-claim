package main

import (
	"fmt"

	"github.com/cloudogu/nexus-claim/domain"
	"github.com/cloudogu/nexus-claim/infrastructure"
)

func main() {
	dao := createFileModelDAO()
	client := createNexusAPIClient()

	plan, err := domain.CreatePlan(dao, client)
	if err != nil {
		panic(err)
	}

	for _, action := range plan.GetActions() {
		fmt.Print(createOperatorFromActionType(action.Type))
		fmt.Print(" ")
		fmt.Println(action.Repository.ID)
	}
}

func createOperatorFromActionType(actionType domain.ActionType) string {
	switch actionType {
	case domain.ActionCreate:
		return "+"
	case domain.ActionModify:
		return "~"
	case domain.ActionRemove:
		return "-"
	default:
		return "#"
	}
}

func createNexusAPIClient() domain.NexusAPIClient {
	return infrastructure.NewHTTPNexusAPIClient("http://localhost:8081/nexus", "admin", "admin123")
}

func createFileModelDAO() domain.ModelDAO {
	return infrastructure.NewFileModelDAO("resources/nexus-initial-example.hcl")
}
