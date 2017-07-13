package main

import (
	cli "gopkg.in/urfave/cli.v2"

	"os"

	"github.com/cloudogu/nexus-claim/cmd"
)

var (
	// Version of the application
	Version string = "0.0.0"

	// CommitID git sha1 hash
	CommitID string = ""
)

func main() {
	cmdApp := cmd.GetApplication()

	app := cli.NewApp()
	app.Name = "nexus-claim"
	app.Usage = "Define your Sonatype Nexus repository structure as code"
	app.Version = createVersion(Version, CommitID)
	app.Flags = cmdApp.GlobalFlags()
	app.Commands = cmdApp.Commands()

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}

func createVersion(version string, commitId string) string {
	v := version
	length := len(commitId)
	if length > 0 {
		limit := 8
		if length < 8 {
			limit = length
		}
		v += "-" + commitId[:limit]
	}
	return v
}
