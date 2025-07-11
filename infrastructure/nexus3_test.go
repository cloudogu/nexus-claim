package infrastructure

import (
  "github.com/cloudogu/nexus-claim/domain"
  "github.com/stretchr/testify/assert"
  "github.com/stretchr/testify/require"
  "io/ioutil"
  "net/http"
  "net/http/httptest"
  "testing"
)

func TestNexus3APIClient_Get(t *testing.T) {
  t.Run("read repository - 'null' result", func(t *testing.T) {
    server := serveReadRepository(t, "answerFromReadRepository.json")
    defer server.Close()

    client := NewNexus3APIClient(server.URL, "admin", "admin123", 30)
    repository, err := client.Get(domain.TypeRepository, domain.RepositoryID("testRepo"))
    require.Nil(t, err)

    assert.Equal(t, "testRepo", string(repository.ID))
  })
  t.Run("read repository - '' result", func(t *testing.T) {
    server := serveReadRepository(t, "answer2FromReadRepository.json")
    defer server.Close()

    client := NewNexus3APIClient(server.URL, "admin", "admin123", 30)
    repository, err := client.Get(domain.TypeRepository, domain.RepositoryID("testRepo"))
    require.Nil(t, err)

    assert.Equal(t, "testRepo", string(repository.ID))
  })

}

func TestNexus3APIClient_GetNotFound(t *testing.T) {
  server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    fullPath := r.Method + " " + r.URL.Path

    if fullPath == "GET /service/rest/v1/script/readRepository" {
      w.WriteHeader(404)
      return
    }

  }))
  defer server.Close()

  client := NewNexus3APIClient(server.URL, "admin", "admin123", 30)
  _, err := client.Get(domain.TypeRepository, domain.RepositoryID("testRepo"))
  assert.NotNil(t, err)

}

func TestNexus3APIClient_Create(t *testing.T) {
  t.Run("create repository - 'null' result", func(t *testing.T) {
    server := serveCreateRepository(t, "answerFromCreateRepository.json")
    defer server.Close()

    client := NewNexus3APIClient(server.URL, "admin", "admin123", 30)

    properties := make(domain.Properties)
    properties["id"] = "testRepo"
    properties["recipeName"] = "docker-hosted"
    repository := domain.Repository{
      ID:         domain.RepositoryID("testRepo"),
      Properties: properties,
      Type:       domain.TypeRepository,
    }
    err := client.Create(repository)
    require.Nil(t, err)
  })

  t.Run("create repository - '' result", func(t *testing.T) {
    server := serveCreateRepository(t, "answer2FromCreateRepository.json")
    defer server.Close()

    client := NewNexus3APIClient(server.URL, "admin", "admin123", 30)

    properties := make(domain.Properties)
    properties["id"] = "testRepo"
    properties["recipeName"] = "docker-hosted"
    repository := domain.Repository{
      ID:         domain.RepositoryID("testRepo"),
      Properties: properties,
      Type:       domain.TypeRepository,
    }
    err := client.Create(repository)
    require.Nil(t, err)
  })
}

func TestNexus3APIClient_Create_error_no_recipeName(t *testing.T) {
  server := serveCreateRepository(t, "answerFromCreateRepository.json")
  defer server.Close()

  client := NewNexus3APIClient(server.URL, "admin", "admin123", 30)

  properties := make(domain.Properties)
  properties["id"] = "testRepo"
  repository := domain.Repository{
    ID:         domain.RepositoryID("testRepo"),
    Properties: properties,
    Type:       domain.TypeRepository,
  }

  err := client.Create(repository)

  assert.NotNil(t, err)
  assert.Contains(t, err.Error(), "could not find property 'recipeName' in repository")
}

func TestNexus3APIClient_CreateWithWrongStatusCode(t *testing.T) {
  server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(502)
  }))
  defer server.Close()

  client := NewNexus3APIClient(server.URL, "admin", "admin123", 30)
  err := client.Create(domain.Repository{ID: domain.RepositoryID("test")})
  require.NotNil(t, err)
  require.Contains(t, err.Error(), "502")
}

