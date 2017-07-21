package cmd

import (
	"bytes"

	"testing"

	"io/ioutil"
	"os"

	"github.com/cloudogu/nexus-claim/domain"
	"github.com/cloudogu/nexus-claim/infrastructure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	cli "gopkg.in/urfave/cli.v2"
)

func TestApplyWithoutPlan(t *testing.T) {
	_, ec, err := execApply("apply")
	assert.Equal(t, 1, ec)
	assert.EqualError(t, err, "plan is required")
}

func TestApplyWithNonExistingPlan(t *testing.T) {
	_, ec, err := execApply("apply", "some/non/existing/file")
	assert.Equal(t, 1, ec)
	assert.EqualError(t, err, "could not find plan at some/non/existing/file")
}

func TestApplyWithInvalidPlan(t *testing.T) {
	temp, err := ioutil.TempFile("", "plan")
	require.Nil(t, err)
	defer os.Remove(temp.Name())

	_, err = temp.Write([]byte("invalid json"))
	require.Nil(t, err)

	_, ec, err := execApply("apply", temp.Name())
	assert.Equal(t, 1, ec)
	assert.Contains(t, err.Error(), "failed to unmarshal plan "+temp.Name())
}

func TestApply(t *testing.T) {
	temp, err := ioutil.TempFile("", "plan")
	require.Nil(t, err)
	defer os.Remove(temp.Name())

	plan := &domain.Plan{}
	plan.Create(domain.Repository{ID: domain.RepositoryID("crepo")})
	plan.Modify(domain.Repository{ID: domain.RepositoryID("mrepo")})
	plan.Remove(domain.Repository{ID: domain.RepositoryID("rrepo")})

	serializedPlan, err := infrastructure.SerializePlan(plan)
	require.Nil(t, err)
	_, err = temp.Write(serializedPlan)
	require.Nil(t, err)

	nexusApiClient := &mockNexusAPIClient{}
	_, ec, err := execApplyWithNexusApiClient(nexusApiClient, "apply", temp.Name())
	require.Nil(t, err)
	require.Equal(t, 0, ec)

	require.Equal(t, 1, len(nexusApiClient.Created))
	assert.Equal(t, domain.RepositoryID("crepo"), nexusApiClient.Created[0].ID)
	require.Equal(t, 1, len(nexusApiClient.Modified))
	assert.Equal(t, domain.RepositoryID("mrepo"), nexusApiClient.Modified[0].ID)
	require.Equal(t, 1, len(nexusApiClient.Removed))
	assert.Equal(t, domain.RepositoryID("rrepo"), nexusApiClient.Removed[0])
}

func execApply(args ...string) (*bytes.Buffer, int, error) {
	return execApplyWithNexusApiClient(&mockNexusAPIClient{}, args...)
}

func execApplyWithNexusApiClient(nexusApiClient domain.NexusAPIClient, args ...string) (*bytes.Buffer, int, error) {
	// capture exitCode and do not exit
	exitCode := 0
	cli.OsExiter = func(ec int) {
		exitCode = ec
	}

	var buffer bytes.Buffer
	cmdApp := Application{
		Output:         &buffer,
		nexusAPIClient: nexusApiClient,
	}

	app := cli.NewApp()
	app.Commands = []cli.Command{
		createApplyCommand(cmdApp.apply),
	}

	// add addition arg to first index to the slice, because os.Args contains the path to
	// the application on index 0
	err := app.Run(append([]string{""}, args...))
	return &buffer, exitCode, err
}
