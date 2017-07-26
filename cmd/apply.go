package cmd

import (
	"os"

	"io/ioutil"

	"github.com/cloudogu/nexus-claim/domain"
	"github.com/cloudogu/nexus-claim/infrastructure"
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
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "input, i",
				Usage: "`PLAN` to apply, use '-' to read from stdin",
			},
		},
	}
}

func (app *Application) apply(c *cli.Context) error {
	planInput := c.String("input")
	if planInput == "" {
		return cli.NewExitError("plan is required", 1)
	}

	plan, err := app.createPlan(planInput)
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

func (app *Application) createPlan(input string) (*domain.Plan, error) {
	serializedPlan, err := app.readPlan(input)
	if err != nil {
		return nil, err
	}

	plan, err := infrastructure.DeserializePlan(serializedPlan)
	if err != nil {
		return nil, cliError("failed to unmarshal plan %s: %v", input, err)
	}

	return plan, nil
}

func (app *Application) readPlan(input string) ([]byte, error) {
	if input == "-" {
		return app.readPlanFromInput()
	}
	return app.readPlanFromFile(input)
}

func (app *Application) readPlanFromInput() ([]byte, error) {
	serializedPlan, err := ioutil.ReadAll(app.Input)
	if err != nil {
		return nil, cliError("failed to read plan from input: %v", err)
	}
	return serializedPlan, nil
}

func (app *Application) readPlanFromFile(file string) ([]byte, error) {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return nil, cliError("could not find plan at %s", file)
	}

	serializedPlan, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, cliError("failed to read plan %s: %v", file, err)
	}

	return serializedPlan, err
}