func TestNexus3APIClient_Delete(t *testing.T) {
  t.Run("delete repository - 'null' result", func(t *testing.T) {
    server := serveDeleteRepository(t, "answerFromDeleteRepository.json")
    defer server.Close()

    client := NewNexus3APIClient(server.URL, "admin", "admin123", 30)

    err := client.Remove(domain.Repository{ID: domain.RepositoryID("test-repo")})
    require.Nil(t, err)
  })

  t.Run("delete repository - '' result", func(t *testing.T) {
    server := serveDeleteRepository(t, "answer2FromDeleteRepository.json")
    defer server.Close()

    client := NewNexus3APIClient(server.URL, "admin", "admin123", 30)

    err := client.Remove(domain.Repository{ID: domain.RepositoryID("test-repo")})
    require.Nil(t, err)
  })
}

func TestNexus3APIClient_DeleteWithInvalidStatusCode(t *testing.T) {
  server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(404)
  }))
  defer server.Close()

  client := NewNexus3APIClient(server.URL, "admin", "admin123", 30)
  err := client.Remove(domain.Repository{ID: domain.RepositoryID("test-repo")})
  require.NotNil(t, err)
  require.Contains(t, err.Error(), "404")
}

func TestNexus3APIClient_Modify(t *testing.T) {
  t.Run("modify repository - 'null' result", func(t *testing.T) {
    server := serveModifyRepository(t, "answerFromModifyRepository.json")
    defer server.Close()

    client := NewNexus3APIClient(server.URL, "admin", "admin123", 30)

    properties := make(domain.Properties)
    properties["id"] = "testRepo"
    properties["recipeName"] = "recipe-name"
    repository := domain.Repository{
      ID:         domain.RepositoryID("testRepo2"),
      Properties: properties,
      Type:       domain.TypeRepository,
    }

    err := client.Modify(repository)

    assert.Equal(t, "testRepo2", string(repository.ID))

    require.Nil(t, err)
  })

  t.Run("modify repository - '' result", func(t *testing.T) {
    server := serveModifyRepository(t, "answer2FromModifyRepository.json")
    defer server.Close()

    client := NewNexus3APIClient(server.URL, "admin", "admin123", 30)

    properties := make(domain.Properties)
    properties["id"] = "testRepo"
    properties["recipeName"] = "recipe-name"
    repository := domain.Repository{
      ID:         domain.RepositoryID("testRepo2"),
      Properties: properties,
      Type:       domain.TypeRepository,
    }

    err := client.Modify(repository)

    assert.Equal(t, "testRepo2", string(repository.ID))

    require.Nil(t, err)
  })
}

func TestNexus3APIClient_ModifyWithWrongStatusCode(t *testing.T) {
  server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(400)
  }))
  defer server.Close()

  client := NewNexus3APIClient(server.URL, "admin", "admin123", 30)
  err := client.Modify(domain.Repository{ID: domain.RepositoryID("test")})
  require.NotNil(t, err)
  require.Contains(t, err.Error(), "400")
}

func TestNexus3APIClient_addRepositoryNamesFromID(t *testing.T) {
  server := serveCreateRepository(t, "answerFromCreateRepository.json")
  defer server.Close()

  client := nexus3APIClient{url: server.URL, username: "admin", password: "admin123"}

  properties := make(domain.Properties)
  properties["recipeName"] = "maven2-hosted"
  nestedProperty := make(map[string]interface{})
  properties["someMavenProperty"] = nestedProperty

  repository := domain.Repository{
    ID:         domain.RepositoryID("testRepo"),
    Properties: properties,
    Type:       domain.TypeRepository,
  }

  actualRepo := client.addRepositoryNamesFromID(repository)

  assert.NotNil(t, actualRepo)
  assert.Contains(t, actualRepo.Properties, "id")
  assert.Contains(t, actualRepo.Properties, "repositoryName")
  assert.Equal(t, "testRepo", actualRepo.Properties["id"])
  assert.Equal(t, "testRepo", actualRepo.Properties["repositoryName"])
}

