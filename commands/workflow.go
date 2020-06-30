package commands

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"

	"github.com/fringeproject/fringe-runner/common"
	"github.com/fringeproject/fringe-runner/modules"
)

type WorkflowCommand struct {
	context *cli.Context
	session *common.Session
	config  *common.FringeConfig
}

type StringArray []string

// From: https://github.com/go-yaml/yaml/issues/100
func (a *StringArray) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var multi []string
	err := unmarshal(&multi)
	if err != nil {
		var single string
		err := unmarshal(&single)
		if err != nil {
			return err
		}
		*a = []string{single}
	} else {
		*a = multi
	}
	return nil
}

type WorkflowNewAssetTrigger struct {
	Types StringArray `yaml:"types"`
	Regex StringArray `yaml:"regex"`
}

type WorkflowTrigger struct {
	NewAsset WorkflowNewAssetTrigger `yaml:"new_asset"`
}

type WorkflowJob struct {
	Name   string `yaml:"name"`
	Module string `yaml:"module"`
}

type FringeWorkflow struct {
	Name string          `yaml:"name"`
	On   WorkflowTrigger `yaml:"on"`
	Jobs []WorkflowJob   `yaml:"jobs"`
}

func ParseWorkflowFile(workflowPath string, config *common.FringeConfig) (*FringeWorkflow, error) {
	workflow := &FringeWorkflow{}

	workflowFile, err := ioutil.ReadFile(workflowPath)
	if err != nil {
		logrus.Warning(err)
		err = fmt.Errorf("Cannot read the workflow file: %s.", workflowPath)
		return nil, err
	}

	err = yaml.Unmarshal([]byte(workflowFile), &workflow)
	if err != nil {
		err = fmt.Errorf("The workflow file is not a valid Workflow YAML file.")
		return nil, err
	}

	return workflow, nil
}

func (s *WorkflowCommand) executeWorkflow(rawAsset string, workflow *FringeWorkflow) error {
	logrus.Infof("Executing workflow: %s", workflow.Name)

	newAsset := &workflow.On.NewAsset
	if newAsset != nil {
		logrus.Info("There is a NewAsset event found.")

		isCorrectType := false
		for _, eventType := range newAsset.Types {
			if eventType == "hostname" {
				isCorrectType = isCorrectType || common.IsHostname(rawAsset)
			} else if eventType == "ip" {
				isCorrectType = isCorrectType || common.IsIPv4(rawAsset)
			} else if eventType == "url" {
				isCorrectType = isCorrectType || common.IsURL(rawAsset)
			} else {
				err := fmt.Errorf("Invalid asset type, must be one of hostname, ip or url.")
				return err
			}
		}

		if !isCorrectType {
			err := fmt.Errorf("The new asset does not have the correct types.")
			return err
		}

		for _, eventRegex := range newAsset.Regex {
			re, err := regexp.Compile(eventRegex)
			if err != nil {
				logrus.Debug(err)
				return fmt.Errorf("The regex is not valid.")
			}

			if !re.MatchString(rawAsset) {
				return fmt.Errorf("The asset does not match the regex.")
			}
		}

		asset := common.Asset{
			Value: rawAsset,
			Type:  "",
		}

		// Create a module context for the execution
		ctx, err := common.NewModuleContext(asset, s.config)
		if err != nil {
			logrus.Warn("Cannot create module context.")
			logrus.Debug(err)
			return err
		}

		for _, job := range workflow.Jobs {
			logrus.Infof("Execute job \"%s\"", job.Name)

			module, err := s.session.Module(job.Module)
			if err != nil {
				err := fmt.Errorf("Cannot find the module \"%s\" in the workflow.", job.Module)
				return err
			}

			// Run the module
			err = module.Run(ctx)
			if err != nil {
				logrus.Warn("Module execution return an error.")
				logrus.Debug(err)
				return err
			}
		}

		// Get the new assets and convert to print it as a JSON list
		updateJob := common.FringeClientUpdateJobRequest{
			ID:          "workflow-" + workflow.Name,
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

func (s *WorkflowCommand) Execute(c *cli.Context, config *common.FringeConfig) error {
	filePath := c.String("path")

	if !common.FileExists(filePath) {
		return fmt.Errorf("The file does not exists.")
	}

	if c.NArg() != 1 {
		return fmt.Errorf("Error getting asset for workflow.")
	}

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

	workflow, err := ParseWorkflowFile(filePath, s.config)
	if err != nil {
		return err
	}

	err = s.executeWorkflow(c.Args().Get(0), workflow)
	if err != nil {
		return err
	}

	return nil
}

func newWorkflowCommand() *WorkflowCommand {
	return &WorkflowCommand{}
}

// Register the command
func init() {
	common.RegisterCommand("workflow", "Execute a workflow manually", newWorkflowCommand(), []cli.Flag{
		&cli.StringFlag{
			Name:    "path",
			Aliases: []string{"p"},
			Usage:   "Path to the rule file",
		},
	})
}
