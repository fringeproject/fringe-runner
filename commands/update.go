package commands

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/fringeproject/fringe-runner/common"
	"github.com/fringeproject/fringe-runner/modules"
	"github.com/fringeproject/fringe-runner/session"
)

type UpdateCommand struct {
}

func (s *UpdateCommand) Execute(c *cli.Context, config *common.FringeConfig) error {
	sess, err := session.NewSession()
	if err != nil {
		logrus.Warn(err)
		os.Exit(1)
	}
	defer sess.Close()

	// Load Fringe modules in the session
	modules.LoadModules(sess)

	for _, module := range sess.Modules {
		resources := module.ResourceURLs()
		if resources != nil {
			logrus.Infof("The module %s needs %d resources.", module.Slug(), len(resources))

			for _, resource := range resources {
				logrus.Infof("Try to download %s from %s.", resource.Name, resource.URL)
				err = common.DownloadResource(resource, config)

				if err != nil {
					logrus.Warningf("There was an error while downloading resource %s", resource.Name)
					logrus.Debug(err)
				} else {
					logrus.Infof("Download of %s is a success.", resource.Name)
				}
			}
		}
	}

	return nil
}

func newUpdateCommand() *UpdateCommand {
	return &UpdateCommand{}
}

// Register the command
func init() {
	common.RegisterCommand("update", "Update the runner ressources", newUpdateCommand(), nil)
}
