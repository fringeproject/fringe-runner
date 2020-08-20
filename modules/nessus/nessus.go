package nessus

import (
	"fmt"
	"github.com/sirupsen/logrus"

	"github.com/fringeproject/fringe-runner/common"
)

type Nessus struct {
}

func NewNessus() *Nessus {
	return &Nessus{}
}

func (m *Nessus) Name() string {
	return "Nessus"
}

func (m *Nessus) Slug() string {
	return "nessus"
}

func (m *Nessus) Description() string {
	return "Connect to a nessus server to start and retreive scans results."
}

func (m *Nessus) ResourceURLs() []common.ModuleResource {
	return nil
}

func (m *Nessus) Run(ctx *common.ModuleContext) error {
	target, err := ctx.GetAssetAsIP()
	if err != nil {
		return err
	}

	endpoint, err := ctx.GetConfigurationValue("nessus_endpoint")
	if err != nil {
		err := fmt.Errorf("You must provide a nessus_endpoint value.")
		return err
	}

	username, err := ctx.GetConfigurationValue("nessus_username")
	if err != nil {
		err := fmt.Errorf("You must provide a nessus_username value.")
		return err
	}

	password, err := ctx.GetConfigurationValue("nessus_password")
	if err != nil {
		err := fmt.Errorf("You must provide a nessus_password value.")
		return err
	}

	api, err := NewNessusAPI(ctx, endpoint, username, password)
	if err != nil {
		return err
	}

	templates, err := api.ListScanTemplates()
	if err != nil {
		logrus.Debug(err)
		err = fmt.Errorf("Cannot list Nessus templates.")
		return err
	}

	templateName := "advanced_dynamic"
	if templateUID, ok := templates[templateName]; ok {
		scanName := fmt.Sprintf("Fringe job: %s", target)
		scanUUID, err := api.CreateScan(templateUID, target, scanName)
		if err != nil {
			logrus.Debug(err)
			err = fmt.Errorf("Cannot create a new Nessus scan.")
			logrus.Warn(err)
			return err
		}

		logrus.Infof("Created new Nessus scan with ID: %s", scanUUID)
	} else {
		err = fmt.Errorf("Cannot find the template \"%s\" on Nessus.", templateName)
		logrus.Debugf("Available templates: %s", templates)
		return err
	}

	return nil
}
