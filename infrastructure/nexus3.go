//go:generate go run ../scripts/generate.go ../infrastructure/groovy_scripts ../scripts
package infrastructure


import (

	"net/http"

	"encoding/json"

	"github.com/cloudogu/nexus-claim/domain"
	"github.com/pkg/errors"
  "github.com/cloudogu/nexus-scripting/manager"
  "fmt"
  "reflect"
  "strings"
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

  fmt.Println("getting " + StringID)

  jsonData, err := script.ExecuteWithStringPayload(StringID)
  if err != nil {
    return nil, err
  }
  if client.isStatusNotFound(jsonData){
    return nil, nil
  }

  repository, err := client.parseRepositoryJson(jsonData)
  if err != nil {
    return nil, err
  }

  return repository,nil

}

func (client *httpNexus3APIClient) parseRepositoryJson(jsonData string) (*domain.Repository, error) {

  dto := newRepository3DTO()

  dto,err := dto.from(jsonData)

  if err != nil {
    return nil,err
  }

	return dto.to(), nil
}

func (client *httpNexus3APIClient) Create(repository domain.Repository) error {

  createRepositoryScript := CREATE_REPOSITORY;
  script,err := client.manager.Create("createRepository",createRepositoryScript)
  if err != nil {
    return err
  }

  fmt.Println("creating " + repository.ID)

  readAbleJson, err := repositoryToJson(repository)
  if err != nil {
    return err
  }

  output, err := script.ExecuteWithStringPayload(readAbleJson)
  if err != nil {
    return err
  }

  if !(strings.Contains(output,"successfully")){
    return errors.Wrapf(err, "error: %s", output)
  }

  return nil
}

func (client *httpNexus3APIClient) Modify(repository domain.Repository) error {

  modifyRepositoryScript := MODIFY_REPOSITORY
  script, err := client.manager.Create("modifyRepository",modifyRepositoryScript)
  if err != nil {
    return err
  }

  fmt.Println("modifying " + repository.ID)

  readAbleJson, err := repositoryToJson(repository)
  if err != nil {
    return err
  }


  output, err := script.ExecuteWithStringPayload(readAbleJson)
  if err != nil {
    return err
  }
  if !(strings.Contains(output,"successfully")){
    return errors.Wrapf(err, "error: %s", output)
  }

  return nil
}

func (client *httpNexus3APIClient) Remove(repository domain.Repository) error {
  
  deleteRepositoryScript := DELETE_REPOSITORY;
  StringID := string(repository.ID)
  script,err := client.manager.Create("deleteRepository",deleteRepositoryScript)
  if err != nil {
    return err
  }

  fmt.Println("deleting " + StringID)

  _, err = script.ExecuteWithStringPayload(StringID)
  if err != nil {
    return err
  }

  return nil
}

func (client *httpNexus3APIClient) isStatusNotFound(output string) bool {
  return strings.Contains(output,"404")
}

func repositoryToJson(repository domain.Repository) (string, error){
  
  jsonData,err := json.Marshal(repository.Properties)
  if err != nil {
    return "", errors.Wrap(err, "failed to marshal the json data")
  }
  readAbleJson:= string(jsonData)
  
  return readAbleJson,nil
}

func newRepository3DTO() *repository3DTO {
  return &repository3DTO{}
}

type repository3DTO struct {
  Data domain.Properties
}

func (dto *repository3DTO) from(jsonData string) (*repository3DTO,error) {

  var jsonMap map[string]interface{}

  err := json.Unmarshal([]byte(jsonData),&jsonMap)
  if err != nil {
    return nil,err
  }

  dto.Data = jsonMap
  dto.Data["id"] = jsonMap["id"].(string)

  return dto,nil
}

func (dto *repository3DTO) to() *domain.Repository {
  properties := dto.convertFloatToInt()
  return &domain.Repository{
    ID:         domain.RepositoryID(dto.Data["id"].(string)),
    Properties: properties,
    Type:       domain.TypeRepository,
  }
}

func (dto *repository3DTO) convertFloatToInt() domain.Properties {
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
