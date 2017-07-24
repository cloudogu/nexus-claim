package infrastructure

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"io/ioutil"

	"encoding/json"

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
	assert.Equal(t, domain.RepositoryTypeRepository, repository.Type)
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

func TestHttpNexusAPIClient_Create(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/service/local/repositories", r.URL.Path)
		assert.Equal(t, "POST", r.Method)

		defer r.Body.Close()

		bytes, err := ioutil.ReadAll(r.Body)
		require.Nil(t, err)

		jsonBody := make(map[string]interface{})
		err = json.Unmarshal(bytes, &jsonBody)
		require.Nil(t, err)

		data, ok := jsonBody["data"].(map[string]interface{})
		require.True(t, ok)

		assert.Equal(t, "test-repo", data["id"])
		assert.Equal(t, "Test Repository", data["name"])

		w.WriteHeader(201)
	}))
	defer server.Close()

	client := NewHTTPNexusAPIClient(server.URL, "admin", "admin123")

	properties := make(domain.Properties)
	properties["name"] = "Test Repository"
	repository := domain.Repository{
		ID:         domain.RepositoryID("test-repo"),
		Properties: properties,
	}
	err := client.Create(repository)
	require.Nil(t, err)
}

func TestHttpNexusAPIClient_CreateWithWrongStatusCode(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(502)
	}))
	defer server.Close()

	client := NewHTTPNexusAPIClient(server.URL, "admin", "admin123")
	err := client.Create(domain.Repository{ID: domain.RepositoryID("test")})
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "invalid status code 502")
}

func TestHttpNexusAPIClient_CreateContentType(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "application/json; charset=UTF-8", r.Header.Get("Content-Type"))
		w.WriteHeader(201)
	}))
	defer server.Close()

	client := NewHTTPNexusAPIClient(server.URL, "admin", "admin123")
	err := client.Create(domain.Repository{ID: domain.RepositoryID("test")})
	require.Nil(t, err)
}

func TestHttpNexusAPIClient_CreateAuthentication(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		assert.Equal(t, "admin", username)
		assert.Equal(t, "admin123", password)
		assert.True(t, ok)

		w.WriteHeader(201)
	}))
	defer server.Close()

	client := NewHTTPNexusAPIClient(server.URL, "admin", "admin123")
	err := client.Create(domain.Repository{ID: domain.RepositoryID("test")})
	require.Nil(t, err)
}

func TestHttpNexusAPIClient_ModifyWithWrongStatusCode(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
	}))
	defer server.Close()

	client := NewHTTPNexusAPIClient(server.URL, "admin", "admin123")
	err := client.Modify(domain.Repository{ID: domain.RepositoryID("test")})
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "invalid status code 400")
}

func TestHttpNexusAPIClient_Modify(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/service/local/repositories/test-repo", r.URL.Path)
		assert.Equal(t, "PUT", r.Method)

		defer r.Body.Close()

		bytes, err := ioutil.ReadAll(r.Body)
		require.Nil(t, err)

		jsonBody := make(map[string]interface{})
		err = json.Unmarshal(bytes, &jsonBody)
		require.Nil(t, err)

		data, ok := jsonBody["data"].(map[string]interface{})
		require.True(t, ok)

		assert.Equal(t, "test-repo", data["id"])
		assert.Equal(t, "Test Repository", data["name"])

		w.WriteHeader(200)
	}))
	defer server.Close()

	client := NewHTTPNexusAPIClient(server.URL, "admin", "admin123")

	properties := make(domain.Properties)
	properties["name"] = "Test Repository"
	repository := domain.Repository{
		ID:         domain.RepositoryID("test-repo"),
		Properties: properties,
	}
	err := client.Modify(repository)
	require.Nil(t, err)
}

func TestHttpNexusAPIClient_RemoveWithInvalidStatusCode(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
	}))
	defer server.Close()

	client := NewHTTPNexusAPIClient(server.URL, "admin", "admin123")
	err := client.Remove(domain.Repository{ID: domain.RepositoryID("test-repo")})
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "invalid status code 404")
}

func TestHttpNexusAPIClient_Remove(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/service/local/repositories/test-repo", r.URL.Path)
		assert.Equal(t, "DELETE", r.Method)
		w.WriteHeader(204)
	}))
	defer server.Close()

	client := NewHTTPNexusAPIClient(server.URL, "admin", "admin123")
	err := client.Remove(domain.Repository{ID: domain.RepositoryID("test-repo")})
	require.Nil(t, err)
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
