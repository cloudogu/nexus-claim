package infrastructure

//go:generate go run ../scripts/generate.go ../infrastructure/groovy_scripts ../scripts

import (
	"encoding/json"
	"github.com/cloudogu/nexus-claim/domain"
	"github.com/cloudogu/nexus-scripting/manager"
	"github.com/pkg/errors"
	"reflect"
	"strings"
)

// NewNexus3APIClient creates a new nexus3APIClient
func NewNexus3APIClient(url string, username string, password string, timeout int) domain.NexusAPIClient {
	clientManager := manager.New(url, username, password)
	clientManager.WithTimeout(timeout)
	return &nexus3APIClient{url, username, password, clientManager}
}

type nexus3APIClient struct {
	url      string
	username string
	password string
	manager  *manager.Manager
}

func (client *nexus3APIClient) Get(repositoryType domain.RepositoryType, id domain.RepositoryID) (*domain.Repository, error) {

	stringID := string(id)

	script, err := client.manager.Create("readRepository", READ_REPOSITORY)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create readRepository.groovy with %s", stringID)
	}

	jsonData, err := script.ExecuteWithStringPayload(stringID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to execute readRepository.groovy with %s", stringID)
	}

	if client.isStatusNotFound(jsonData) {
		return nil, nil
	}

	repository, err := client.JSONToRepository(jsonData)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse JSON file of %s", stringID)
	}

	return repository, nil

}

func (client *nexus3APIClient) Create(repository domain.Repository) error {
	script, err := client.manager.Create("createRepository", CREATE_REPOSITORY)
	if err != nil {
		return errors.Wrapf(err, "failed to create createRepository.groovy from %s", repository.ID)
	}

	enrichedRepository := client.addRepositoryNamesFromID(repository)
	enrichedRepository, err = client.addRepoInfosFromRecipeName(enrichedRepository)
	if err != nil {
		return errors.Wrapf(err, "failed to create repository %s", enrichedRepository.ID)
	}

	readAbleJSON, err := client.repositoryToJSON(enrichedRepository)
	if err != nil {
		return errors.Wrapf(err, "failed to parse to JSON from repository %s to create it", enrichedRepository.ID)
	}

	output, err := script.ExecuteWithStringPayload(readAbleJSON)
	if err != nil {
		return errors.Wrapf(err, "failed to execute createRepository.groovy with %s", enrichedRepository.ID)
	}

	if !(output == "null") {
		return errors.Errorf("createRepository.groovy with repository %s executed but an error occured: %s", enrichedRepository.ID, output)
	}

	return nil
}

func (client *nexus3APIClient) Modify(repository domain.Repository) error {
	script, err := client.manager.Create("modifyRepository", MODIFY_REPOSITORY)
	if err != nil {
		return errors.Wrapf(err, "failed to create modifyRepository.groovy from %s", repository.ID)
	}

	enrichedRepository := client.addRepositoryNamesFromID(repository)
	enrichedRepository, err = client.addRepoInfosFromRecipeName(enrichedRepository)
	if err != nil {
		return errors.Wrapf(err, "failed to modify repository %s", enrichedRepository.ID)
	}

	readAbleJSON, err := client.repositoryToJSON(enrichedRepository)
	if err != nil {
		return errors.Wrapf(err, "failed to parse to JSON from repository %s to modify it", enrichedRepository.ID)
	}

	output, err := script.ExecuteWithStringPayload(readAbleJSON)
	if err != nil {
		return errors.Wrapf(err, "failed to execute modifyRepository.groovy with %s", string(enrichedRepository.ID))
	}
	if !(output == "null") {
		return errors.Errorf("modifyRepository.groovy with repository %s executed but an error occured: %s", enrichedRepository.ID, output)
	}

	return nil
}

// addRepositoryNamesFromID takes a repository and adds repository (meta) data which are redundantly(if absent) from the
// repository ID. If this meta data was missing it would render the repository completely unmaintainable. This is
// because the Nexus 3 API allows repository creation without meta data.
func (client *nexus3APIClient) addRepositoryNamesFromID(repository domain.Repository) domain.Repository {
	id := string(repository.ID)

	newRepo := repository.Clone()
	newRepo.Properties["id"] = id
	newRepo.Properties["repositoryName"] = id

	return newRepo
}

// addRepoInfosFromRecipeName takes a repository and adds repository meta data about the format and type of the
// repository (f. i. "docker-hosted") without explicitly stating it additionally.
func (client *nexus3APIClient) addRepoInfosFromRecipeName(repository domain.Repository) (domain.Repository, error) {
	recipeName, err := repository.GetRecipeName()
	if err != nil {
		return domain.Repository{}, err
	}

	recipePartDelimiter := "-"
	recipeNameParts := strings.Split(recipeName, recipePartDelimiter)
	repoFormat := recipeNameParts[0]
	repoType := recipeNameParts[1]

	newRepo := repository.Clone()
	newRepo.Properties["format"] = repoFormat
	newRepo.Properties["type"] = repoType

	return newRepo, nil
}

func (client *nexus3APIClient) Remove(repository domain.Repository) error {
	stringID := string(repository.ID)
	script, err := client.manager.Create("deleteRepository", DELETE_REPOSITORY)
	if err != nil {
		return errors.Wrapf(err, "failed to create deleteRepository.groovy from %s", repository.ID)
	}

	output, err := script.ExecuteWithStringPayload(stringID)
	if err != nil {
		return errors.Wrapf(err, "failed to execute deleteRepository.groovy with %s", string(repository.ID))
	}
	if !(output == "null") {
		return errors.Errorf("deleteRepository.groovy with repository %s executed but an error occured: %s", repository.ID, output)
	}

	return nil
}

func (client *nexus3APIClient) isStatusNotFound(output string) bool {
	return strings.Contains(output, "404: no repository found")
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
		return nil, errors.Wrap(err, "Error occurred during parsing JSON data into a Nexus3DPO")
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
		return nil, errors.Wrap(err, "Error occured at unmarshalling jsonData")
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
