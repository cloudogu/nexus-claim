//go:generate go run ../scripts/generate.go ../infrastructure/groovy_scripts ../scripts
package infrastructure

// -s http://localhost:8082/nexus -u admin -p admin123 plan -i ./groovy/nexus-initial-example.hcl -o nexus-initial-example.json
// -s http://localhost:8081 -u admin -p admin123 plan -i ./groovy/nexus3-initial-example.hcl -o nexus3-initial-example.json


import (

	"net/http"

	"encoding/json"

	"github.com/cloudogu/nexus-claim/domain"
	"github.com/pkg/errors"
  "github.com/cloudogu/nexus-scripting/manager"
  "fmt"
  "reflect"
)


// NewHTTPNexus3APIClient creates a new Nexus3APIClient
func NewHTTPNexus3APIClient(url string, username string, password string) domain.NexusAPIClient {
  manager := manager.New(url,username,password)
  return &httpNexus3APIClient{url, username, password, &http.Client{},manager}
}

type httpNexus3APIClient struct {
	url        string
	username   string
	password   string
	httpClient *http.Client
	manager    *manager.Manager
}

func (client *httpNexus3APIClient) Get(repositoryType domain.RepositoryType, id domain.RepositoryID) (*domain.Repository, error) {

  readRepositoryScript := READ_REPOSITORY;
  StringID := string(id)
  script,err := client.manager.Create("readRepository",readRepositoryScript)

  if err != nil {
    return nil, err
  }

  jsonData, err := script.ExecuteWithStringPayload(StringID)

  if err != nil {
    return nil, err
  }

  // Todo ParseRepositoryJson fertig stellen
  repository, err := client.parseRepositoryJson(jsonData)

  if err != nil {
    return nil, err
  }

  return repository,nil

}

func (client *httpNexus3APIClient) parseRepositoryJson(jsonData string) (*domain.Repository, error) {

  fmt.Println(jsonData)

	dto := newRepository3DTO()
	bytes,err := json.Marshal(jsonData)
  if err != nil {
    return nil, errors.Wrap(err, "failed to marshal the json data")
  }

  err = json.Unmarshal(bytes, dto)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal the json data")
	}
	fmt.Println(dto.attributes)

	return dto.to(), nil
}


func (client *httpNexus3APIClient) Create(repository domain.Repository) error {

  createRepositoryScript := CREATE_REPOSITORY;
  script,err := client.manager.Create("createRepository",createRepositoryScript)
  if err != nil {
    return err
  }

  jsonData,err := json.Marshal(repository.Properties)
  if err != nil {
    return errors.Wrap(err, "failed to marshal the json data")
  }

  readAbleJson:= string(jsonData)

  fmt.Println(readAbleJson)

  _, err = script.ExecuteWithStringPayload(readAbleJson)
  if err != nil {
    return err
  }

  return nil
}

func (client *httpNexus3APIClient) Modify(repository domain.Repository) error {
	return nil
}

func (client *httpNexus3APIClient) Remove(repository domain.Repository) error {

  deleteRepositoryScript := DELETE_REPOSITORY;
  StringID := string(repository.ID)
  script,err := client.manager.Create("deleteRepository",deleteRepositoryScript)
  if err != nil {
    return err
  }

  _, err = script.ExecuteWithStringPayload(StringID)
  if err != nil {
    return err
  }

  return nil
}

func newRepository3DTO() *repository3DTO {
  return &repository3DTO{}
}

type repository3DTO struct {
  Format string
  Id string
  Types string
  Url string
  attributes domain.Properties
}

func (dto *repository3DTO) from(repository domain.Repository) *repository3DTO {
  properties := repository.Properties
  if properties == nil {
    properties = make(domain.Properties)
  }
  dto.attributes = properties
  dto.Id = string(repository.ID)
  return dto
}

func (dto *repository3DTO) to() *domain.Repository {
  properties := dto.convertFloatToInt()
  return &domain.Repository{
    ID:         domain.RepositoryID(dto.Id),
    Properties: properties,
    Type:       domain.TypeRepository,
  }
}

func (dto *repository3DTO) convertFloatToInt() domain.Properties {
  properties := make(domain.Properties)
  for key, value := range dto.attributes {
    if reflect.TypeOf(value).Kind() == reflect.Float64 {
      properties[key] = int(value.(float64))
    } else {
      properties[key] = value
    }
  }
  return properties
}
