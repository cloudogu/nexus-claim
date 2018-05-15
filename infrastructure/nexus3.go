//go:generate go run ../scripts/generate.go ../infrastructure/groovy_scripts ../scripts
package infrastructure

import (
	"encoding/json"

	"github.com/cloudogu/nexus-claim/domain"
	"github.com/cloudogu/nexus-scripting/manager"
	"github.com/pkg/errors"
	"reflect"
	"strings"
)

// NewNexus3APIClient creates a new nexus3APIClient
func NewNexus3APIClient(url string, username string, password string) domain.NexusAPIClient {
	clientManager := manager.New(url, username, password)
	return &nexus3APIClient{url, username, password, clientManager}
}

type nexus3APIClient struct {
	url        string
	username   string
	password   string
	manager    *manager.Manager
}

func (client *nexus3APIClient) Get(repositoryType domain.RepositoryType, id domain.RepositoryID) (*domain.Repository, error) {

	stringID := string(id)

	script, err := client.manager.Create("readRepository", READ_REPOSITORY)
	if err != nil {
		return nil, err
	}

	jsonData, err := script.ExecuteWithStringPayload(stringID)
	if err != nil {
		return nil, err
	}
	if client.isStatusNotFound(jsonData) {
		return nil, nil
	}

	repository, err := client.JSONToRepository(jsonData)
	if err != nil {
		return nil, err
	}

	return repository, nil

}

func (client *nexus3APIClient) Create(repository domain.Repository) error {

	script, err := client.manager.Create("createRepository", CREATE_REPOSITORY)

	if err != nil {
		return err
	}

	readAbleJSON, err := client.repositoryToJSON(repository)
  if err != nil {
		return err
	}
	output, err := script.ExecuteWithStringPayload(readAbleJSON)
	if err != nil {
		return err
	}
	if strings.Contains(output, "exception") {
		return errors.New(output)
	}

	return nil
}

func (client *nexus3APIClient) Modify(repository domain.Repository) error {

	script, err := client.manager.Create("modifyRepository", MODIFY_REPOSITORY)
	if err != nil {
		return err
	}

	readAbleJSON, err := client.repositoryToJSON(repository)
  if err != nil {
		return err
	}

	output, err := script.ExecuteWithStringPayload(readAbleJSON)
	if err != nil {
		return err
	}
	if strings.Contains(output, "exception") {
		return errors.Wrapf(err, "error: %s", output)
	}

	return nil
}

func (client *nexus3APIClient) Remove(repository domain.Repository) error {

  stringID := string(repository.ID)
	script, err := client.manager.Create("deleteRepository", DELETE_REPOSITORY)
	if err != nil {
		return err
	}

	output, err := script.ExecuteWithStringPayload(stringID)
	if err != nil {
		return err
	}
  if strings.Contains(output, "exception") {
    return errors.Wrapf(err, "error: %s", output)
  }

	return nil
}

func (client *nexus3APIClient) isStatusNotFound(output string) bool {
	return strings.Contains(output,"404: no repository found")
}

func (client *nexus3APIClient) repositoryToJSON(repository domain.Repository) (string, error) {

	jsonData, err := json.Marshal(repository.Properties)
	if err != nil {
		return "", errors.Wrap(err, "failed to marshal the json data")
	}
	readAbleJSON := string(jsonData)

	return readAbleJSON, nil
}

func (client *nexus3APIClient) JSONToRepository(jsonData string) (*domain.Repository, error) {

  dto := newNexus3RepositoryDTO()
  dto, err := dto.from(jsonData)
  if err != nil {
    return nil, err
  }
  return dto.to(), nil
}

func newNexus3RepositoryDTO() *nexus3RepositoryDTO {
	return &nexus3RepositoryDTO{}
}

type nexus3RepositoryDTO struct {
	Data domain.Properties
}

func (dto *nexus3RepositoryDTO) from(jsonData string) (*nexus3RepositoryDTO, error) {

	var jsonMap map[string]interface{}
	err := json.Unmarshal([]byte(jsonData), &jsonMap)
	if err != nil {
		return nil, err
	}

	dto.Data = jsonMap
	dto.Data["id"] = jsonMap["id"].(string)

	return dto, nil
}

func (dto *nexus3RepositoryDTO) to() *domain.Repository {
	properties := dto.convertFloatToInt()
	return &domain.Repository{
		ID:         domain.RepositoryID(dto.Data["id"].(string)),
		Properties: properties,
		Type:       domain.TypeRepository,
	}
}

func (dto *nexus3RepositoryDTO) convertFloatToInt() domain.Properties {
	properties := make(domain.Properties)
	for key, value := range dto.Data {
		if reflect.TypeOf(value).Kind() == reflect.Float64 {
			properties[key] = int(value.(float64))
		} else {
			properties[key] = value
		}
	}
	return properties
}
