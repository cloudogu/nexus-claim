package cmd

import (
	"fmt"
	"os"

	"encoding/json"
	"io/ioutil"

	"github.com/cloudogu/nexus-claim/domain"
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
	planPath := c.Args().First()
	if planPath == "" {
		return cli.NewExitError("plan is required", 1)
	}

	plan, err := createPlanFromPath(planPath)
	if err != nil {
		return err
	}

	nexusAPIClient := app.createNexusAPIClient(c)
	err = domain.ApplyPlan(nexusAPIClient, plan)
	if err != nil {
		return errors.Wrapf(err, "failed to execute plan")
	}

	return nil
}

func createPlanFromPath(path string) (*domain.Plan, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, cliError("could not find plan at %s", path)
	}

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, cliError("failed to read plan %s: %v", path, err)
	}

	plan := &domain.Plan{}
	err = json.Unmarshal(bytes, plan)
	if err != nil {
		return nil, cliError("failed to unmarshal plan %s: %v", path, err)
	}

	return plan, nil
}

func cliError(format string, args ...interface{}) error {
	return cli.NewExitError(fmt.Sprintf(format, args...), 1)
}
