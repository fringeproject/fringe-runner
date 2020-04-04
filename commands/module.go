package commands

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/fringeproject/fringe-runner/common"
)

type ModuleCommand struct {
	context *cli.Context
}

func (s *ModuleCommand) listModules() error {
	// TODO: list modules
	return nil
}

func (s *ModuleCommand) executeModule() error {
	if s.context.NArg() != 2 {
		return fmt.Errorf("Error getting arguments for execution <module> <asset>.")
	}

	moduleSlug := s.context.Args().Get(0)
	asset := s.context.Args().Get(1)

	logrus.Infof("Executing module \"%s\" with asset \"%s\".", moduleSlug, asset)

	return nil
}

func (s *ModuleCommand) Execute(c *cli.Context) error {
	s.context = c

	// Check command args to know what to do
	if c.Bool("list") {
		return s.listModules()
	} else if c.Bool("exec") {
		return s.executeModule()
	} else {
		// There is no flag set, then we print the help menu
		return cli.ShowSubcommandHelp(c)
	}
}

func newModuleCommand() *ModuleCommand {
	return &ModuleCommand{}
}

// Register the command
func init() {
	common.RegisterCommand("module", "Use fringe modules", newModuleCommand(), []cli.Flag{
		&cli.BoolFlag{
			Name:    "list",
			Aliases: []string{"l"},
			Usage:   "List available modules then exit",
		},
		&cli.BoolFlag{
			Name:    "exec",
			Aliases: []string{"e"},
			Usage:   "Execute a module",
		},
	})
}
