package common

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var commands []*cli.Command

type Commander interface {
	Execute(c *cli.Context, config *FringeConfig) error
}

func setLogLevel(level string) error {
	logrus.SetFormatter(&logrus.TextFormatter{})
	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		return fmt.Errorf("The config value \"%s\" is not valid.", level)
	}
	logrus.SetLevel(logLevel)

	return nil
}

func RegisterCommand(name, usage string, data Commander, flags []cli.Flag) {
	command := &cli.Command{
		Name:  name,
		Usage: usage,
		Action: func(c *cli.Context) error {
			configPath := c.String("config")

			config, err := ReadConfigFile(configPath)
			if err != nil {
				return err
			}

			err = setLogLevel(config.LogLevel)
			if err != nil {
				return err
			}

			// TODO: remove flag `config` from the context before passing it to
			// the command
			return data.Execute(c, config)
		},
		Flags: flags,
	}

	configFlag := &cli.StringFlag{
		Name:    "config",
		Aliases: []string{"c"},
		Usage:   "load configuration from `FILE`",
		EnvVars: []string{"FRINGE_CONFIG"},
	}

	command.Flags = append(command.Flags, configFlag)

	commands = append(commands, command)
}

func GetCommands() []*cli.Command {
	return commands
}
