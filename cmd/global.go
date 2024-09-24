package cmd

import (
  "sort"

  "io"

  "os"

  "fmt"

  "github.com/cloudogu/nexus-claim/domain"
  "github.com/cloudogu/nexus-claim/infrastructure"
  "gopkg.in/urfave/cli.v2"
)

// Application holds all command actions and gives context to them.
// The main use for a central command holder is testing.
type Application struct {
	Output         io.Writer
	Input          io.Reader
	commands       cli.Commands
	nexusAPIClient domain.NexusAPIClient
}

var (
	application = &Application{
		Output:   os.Stdout,
		Input:    os.Stdin,
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
			Value:  "http://localhost:8081",
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
		cli.IntFlag{
			Name:  "timeout",
			Value: 30,
			Usage: "Timeout for HTTP client requests",
		},
		cli.BoolFlag{
			Name:   "nexus2",
			Usage:  "use this flag to use nexus-claim with nexus 2",
			EnvVar: "NEXUS_V2",
		},
	}
}

func (app *Application) createNexusAPIClient(c *cli.Context) domain.NexusAPIClient {

	if app.nexusAPIClient != nil {

		return app.nexusAPIClient
	} else if c.Bool("nexus2") {
		return infrastructure.NewHTTPNexusAPIClient(
			c.GlobalString("server"),
			c.GlobalString("username"),
			c.GlobalString("password"),
		)
	}

	return infrastructure.NewNexus3APIClient(
		c.GlobalString("server"),
		c.GlobalString("username"),
		c.GlobalString("password"),
		c.GlobalInt("timeout"),
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
	_, err := app.Output.Write([]byte(action.ToString() + "\n"))
	if err != nil {
		return fmt.Errorf("failed to write action: %w", err)
	}
	return err
}

func cliError(format string, args ...interface{}) error {
	return cli.NewExitError(fmt.Sprintf(format, args...), 1)
}
