package cmd

import (
	"encoding/json"
	"io/ioutil"

	"github.com/cloudogu/nexus-claim/domain"
	"github.com/pkg/errors"
	cli "gopkg.in/urfave/cli.v2"
)

// init registers the subcommand show to the application
func init() {
	application.registerCommand(createShowCommand(application.show))
}

func createShowCommand(actionFunc cli.ActionFunc) cli.Command {
	return cli.Command{
		Name:    "show",
		Aliases: []string{"s"},
		Usage:   "Reads a plan and prints it to stdout",
		Action:  actionFunc,
	}
}

func (app *Application) show(c *cli.Context) error {
	planFile := c.Args().First()
	if planFile == "" {
		return errors.New("usage: nexus-claim show path/to/plan.json")
	}

	bytes, err := ioutil.ReadFile(planFile)
	if err != nil {
		return errors.Wrapf(err, "failed to read plan file %s", planFile)
	}

	plan := &domain.Plan{}
	err = json.Unmarshal(bytes, plan)
	if err != nil {
		return errors.Wrapf(err, "failed to unmarshal plan file %s", planFile)
	}

	err = app.printPlan(plan)
	if err != nil {
		return err
	}

	return nil
}
