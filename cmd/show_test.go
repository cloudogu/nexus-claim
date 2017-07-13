package cmd

import (
	"bytes"

	"testing"

	"io/ioutil"
	"os"

	"bufio"

	"github.com/cloudogu/nexus-claim/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	cli "gopkg.in/urfave/cli.v2"
)

func TestShow(t *testing.T) {
	file, err := ioutil.TempFile("", "ncp-")
	require.Nil(t, err)

	defer os.Remove(file.Name())

	plan := &domain.Plan{}
	plan.Create(domain.Repository{ID: domain.RepositoryID("abc")})

	bytes, err := plan.MarshalJSON()
	require.Nil(t, err)

	_, err = file.Write(bytes)
	require.Nil(t, err)

	buffer, err := execShow(file.Name())
	require.Nil(t, err)

	scanner := bufio.NewScanner(buffer)
	require.True(t, scanner.Scan())
	assert.Equal(t, "+ abc", scanner.Text())
	assert.False(t, scanner.Scan())
}

func execShow(args ...string) (*bytes.Buffer, error) {
	var buffer bytes.Buffer
	cmdApp := Application{
		Output:         &buffer,
		nexusAPIClient: &mockNexusAPIClient{},
	}

	app := cli.NewApp()
	app.Commands = []cli.Command{
		createShowCommand(cmdApp.show),
	}

	err := app.Run(append([]string{"", "show"}, args...))
	return &buffer, err
}
