package cmd

import (
	"encoding/json"
	"io/ioutil"

	"github.com/cloudogu/nexus-claim/domain"
	"github.com/cloudogu/nexus-claim/infrastructure"
	"github.com/pkg/errors"
	"gopkg.in/urfave/cli.v2"
)

// init registers the subcommand plan to the application
func init() {
	application.registerCommand(createPlanCommand(application.plan))
}

func createPlanCommand(actionFunc cli.ActionFunc) cli.Command {
	return cli.Command{
		Name:    "plan",
		Aliases: []string{"p"},
		Usage:   "Reads model and creates a plan for synchronization",
		Action:  actionFunc,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "input, i",
				Usage: "Write plan to `OUTPUT`",
				Value: "nexus-claim.hcl",
			},
			cli.BoolFlag{
				Name:  "quiet, q",
				Usage: "Do no print the plan to stdout",
			},
			cli.StringFlag{
				Name:  "output, o",
				Usage: "Write plan to `OUTPUT`",
			},
		},
	}
}

func (app *Application) plan(c *cli.Context) error {
	dao := createFileModelDAO(c)
	client := app.createNexusAPIClient(c)

	plan, err := domain.CreatePlan(dao, client)
	if err != nil {
		return err
	}

	if !c.Bool("quiet") {
		err = app.printPlan(plan)
		if err != nil {
			return err
		}
	}

	output := c.String("output")
	if output != "" {
		err = writePlan(output, plan)
		if err != nil {
			return err
		}
	}
	return nil
}

func createFileModelDAO(c *cli.Context) domain.ModelDAO {
	return infrastructure.NewFileModelDAO(c.String("input"))
}

func writePlan(output string, plan *domain.Plan) error {
	bytes, err := json.Marshal(plan)
	if err != nil {
		return errors.Wrap(err, "failed to marshal plan")
	}

	err = ioutil.WriteFile(output, bytes, 0644)
	if err != nil {
		return errors.Wrapf(err, "failed to write plan to %s", output)
	}

	return nil
}
