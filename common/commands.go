package common

import (
	"github.com/urfave/cli/v2"
)

var commands []*cli.Command

type Commander interface {
	Execute(c *cli.Context) error
}

func RegisterCommand(name, usage string, data Commander, flags []cli.Flag) {
	command := &cli.Command{
		Name:   name,
		Usage:  usage,
		Action: data.Execute,
		Flags:  flags,
	}

	commands = append(commands, command)
}

func GetCommands() []*cli.Command {
	return commands
}
