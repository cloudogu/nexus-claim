package main

import (
	cli "gopkg.in/urfave/cli.v2"

	"os"

	"github.com/cloudogu/nexus-claim/cmd"
)

func main() {
	cmdApp := cmd.GetApplication()

	app := cli.NewApp()
	app.Name = "nexus-claim"
	app.Usage = "Define your Sonatype Nexus repository structure as code"
	app.Flags = cmdApp.GlobalFlags()
	app.Commands = cmdApp.Commands()

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
