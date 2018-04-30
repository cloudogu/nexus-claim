package infrastructure

// -s http://localhost:8082/nexus -u admin -p admin123 plan -i ./groovy/nexus-initial-example.hcl -o nexus-initial-example.json
// -s http://localhost:8081 -u admin -p admin123 plan -i ./groovy/nexus3-initial-example.hcl -o nexus3-initial-example.json


import (

	"net/http"

	"encoding/json"

	"log"

	"io"
	"bytes"
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

  repository, err := client.parseRepositoryJson(jsonData)

  if err != nil {
    return nil, err
  }

  return repository,nil

}

func (client *httpNexus3APIClient) createRepositoryURL(repositoryType domain.RepositoryType, id domain.RepositoryID) string {
	return client.createRepositoryServiceURL(repositoryType) + "/" + string(id)
}

func (client *httpNexus3APIClient) createReadRequest(url string) (*http.Request, error) {
	request, err := client.createRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Accept", contentType)
	return request, nil
}

func (client *httpNexus3APIClient) createRequest(method string, url string, body io.Reader) (*http.Request, error) {
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create %s request for %s", method, url)
	}

	if client.username != "" {
		request.SetBasicAuth(client.username, client.password)
	}

	return request, nil
}

func (client *httpNexus3APIClient) closeAndLogOnError(closer io.Closer) {
	err := closer.Close()
	if err != nil {
		log.Printf("failed to close response body: %v", err)
	}
}

func (client *httpNexus3APIClient) isStatusNotFound(response *http.Response) bool {
	return response.StatusCode == 404
}

func (client *httpNexus3APIClient) isStatusOK(response *http.Response) bool {
	return response.StatusCode == 200
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
	dto := newRepositoryDTO().from(repository)
	request, err := client.createWriteRequest("POST", client.createRepositoryServiceURL(repository.Type), dto)
	if err != nil {
		return err
	}

	response, err := client.httpClient.Do(request)
	if err != nil {
		return errors.Wrapf(err, "failed to create repository %s", repository.ID)
	}

	if response.StatusCode != 201 {
		return client.statusCodeError(response.StatusCode)
	}

	return nil
}

func (client *httpNexus3APIClient) createRepositoryServiceURL(repositoryType domain.RepositoryType) string {
	var typeURLPart string
	switch repositoryType {
	case domain.TypeRepository:
		typeURLPart = "repositories"
	case domain.TypeGroup:
		typeURLPart = "repo_groups"
	}
	return client.url + repositoryServiceURL + typeURLPart
}

func (client *httpNexus3APIClient) createWriteRequest(method, url string, body interface{}) (*http.Request, error) {
	reader, err := client.createJSONBody(body)
	if err != nil {
		return nil, err
	}

	request, err := client.createRequest(method, url, reader)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", contentType)
	return request, nil
}

func (client *httpNexus3APIClient) createJSONBody(object interface{}) (io.Reader, error) {
	data, err := json.Marshal(object)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal object")
	}

	return bytes.NewBuffer(data), nil
}

func (client *httpNexus3APIClient) Modify(repository domain.Repository) error {
	dto := newRepositoryDTO().from(repository)
	request, err := client.createWriteRequest("PUT", client.createRepositoryURL(repository.Type, repository.ID), dto)
	if err != nil {
		return err
	}

	response, err := client.httpClient.Do(request)
	if err != nil {
		return errors.Wrapf(err, "failed to modify repository %s", repository.ID)
	}

	if response.StatusCode != 200 {
		return client.statusCodeError(response.StatusCode)
	}

	return nil
}

func (client *httpNexus3APIClient) statusCodeError(statusCode int) error {
	return errors.Errorf("invalid status code %d", statusCode)
}

func (client *httpNexus3APIClient) Remove(repository domain.Repository) error {
	id := repository.ID

	request, err := client.createRequest("DELETE", client.createRepositoryURL(repository.Type, id), nil)
	if err != nil {
		return err
	}

	response, err := client.httpClient.Do(request)
	if err != nil {
		return errors.Wrapf(err, "failed to remove repository %s", id)
	}

	if response.StatusCode != 204 {
		return client.statusCodeError(response.StatusCode)
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
