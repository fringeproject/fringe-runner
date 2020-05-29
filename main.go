package main

import (
	"os"
	"path"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	// Import commands
	_ "github.com/fringeproject/fringe-runner/commands"
	"github.com/fringeproject/fringe-runner/common"
)

func main() {
	app := &cli.App{
		Name:    path.Base(os.Args[0]),
		Usage:   "a Fringe Runner",
		Version: common.AppVersion.ShortLine(),
		Authors: []*cli.Author{
			&cli.Author{
				Name:  "Fringe Project",
				Email: "contact@fringeproject.com",
			},
		},
		Commands: common.GetCommands(),
		CommandNotFound: func(context *cli.Context, command string) {
			logrus.Errorf("Command % snot found.", command)
		},
	}
	cli.VersionPrinter = common.AppVersion.Printer

	if err := app.Run(os.Args); err != nil {
		logrus.Error(err)
	}
}
