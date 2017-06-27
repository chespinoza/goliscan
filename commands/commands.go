package commands

import (
	"github.com/urfave/cli"

	"gitlab.com/ayufan/golang-cli-helpers"
)

var commands []cli.Command

type Command interface {
	Execute(c *cli.Context)
}

func RegisterSimpleCommand(command cli.Command) {
	commands = append(commands, command)
}

func RegisterCommand(name, usage string, data Command, flags ...cli.Flag) {
	RegisterSimpleCommand(cli.Command{
		Name:   name,
		Usage:  usage,
		Action: data.Execute,
		Flags:  append(flags, clihelpers.GetFlagsFromStruct(data)...),
	})
}

func GetCommands() []cli.Command {
	return commands
}
