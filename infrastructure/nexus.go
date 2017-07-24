package infrastructure

import (
	"net/http"

	"io/ioutil"

	"encoding/json"

	"log"

	"io"

	"bytes"

	"github.com/cloudogu/nexus-claim/domain"
	"github.com/pkg/errors"
)

const (
	repositoryServiceURL = "/service/local/"
	contentType          = "application/json; charset=UTF-8"
)

// NewHTTPNexusAPIClient creates a new http based NexusAPIClient
func NewHTTPNexusAPIClient(url string, username string, password string) domain.NexusAPIClient {
	return &httpNexusAPIClient{url, username, password, &http.Client{}}
}

type httpNexusAPIClient struct {
	url        string
	username   string
	password   string
	httpClient *http.Client
}

func (client *httpNexusAPIClient) Get(repositoryType domain.RepositoryType, id domain.RepositoryID) (*domain.Repository, error) {
	repositoryURL := client.createRepositoryURL(repositoryType, id)
	request, err := client.createReadRequest(repositoryURL)
	if err != nil {
		return nil, err
	}

	response, err := client.httpClient.Do(request)
	if err != nil {
		return nil, errors.Wrapf(err, "get request for repository %s failed", id)
	}
	defer client.closeAndLogOnError(response.Body)

	if client.isStatusNotFound(response) {
		return nil, nil
	}

	if !client.isStatusOK(response) {
		return nil, client.statusCodeError(response.StatusCode)
	}

	repository, err := client.parseRepositoryResponse(response)
	if err != nil {
		return nil, err
	}
	repository.Type = repositoryType
	return repository, nil
}

func (client *httpNexusAPIClient) createRepositoryURL(repositoryType domain.RepositoryType, id domain.RepositoryID) string {
	return client.createRepositoryServiceURL(repositoryType) + "/" + string(id)
}

func (client *httpNexusAPIClient) createReadRequest(url string) (*http.Request, error) {
	request, err := client.createRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Accept", contentType)
	return request, nil
}

func (client *httpNexusAPIClient) createRequest(method string, url string, body io.Reader) (*http.Request, error) {
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create %s request for %s", method, url)
	}

	if client.username != "" {
		request.SetBasicAuth(client.username, client.password)
	}

	return request, nil
}

func (client *httpNexusAPIClient) closeAndLogOnError(closer io.Closer) {
	err := closer.Close()
	if err != nil {
		log.Printf("failed to close response body: %v", err)
	}
}

func (client *httpNexusAPIClient) isStatusNotFound(response *http.Response) bool {
	return response.StatusCode == 404
}

func (client *httpNexusAPIClient) isStatusOK(response *http.Response) bool {
	return response.StatusCode == 200
}

func (client *httpNexusAPIClient) parseRepositoryResponse(response *http.Response) (*domain.Repository, error) {
	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse response body")
	}

	dto := newRepositoryDTO()
	err = json.Unmarshal(content, dto)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal response body")
	}

	return dto.to(), nil
}

func newRepositoryDTO() *repositoryDTO {
	return &repositoryDTO{}
}

type repositoryDTO struct {
	Data domain.Properties `json:"data"`
}

func (dto *repositoryDTO) from(repository domain.Repository) *repositoryDTO {
	properties := repository.Properties
	if properties == nil {
		properties = make(domain.Properties)
	}
	dto.Data = properties
	dto.Data["id"] = repository.ID
	return dto
}

func (dto *repositoryDTO) to() *domain.Repository {
	return &domain.Repository{
		ID:         domain.RepositoryID(dto.Data["id"].(string)),
		Properties: dto.Data,
		Type:       domain.TypeRepository,
	}
}

func (client *httpNexusAPIClient) Create(repository domain.Repository) error {
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

func (client *httpNexusAPIClient) createRepositoryServiceURL(repositoryType domain.RepositoryType) string {
	var typeUrlPart string
	switch repositoryType {
	case domain.TypeRepository:
		typeUrlPart = "repositories"
	case domain.TypeGroup:
		typeUrlPart = "repo_groups"
	}
	return client.url + repositoryServiceURL + typeUrlPart
}

func (client *httpNexusAPIClient) createWriteRequest(method, url string, body interface{}) (*http.Request, error) {
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

func (client *httpNexusAPIClient) createJSONBody(object interface{}) (io.Reader, error) {
	data, err := json.Marshal(object)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal object")
	}

	return bytes.NewBuffer(data), nil
}

func (client *httpNexusAPIClient) Modify(repository domain.Repository) error {
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

func (client *httpNexusAPIClient) statusCodeError(statusCode int) error {
	return errors.Errorf("invalid status code %d", statusCode)
}

func (client *httpNexusAPIClient) Remove(repository domain.Repository) error {
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
