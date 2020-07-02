package commands

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/fringeproject/fringe-runner/common"
	"github.com/fringeproject/fringe-runner/modules"
)

type ModuleCommand struct {
	context *cli.Context
	session *common.Session
	config  *common.FringeConfig
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
	// Check if module and asset argument are set
	if s.context.String("module") == "" {
		return fmt.Errorf("Please specify a module to execute.")
	}

	if s.context.String("asset") == "" {
		return fmt.Errorf("Please specify an asset or a file containing assets.")
	}

	// Get the module from the slug
	module, err := s.session.Module(s.context.String("module"))
	if err != nil {
		return err
	}

	workflow := &common.FringeWorkflow{}
	if len(s.context.String("workflow")) > 0 {
		// Parse the workflow file
		workflow, err = common.NewFringeWorkflow(s.context.String("workflow"))
		if err != nil {
			logrus.Debug(err)
			return fmt.Errorf("Cannot parse workflow file.")
		}
	}

	// The asset argument can be a raw asset or a path to a file
	assets := []common.Asset{}
	if common.FileExists(s.context.String("asset")) {
		file, err := os.Open(s.context.String("asset"))
		if err != nil {
			return err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			asset := common.Asset{
				Value: scanner.Text(),
				Type:  "",
			}
			assets = append(assets, asset)
		}

		if err := scanner.Err(); err != nil {
			return err
		}
	} else {
		asset := common.Asset{
			Value: s.context.String("asset"),
			Type:  "",
		}
		assets = append(assets, asset)
	}

	for _, asset := range assets {
		// Create a module context for the execution
		ctx, err := common.NewModuleContext(asset, s.config)
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
			continue
		}

		err = workflow.Run(ctx.NewAssets, s.session, s.config)
		if err != nil {
			logrus.Infof("Workflow returns an error: %s", err)
		}

		// Get the new assets and convert to print it as a JSON list
		updateJob := common.FringeClientUpdateJobRequest{
			ID:          "",
			Status:      common.JOB_STATUS_SUCCESS,
			Assets:      ctx.NewAssets,
			Tags:        ctx.NewTags,
			Description: "",
			StartedAt:   0,
			EndedAt:     0,
		}

		fmt.Println(updateJob.JSON())
	}

	return nil
}

func (s *ModuleCommand) Execute(c *cli.Context, config *common.FringeConfig) error {
	// Create a new session that hold the modules
	// Even if we don't know the command the user want to pass, we create the
	// session and load the modules.
	sess, err := common.NewSession()
	if err != nil {
		logrus.Warn(err)
		os.Exit(1)
	}
	defer sess.Close()

	// Add the context, session and config to the current command for re-use
	s.context = c
	s.session = sess
	s.config = config

	// Load Fringe modules in the session
	modules.LoadModules(s.session)

	// Check command args to know what to do
	if c.Bool("list-modules") {
		return s.listModules()
	}

	// By default we execute a module
	return s.executeModule()
}

func newModuleCommand() *ModuleCommand {
	return &ModuleCommand{}
}

// Register the command
func init() {
	common.RegisterCommand("module", "Uses Fringe modules manually", newModuleCommand(), []cli.Flag{
		&cli.BoolFlag{
			Name:    "list-modules",
			Aliases: []string{"L"},
			Usage:   "List available modules",
		},
		&cli.StringFlag{
			Name:    "asset",
			Aliases: []string{"a"},
			Usage:   "Asset or file containing assets",
		},
		&cli.StringFlag{
			Name:    "module",
			Aliases: []string{"m"},
			Usage:   "Module slug to execute",
		},
		&cli.StringFlag{
			Name:    "workflow",
			Aliases: []string{"w"},
			Usage:   "Workflow file to use",
		},
	})
}
