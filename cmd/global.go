package cmd

import (
	"sort"

	"io"

	"os"

	"fmt"

	"github.com/cloudogu/nexus-claim/domain"
	"github.com/cloudogu/nexus-claim/infrastructure"
	"github.com/pkg/errors"
	"gopkg.in/urfave/cli.v2"
)

// Application holds all command actions and gives context to them.
// The main use for a central command holder is testing.
type Application struct {
	Output         io.Writer
	commands       cli.Commands
	nexusAPIClient domain.NexusAPIClient
}

var (
	application = &Application{
		Output:   os.Stdout,
		commands: []cli.Command{},
	}
)

// GetApplication returns the holder for commands and global flags
func GetApplication() *Application {
	return application
}

func (app *Application) registerCommand(cmd cli.Command) {
	app.commands = append(app.commands, cmd)
}

// Commands returns all registered commands ordered by name
func (app *Application) Commands() []cli.Command {
	commands := app.commands
	sort.Sort(commandsByName(commands))
	return commands
}

type commandsByName []cli.Command

func (c commandsByName) Len() int {
	return len(c)
}

func (c commandsByName) Less(i, j int) bool {
	return c[i].Name < c[j].Name
}

func (c commandsByName) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

// GlobalFlags returns all global flags, which are required by the commands
func (app *Application) GlobalFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:   "server, s",
			Value:  "http://localhost:8081/nexus",
			Usage:  "Url to the nexus server",
			EnvVar: "NEXUS_SERVER",
		},
		cli.StringFlag{
			Name:   "username, u",
			Value:  "admin",
			Usage:  "Username of a nexus admin user",
			EnvVar: "NEXUS_USER",
		},
		cli.StringFlag{
			Name:   "password, p",
			Value:  "admin123",
			Usage:  "Password of the nexus user",
			EnvVar: "NEXUS_PASSWORD",
		},
	}
}

func (app *Application) createNexusAPIClient(c *cli.Context) domain.NexusAPIClient {
	if app.nexusAPIClient != nil {
		return app.nexusAPIClient
	}
	return infrastructure.NewHTTPNexusAPIClient(
		c.GlobalString("server"),
		c.GlobalString("username"),
		c.GlobalString("password"),
	)
}

func (app *Application) printPlan(plan *domain.Plan) error {
	for _, action := range plan.GetActions() {
		err := app.printAction(action)
		if err != nil {
			return err
		}
	}
	return nil
}

func (app *Application) printAction(action domain.Action) error {
	value := fmt.Sprintf("%s %s\n", createOperatorFromActionType(action.GetType()), action.GetRepository().ID)
	_, err := app.Output.Write([]byte(value))
	if err != nil {
		return errors.Wrap(err, "failed to write action to output")
	}
	return nil
}

func createOperatorFromActionType(actionType domain.ActionType) string {
	switch actionType {
	case domain.ActionCreate:
		return "+"
	case domain.ActionModify:
		return "~"
	case domain.ActionRemove:
		return "-"
	default:
		return "#"
	}
}
