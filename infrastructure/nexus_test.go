package infrastructure

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"io/ioutil"

	"github.com/cloudogu/nexus-claim/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHttpNexusAPIClient_Get(t *testing.T) {
	server := servce(t, "/service/local/repositories/test-repo", "simple.json")
	defer server.Close()

	client := NewHTTPNexusAPIClient(server.URL, "admin", "admin123")
	repository, err := client.Get(domain.RepositoryID("test-repo"))
	require.Nil(t, err)

	assert.Equal(t, domain.RepositoryID("test-repo"), repository.ID)
	assert.Equal(t, "Simple test repository", repository.Properties["name"])
}

func TestHttpNexusAPIClient_GetNotFound(t *testing.T) {
	server := servce(t, "/service/local/repositories/test-repo", "simple.json")
	defer server.Close()

	client := NewHTTPNexusAPIClient(server.URL, "admin", "admin123")
	repository, err := client.Get(domain.RepositoryID("non-existing-repo"))
	require.Nil(t, err)
	require.Nil(t, repository)
}

func TestHttpNexusAPIClient_GetInvalidStatusCode(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(503)
	}))
	defer server.Close()

	client := NewHTTPNexusAPIClient(server.URL, "admin", "admin123")
	repository, err := client.Get(domain.RepositoryID("some-repo"))
	require.NotNil(t, err)
	require.Nil(t, repository)
	require.Contains(t, err.Error(), "invalid status code 503")
}

func TestHttpNexusAPIClient_GetInvalidBody(t *testing.T) {
	server := servce(t, "/service/local/repositories/invalid-body", "invalid-body.json")
	defer server.Close()

	client := NewHTTPNexusAPIClient(server.URL, "admin", "admin123")
	repository, err := client.Get(domain.RepositoryID("invalid-body"))
	require.NotNil(t, err)
	require.Nil(t, repository)
	require.Contains(t, err.Error(), "failed to unmarshal response body")
}

func TestHttpNexusAPIClient_GetAcceptHeader(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "application/json; charset=UTF-8", r.Header.Get("Accept"))
		w.WriteHeader(404)
	}))
	defer server.Close()

	client := NewHTTPNexusAPIClient(server.URL, "admin", "admin123")
	client.Get(domain.RepositoryID("accept"))
}

func TestHttpNexusAPIClient_GetAuthentication(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		assert.Equal(t, "admin", username)
		assert.Equal(t, "admin123", password)
		assert.True(t, ok)
		w.WriteHeader(404)
	}))
	defer server.Close()

	client := NewHTTPNexusAPIClient(server.URL, "admin", "admin123")
	client.Get(domain.RepositoryID("some-repo"))
}

func servce(t *testing.T, url string, filename string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == url {
			w.WriteHeader(200)
			bytes, err := ioutil.ReadFile("../resources/" + filename)
			require.Nil(t, err)
			w.Write(bytes)
			return
		}
		w.WriteHeader(404)
	}))
}
