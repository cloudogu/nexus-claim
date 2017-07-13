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
	repositoryServiceURL = "/service/local/repositories"
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

func (client *httpNexusAPIClient) Get(id domain.RepositoryID) (*domain.Repository, error) {
	repositoryURL := client.createRepositoryURL(id)
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
		return nil, errors.Errorf("invalid status code %d returned", response.StatusCode)
	}

	return client.parseRepositoryResponse(response)
}

func (client *httpNexusAPIClient) createRepositoryURL(id domain.RepositoryID) string {
	return client.createRepositoryServiceURL() + "/" + string(id)
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
	}
}

func (client *httpNexusAPIClient) Create(repository domain.Repository) error {
	dto := newRepositoryDTO().from(repository)
	request, err := client.createWriteRequest("POST", client.createRepositoryServiceURL(), dto)
	if err != nil {
		return errors.Wrap(err, "failed to create post request")
	}

	response, err := client.httpClient.Do(request)
	if err != nil {
		return errors.Wrapf(err, "failed to create repository %s", repository.ID)
	}

	if response.StatusCode != 201 {
		return errors.Errorf("invalid status code %d", response.StatusCode)
	}

	return nil
}

func (client *httpNexusAPIClient) createRepositoryServiceURL() string {
	return client.url + repositoryServiceURL
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
	return errors.New("not yet implemented")
}

func (client *httpNexusAPIClient) Remove(id domain.RepositoryID) error {
	return errors.New("not yet implemented")
}
