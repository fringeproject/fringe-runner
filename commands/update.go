package commands

import (
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/fringeproject/fringe-runner/common"
)

type UpdateCommand struct {
}

func (s *UpdateCommand) Execute(c *cli.Context) error {
	err := common.UpdateModuleRessources()
	if err != nil {
		return err
	}

	logrus.Info("Downloaded the ressources files.")

	return nil
}

func newUpdateCommand() *UpdateCommand {
	return &UpdateCommand{}
}

// Register the command
func init() {
	common.RegisterCommand("update", "Update the runner ressources", newUpdateCommand(), nil)
}
