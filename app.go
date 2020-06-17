package main

import (
	"gopkg.in/urfave/cli.v2"

	"os"

	"github.com/cloudogu/nexus-claim/cmd"
)

var (
	// Version of the application
	Version = "0.3.1"

	// CommitID git sha1 hash
	CommitID string
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
		panic(err.Error())
	}
}

func createVersion(version string, commitID string) string {
	v := version
	length := len(commitID)
	if length > 0 {
		limit := 8
		if length < 8 {
			limit = length
		}
		v += "-" + commitID[:limit]
	}
	return v
}
