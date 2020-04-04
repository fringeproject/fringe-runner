package main

import (
	"os"
	"path"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	// Import commands
	_ "github.com/fringeproject/fringe-runner/commands"
	"github.com/fringeproject/fringe-runner/common"
)

func main() {
	// Read .env file and set environnement variables
	errGotEnv := godotenv.Load()

	// Logger configuration
	// TODO: add custom logger configuration
	logrus.SetFormatter(&logrus.TextFormatter{})
	logLevel, err := logrus.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		logLevel = logrus.InfoLevel
	}
	logrus.SetLevel(logLevel)

	// Print the godotenv error after we set the loglevel
	if errGotEnv != nil {
		logrus.Debug("Coulnd not load .env file")
	}

	app := cli.NewApp()
	app.Name = path.Base(os.Args[0])
	app.Usage = "a Fringe Runner"
	app.Version = common.AppVersion.ShortLine()
	cli.VersionPrinter = common.AppVersion.Printer
	app.Authors = []*cli.Author{
		&cli.Author{
			Name:  "Fringe Project",
			Email: "contact@fringeproject.com",
		},
	}
	app.Commands = common.GetCommands()
	app.CommandNotFound = func(context *cli.Context, command string) {
		logrus.Errorln("Command", command, "not found.")
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Errorln(err)
	}
}
