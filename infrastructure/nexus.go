package infrastructure

import (
	"net/http"

	"io/ioutil"

	"encoding/json"

	"log"

	"io"

	"github.com/cloudogu/nexus-claim/domain"
	"github.com/pkg/errors"
)

const (
	repositoryServiceURL = "/service/local/repositories/"
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
	request, err := client.createGetRequest(repositoryURL)
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
	return client.url + repositoryServiceURL + string(id)
}

func (client *httpNexusAPIClient) createGetRequest(url string) (*http.Request, error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create get request for %s", url)
	}

	request.Header.Set("Accept", "application/json; charset=UTF-8")

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

	nexusAPIResponse := nexusAPIResponse{}
	err = json.Unmarshal(content, &nexusAPIResponse)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal response body")
	}

	properties := nexusAPIResponse.Data
	return &domain.Repository{
		ID:         domain.RepositoryID(properties["id"].(string)),
		Properties: properties,
	}, nil
}

type nexusAPIResponse struct {
	Data domain.Properties
}
