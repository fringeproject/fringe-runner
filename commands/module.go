package commands

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/fringeproject/fringe-runner/common"
	"github.com/fringeproject/fringe-runner/modules"
	"github.com/fringeproject/fringe-runner/session"
)

type ModuleCommand struct {
	context *cli.Context
	session *session.Session
}

func (s *ModuleCommand) listModules() error {
	// Get the list of modules
	moduleList, err := s.session.GetModules()
	if err != nil {
		return err
	}

	// Convert the list to print it as a JSON list
	moduleJSON, err := json.MarshalIndent(moduleList, "", "\t")
	if err != nil {
		return fmt.Errorf("Couldn't format the module list to JSON.")
	}

	fmt.Println(string(moduleJSON))

	return nil
}

func (s *ModuleCommand) executeModule() error {
	if s.context.NArg() != 2 {
		return fmt.Errorf("Error getting arguments for execution <module> <asset>.")
	}

	moduleSlug := s.context.Args().Get(0)
	asset := s.context.Args().Get(1)

	// Get the module from the slug
	module, err := s.session.Module(moduleSlug)
	if err != nil {
		logrus.Warnf("Cannot find module with slug \"%s\"", moduleSlug)
		logrus.Debug(err)
		return err
	}

	// Create a module context for the execution
	ctx, err := common.NewModuleContext(asset)
	if err != nil {
		logrus.Warn("Cannot crate module context.")
		logrus.Debug(err)
		return err
	}

	// Run the module
	err = module.Run(ctx)
	if err != nil {
		logrus.Warn("Module execution return an error.")
		logrus.Debug(err)
		return err
	}

	return nil
}

func (s *ModuleCommand) Execute(c *cli.Context) error {

	// Create a new session that hold the modules
	// Even if we don't know the command the user want to pass, we create the
	// session and load the modules.
	sess, err := session.NewSession()
	if err != nil {
		logrus.Warn(err)
		os.Exit(1)
	}
	defer sess.Close()

	// Add the context and the session to the current command for re-use
	s.context = c
	s.session = sess

	// Load Fringe modules in the session
	modules.LoadModules(s.session)

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