func TestNexus3APIClient_addRepoInfosFromRecipeName(t *testing.T) {
  server := serveCreateRepository(t, "answerFromCreateRepository.json")
  defer server.Close()

  client := nexus3APIClient{url: server.URL, username: "admin", password: "admin123"}

  properties := make(domain.Properties)
  properties["recipeName"] = "docker-hosted"
  repository := domain.Repository{
    ID:         domain.RepositoryID("testRepo"),
    Properties: properties,
    Type:       domain.TypeRepository,
  }
  actualRepo, err := client.addRepoInfosFromRecipeName(repository)

  assert.Nil(t, err)
  assert.NotNil(t, actualRepo)
  assert.Contains(t, actualRepo.Properties, "type")
  assert.Contains(t, actualRepo.Properties, "format")
  assert.Equal(t, "hosted", actualRepo.Properties["type"])
  assert.Equal(t, "docker", actualRepo.Properties["format"])
}

func TestNexus3APIClient_addRepoInfosFromRecipeName_error(t *testing.T) {
  server := serveCreateRepository(t, "answerFromCreateRepository.json")
  defer server.Close()

  client := nexus3APIClient{url: server.URL, username: "admin", password: "admin123"}

  properties := make(domain.Properties)
  repository := domain.Repository{
    ID:         domain.RepositoryID("testRepo"),
    Properties: properties,
    Type:       domain.TypeRepository,
  }
  _, err := client.addRepoInfosFromRecipeName(repository)

  assert.NotNil(t, err)
  assert.Contains(t, err.Error(), "could not find property 'recipeName' in repository testRepo")
}

func serveCreateRepository(t *testing.T, filename string) *httptest.Server {
  return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

    fullPath := r.Method + " " + r.URL.Path

    if fullPath == "PUT /service/rest/v1/script/createRepository" {
      w.WriteHeader(204)
      return
    } else if fullPath == "GET /service/rest/v1/script/createRepository" {
      w.WriteHeader(200)
      return
    } else if fullPath == "POST /service/rest/v1/script/createRepository/run" {
      w.WriteHeader(200)
      bytes, err := ioutil.ReadFile("../resources/nexus3/" + filename)
      require.Nil(t, err)
      _, _ = w.Write(bytes)
      return
    }

  }))
}

func serveReadRepository(t *testing.T, filename string) *httptest.Server {
  return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

    fullPath := r.Method + " " + r.URL.Path
    if fullPath == "POST /service/rest/v1/script" {
      w.WriteHeader(204)
      return
    } else if fullPath == "GET /service/rest/v1/script/readRepository" {
      w.WriteHeader(200)
      return
    } else if fullPath == "PUT /service/rest/v1/script/readRepository" {
      w.WriteHeader(204)
      return
    } else if fullPath == "POST /service/rest/v1/script/readRepository/run" {
      w.WriteHeader(200)
      bytes, err := ioutil.ReadFile("../resources/nexus3/" + filename)
      require.Nil(t, err)
      _, _ = w.Write(bytes)
      return
    }

  }))
}

func serveDeleteRepository(t *testing.T, filename string) *httptest.Server {
  return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

    fullPath := r.Method + " " + r.URL.Path

    if fullPath == "PUT /service/rest/v1/script/deleteRepository" {
      w.WriteHeader(204)
      return
    } else if fullPath == "GET /service/rest/v1/script/deleteRepository" {
      w.WriteHeader(200)
      return
    } else if fullPath == "POST /service/rest/v1/script/deleteRepository/run" {
      w.WriteHeader(200)
      bytes, err := ioutil.ReadFile("../resources/nexus3/" + filename)
      require.Nil(t, err)
      _, _ = w.Write(bytes)
      return
    }

  }))
}

func serveModifyRepository(t *testing.T, filename string) *httptest.Server {
  return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

    fullPath := r.Method + " " + r.URL.Path

    if fullPath == "PUT /service/rest/v1/script/modifyRepository" {
      w.WriteHeader(204)
      return
    } else if fullPath == "GET /service/rest/v1/script/modifyRepository" {
      w.WriteHeader(200)
      return
    } else if fullPath == "POST /service/rest/v1/script/modifyRepository/run" {
      w.WriteHeader(200)
      bytes, err := ioutil.ReadFile("../resources/nexus3/" + filename)
      require.Nil(t, err)
      _, _ = w.Write(bytes)
      return
    }

  }))
}
