package cmd

import (
	"github.com/pkg/errors"
	cli "gopkg.in/urfave/cli.v2"
)

// init registers the subcommand show to the application
func init() {
	application.registerCommand(createApplyCommand(application.apply))
}

func createApplyCommand(actionFunc cli.ActionFunc) cli.Command {
	return cli.Command{
		Name:    "apply",
		Aliases: []string{"a"},
		Usage:   "Applies a plan to the nexus api",
		Action:  actionFunc,
	}
}

func (app *Application) apply(c *cli.Context) error {
	return errors.New("not yet implemented")
}
