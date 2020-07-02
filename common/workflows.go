package common

import (
	"fmt"
	"regexp"

	"github.com/sirupsen/logrus"
)

type FringeWorkflowAssetType []AssetType

type FringeWorkflowNewAssetTrigger struct {
	Types FringeWorkflowAssetType `yaml:"types"`
	Regex StringArray             `yaml:"regex"`
}

type FringeWorkflowTrigger struct {
	NewAsset FringeWorkflowNewAssetTrigger `yaml:"new_asset"`
}

type FringeWorkflowJob struct {
	Name   string `yaml:"name"`
	Module string `yaml:"module"`
}

type FringeWorkflow struct {
	Name string                `yaml:"name"`
	On   FringeWorkflowTrigger `yaml:"on"`
	Jobs []FringeWorkflowJob   `yaml:"jobs"`
}

// Parse a workflow YAML file, check the structure and returns the object
func NewFringeWorkflow(filePath string) (*FringeWorkflow, error) {
	workflow := &FringeWorkflow{}
	err := ParseYAMLFile(filePath, workflow)
	if err != nil {
		return nil, err
	}

	return workflow, nil
}

func (w *FringeWorkflow) String() string {
	return "A workflow"
}

func (w *FringeWorkflow) Run(assets []Asset, sess *Session, config *FringeConfig) error {
	for _, asset := range assets {
		err := w.runAsset(asset, sess, config)
		if err != nil {
			return err
		}
	}

	return nil
}

func (w *FringeWorkflow) runAsset(asset Asset, sess *Session, config *FringeConfig) error {
	// Get the `new_asset` field to trigger the workflow on the `ctx.NewAssets`
	newAsset := &w.On.NewAsset
	if newAsset != nil {
		isCorrectType := false
		for _, eventType := range newAsset.Types {
			if eventType == ASSET_HOSTNAME {
				isCorrectType = isCorrectType || IsHostname(asset.Value)
				break
			} else if eventType == ASSET_IP {
				isCorrectType = isCorrectType || IsIPv4(asset.Value)
				break
			} else if eventType == ASSET_URL {
				isCorrectType = isCorrectType || IsURL(asset.Value)
				break
			} else {
				err := fmt.Errorf("Invalid asset type, must be one of hostname, ip or url.")
				return err
			}
		}

		if !isCorrectType {
			err := fmt.Errorf("The new asset does not have the correct types.")
			return err
		}

		// TODO: add filter on tags

		for _, eventRegex := range newAsset.Regex {
			re, err := regexp.Compile(eventRegex)
			if err != nil {
				logrus.Debug(err)
				return fmt.Errorf("The workflow regex is not valid.")
			}

			if !re.MatchString(asset.Value) {
				return fmt.Errorf("The asset does not match the regex.")
			}
		}

		ctx, err := NewModuleContext(asset, config)
		if err != nil {
			return err
		}

		for _, job := range w.Jobs {
			logrus.Infof("Execute job \"%s\"", job.Name)

			module, err := sess.Module(job.Module)
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
	}

	return nil
}
