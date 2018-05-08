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

func TestHttpNexus3APIClient_Get(t *testing.T) {
  server := servceReadRepository(t,  "simpleAnswerFromReadRepository.json")
  defer server.Close()

  client := NewNexus3APIClient(server.URL, "admin", "admin123")
  repository, err := client.Get(domain.TypeRepository, domain.RepositoryID("testRepo"))
  require.Nil(t, err)

  assert.Equal(t, "testRepo", string(repository.ID))
}

func TestHttpNexus3APIClient_Create(t *testing.T){
  server := serveCreateRepository(t,  "simpleAnswerFromCreateRepository.json")
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

    //w.WriteHeader(404)
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
