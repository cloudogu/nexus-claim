package cmd

import (
	"testing"

	"bytes"

	"bufio"

	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/cloudogu/nexus-claim/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/urfave/cli.v2"
)

func TestPlan(t *testing.T) {
	buffer, err := execPlan("plan", "-i", "../resources/nexus-initial-example.hcl")
	require.Nil(t, err)

	scanner := bufio.NewScanner(buffer)
	require.True(t, scanner.Scan())
	assert.Equal(t, "- apache-snapshots", scanner.Text())
	require.True(t, scanner.Scan())
	assert.Equal(t, "- central-m1", scanner.Text())
	require.True(t, scanner.Scan())
	assert.Equal(t, "~ thirdparty", scanner.Text())
	require.True(t, scanner.Scan())
	assert.Equal(t, "+ scm-releases", scanner.Text())
}

func TestPlanWithQuietParameter(t *testing.T) {
	buffer, err := execPlan("plan", "-q", "-i", "../resources/nexus-initial-example.hcl")
	require.Nil(t, err)
	assert.Equal(t, 0, buffer.Len())
}

func TestPlanWriteOutput(t *testing.T) {
	file, err := ioutil.TempFile("", "nc-plan")
	require.Nil(t, err)

	defer os.Remove(file.Name())

	_, err = execPlan("plan", "-q", "-i", "../resources/detail.hcl", "-o", file.Name())
	require.Nil(t, err)

	bytes, err := ioutil.ReadAll(file)
	require.Nil(t, err)

	plan := &domain.Plan{}
	err = json.Unmarshal(bytes, plan)
	require.Nil(t, err)

	require.Equal(t, 1, len(plan.GetActions()))
	action := plan.GetActions()[0]

	assert.Equal(t, domain.ActionCreate, action.Type)
	assert.Equal(t, domain.RepositoryID("releases"), action.Repository.ID)
	assert.Equal(t, "Releases", action.Repository.Properties["Name"])
}

func execPlan(args ...string) (*bytes.Buffer, error) {
	var buffer bytes.Buffer
	cmdApp := Application{
		Output:         &buffer,
		nexusAPIClient: &mockNexusAPIClient{},
	}

	app := cli.NewApp()
	app.Commands = []cli.Command{
		createPlanCommand(cmdApp.plan),
	}

	// add addition arg to first index to the slice, because os.Args contains the path to
	// the application on index 0
	err := app.Run(append([]string{""}, args...))
	return &buffer, err
}

type mockNexusAPIClient struct {
	repositories map[domain.RepositoryID]*domain.Repository
}

func (mock *mockNexusAPIClient) add(repo domain.Repository) {
	mock.repositories[repo.ID] = &repo
}

func (mock *mockNexusAPIClient) init() {
	if mock.repositories == nil {
		mock.repositories = make(map[domain.RepositoryID]*domain.Repository)
		mock.add(domain.Repository{
			ID: domain.RepositoryID("apache-snapshots"),
		})
		mock.add(domain.Repository{
			ID: domain.RepositoryID("central-m1"),
		})
		properies := make(domain.Properties)
		properies["name"] = "3rd Party"
		mock.add(domain.Repository{
			ID:         domain.RepositoryID("thirdparty"),
			Properties: properies,
		})
	}
}

func (mock *mockNexusAPIClient) Get(id domain.RepositoryID) (*domain.Repository, error) {
	mock.init()
	return mock.repositories[id], nil
}
