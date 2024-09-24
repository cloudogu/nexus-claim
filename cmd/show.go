package cmd

import (
  "fmt"
  "io/ioutil"

  "github.com/cloudogu/nexus-claim/infrastructure"
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
		return fmt.Errorf("usage: nexus-claim show path/to/plan.json")
	}

	serializedPlan, err := ioutil.ReadFile(planFile)
	if err != nil {
		return fmt.Errorf("failed to read plan file %s: %w", planFile, err)
	}

	plan, err := infrastructure.DeserializePlan(serializedPlan)
	if err != nil {
		return fmt.Errorf("failed to unmarshal plan file %s: %w", planFile, err)
	}

	err = app.printPlan(plan)
	if err != nil {
		return err
	}

	return nil
}
