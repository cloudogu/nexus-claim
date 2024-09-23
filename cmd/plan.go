package cmd

import (
  "fmt"
  "io/ioutil"

  "github.com/cloudogu/nexus-claim/domain"
  "github.com/cloudogu/nexus-claim/infrastructure"
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
				Usage: "Write plan to `OUTPUT`, use '-' to write to stdout",
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

	output := c.String("output")
	if output != "" {
		err = app.writePlan(output, plan)
		if err != nil {
			return cliError("failed to write plan: %v", err)
		}
	}

	if !c.Bool("quiet") && output != "-" {
		err = app.printPlan(plan)
		if err != nil {
			return cliError("failed to print plan: %v", err)
		}
	}

	return nil
}

func createFileModelDAO(c *cli.Context) domain.ModelDAO {
	return infrastructure.NewFileModelDAO(c.String("input"))
}

func (app *Application) writePlan(output string, plan *domain.Plan) error {
	bytes, err := infrastructure.SerializePlan(plan)
	if err != nil {
		return fmt.Errorf("failed to marshal plan: %w", err)
	}

	if output == "-" {
		return app.writePlanToOutput(bytes)
	}

	return app.writePlanToFile(output, bytes)
}

func (app *Application) writePlanToOutput(plan []byte) error {
	_, err := app.Output.Write(plan)
	if err != nil {
		return fmt.Errorf("failed to write plan to output: %w", err)
	}
	return nil
}

func (app *Application) writePlanToFile(file string, plan []byte) error {
	err := ioutil.WriteFile(file, plan, 0644)
	if err != nil {
		return fmt.Errorf("failed to write plan to %s: %w", file, err)
	}
	return nil
}
