package cmd

import (
	"testing"

	"bytes"

	"bufio"

	"io/ioutil"
	"os"

	"github.com/cloudogu/nexus-claim/domain"
	"github.com/cloudogu/nexus-claim/infrastructure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/urfave/cli.v2"
)

func TestPlan(t *testing.T) {
	buffer, err, ec := execPlan("plan", "-i", getResourcesDir()+"nexus-initial-example.hcl")
	require.Nil(t, err)
	require.Equal(t, 0, ec)

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
	buffer, err, ec := execPlan("plan", "-q", "-i", getResourcesDir()+"nexus-initial-example.hcl")
	require.Nil(t, err)
	require.Equal(t, 0, ec)
	assert.Equal(t, 0, buffer.Len())
}

func TestPlanWriteOutput(t *testing.T) {
	file, err := ioutil.TempFile("", "nc-plan")
	require.Nil(t, err)

	defer os.Remove(file.Name())

	_, err, ec := execPlan("plan", "-q", "-i", getResourcesDir()+"detail.hcl", "-o", file.Name())
	require.Nil(t, err)
	require.Equal(t, 0, ec)

	serializedPlan, err := ioutil.ReadAll(file)
	require.Nil(t, err)

	plan, err := infrastructure.DeserializePlan(serializedPlan)
	require.Nil(t, err)

	require.Equal(t, 1, len(plan.GetActions()))
	action := plan.GetActions()[0]

	assert.Equal(t, domain.ActionCreate, action.GetType())
	assert.Equal(t, domain.RepositoryID("releases"), action.GetRepository().ID)
	assert.Equal(t, "Releases", action.GetRepository().Properties["Name"])
}

func TestPlanWriteToStdOut(t *testing.T) {
	buffer, err, ec := execPlan("plan", "-i", getResourcesDir()+"detail.hcl", "-o", "-")
	require.Nil(t, err)
	require.Equal(t, 0, ec)

	plan, err := infrastructure.DeserializePlan(buffer.Bytes())
	require.Nil(t, err)

	actions := plan.GetActions()
	require.Equal(t, 1, len(actions))
}

func execPlan(args ...string) (*bytes.Buffer, error, int) {
	// capture exitCode and do not exit
	exitCode := 0
	cli.OsExiter = func(ec int) {
		exitCode = ec
	}

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
	return &buffer, err, exitCode
}

func getResourcesDir() string {
	return "../resources/nexus2/"
}
