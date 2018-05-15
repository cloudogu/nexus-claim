package infrastructure

import (
  "github.com/stretchr/testify/require"
  "testing"
  "github.com/cloudogu/nexus-claim/domain"
  "github.com/stretchr/testify/assert"
  "net/http/httptest"
  "net/http"
  "io/ioutil"
)

func TestNexus3APIClient_Get(t *testing.T) {
  server := servceReadRepository(t,  "answerFromReadRepository.json")
  defer server.Close()

  client := NewNexus3APIClient(server.URL, "admin", "admin123")
  repository, err := client.Get(domain.TypeRepository, domain.RepositoryID("testRepo"))
  require.Nil(t, err)

  assert.Equal(t, "testRepo", string(repository.ID))
}

func TestNexus3APIClient_GetNotFound(t *testing.T) {
  server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    fullPath := r.Method + " " + r.URL.Path

    if fullPath == "GET /service/rest/v1/script/readRepository"{
      w.WriteHeader(404)
      return
    }

  }))
  defer server.Close()

  client := NewNexus3APIClient(server.URL, "admin", "admin123")
  _, err := client.Get(domain.TypeRepository, domain.RepositoryID("testRepo"))
  assert.NotNil(t,err)

}

func TestNexus3APIClient_Create(t *testing.T){
  server := serveCreateRepository(t,  "answerFromCreateRepository.json")
  defer server.Close()

  client := NewNexus3APIClient(server.URL, "admin", "admin123")

  properties := make(domain.Properties)
  properties["id"] = "testRepo"
  repository := domain.Repository{
    ID:         domain.RepositoryID("testRepo"),
    Properties: properties,
    Type:       domain.TypeRepository,
  }
  err := client.Create(repository)
  require.Nil(t, err)
}

func TestNexus3APIClient_CreateWithWrongStatusCode(t *testing.T) {
  server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(502)
  }))
  defer server.Close()

  client := NewNexus3APIClient(server.URL, "admin", "admin123")
  err := client.Create(domain.Repository{ID: domain.RepositoryID("test")})
  require.NotNil(t, err)
  require.Contains(t, err.Error(), "502")
}

func TestNexus3APIClient_Delete(t *testing.T){
  server := serveDeleteRepository(t,  "answerFromDeleteRepository.json")
  defer server.Close()

  client := NewNexus3APIClient(server.URL, "admin", "admin123")

  err := client.Remove(domain.Repository{ID: domain.RepositoryID("test-repo")})
  require.Nil(t, err)
}

func TestNexus3APIClient_DeleteWithInvalidStatusCode(t *testing.T) {
  server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(404)
  }))
  defer server.Close()

  client := NewNexus3APIClient(server.URL, "admin", "admin123")
  err := client.Remove(domain.Repository{ID: domain.RepositoryID("test-repo")})
  require.NotNil(t, err)
  require.Contains(t, err.Error(), "404")
}

func TestNexus3APIClient_Modify(t *testing.T){
  server := serveModifyRepository(t,  "answerFromModifyRepository.json")
  defer server.Close()

  client := NewNexus3APIClient(server.URL, "admin", "admin123")

  properties := make(domain.Properties)
  properties["id"] = "testRepo"
  repository := domain.Repository{
    ID:         domain.RepositoryID("testRepo2"),
    Properties: properties,
    Type:       domain.TypeRepository,
  }

  err := client.Modify(repository)

  assert.Equal(t,"testRepo2",string(repository.ID))

  require.Nil(t, err)

}

func TestNexus3APIClient_ModifyWithWrongStatusCode(t *testing.T) {
  server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(400)
  }))
  defer server.Close()

  client := NewNexus3APIClient(server.URL, "admin", "admin123")
  err := client.Modify(domain.Repository{ID: domain.RepositoryID("test")})
  require.NotNil(t, err)
  require.Contains(t, err.Error(), "400")
}

func serveCreateRepository(t *testing. T, filename string) *httptest.Server{
  return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

    fullPath := r.Method + " " + r.URL.Path

    if fullPath == "PUT /service/rest/v1/script/createRepository"{
      w.WriteHeader(204)
      return
    }
    if fullPath == "GET /service/rest/v1/script/createRepository"{
      w.WriteHeader(200)
      return
    }
    if fullPath == "POST /service/rest/v1/script/createRepository/run"{
      w.WriteHeader(200)
      bytes, err := ioutil.ReadFile("resources/nexus3/" + filename)
      require.Nil(t, err)
      w.Write(bytes)
      return
    }
    w.WriteHeader(404)

  }))}

func servceReadRepository(t *testing.T, filename string) *httptest.Server {
  return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

    fullPath := r.Method + " " + r.URL.Path
    if fullPath == "POST /service/rest/v1/script"{
      w.WriteHeader(204)
      return
    } else if fullPath == "GET /service/rest/v1/script/readRepository"{
      w.WriteHeader(200)
      return
    } else if fullPath == "PUT /service/rest/v1/script/readRepository"{
      w.WriteHeader(204)
      return
    } else if fullPath == "POST /service/rest/v1/script/readRepository/run"{
      w.WriteHeader(200)
      bytes, err := ioutil.ReadFile("resources/nexus3/" + filename)
      require.Nil(t, err)
      w.Write(bytes)
      return
    }

    w.WriteHeader(404)
  }))
}

func serveDeleteRepository(t *testing. T, filename string) *httptest.Server{
  return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

    fullPath := r.Method + " " + r.URL.Path

    if fullPath == "PUT /service/rest/v1/script/deleteRepository"{
      w.WriteHeader(204)
      return
    }
    if fullPath == "GET /service/rest/v1/script/deleteRepository"{
      w.WriteHeader(200)
      return
    }
    if fullPath == "POST /service/rest/v1/script/deleteRepository/run"{
      w.WriteHeader(200)
      bytes, err := ioutil.ReadFile("resources/nexus3/" + filename)
      require.Nil(t, err)
      w.Write(bytes)
      return
    }
    w.WriteHeader(404)

  }))}

func serveModifyRepository(t *testing. T, filename string) *httptest.Server{
  return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

    fullPath := r.Method + " " + r.URL.Path

    if fullPath == "PUT /service/rest/v1/script/modifyRepository"{
      w.WriteHeader(204)
      return
    }
    if fullPath == "GET /service/rest/v1/script/modifyRepository"{
      w.WriteHeader(200)
      return
    }
    if fullPath == "POST /service/rest/v1/script/modifyRepository/run"{
      w.WriteHeader(200)
      bytes, err := ioutil.ReadFile("resources/nexus3/" + filename)
      require.Nil(t, err)
      w.Write(bytes)
      return
    }
    w.WriteHeader(404)

  }))}
