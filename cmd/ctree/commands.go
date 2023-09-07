package main

import (
	"errors"
)

// Command name
const (
	cmdConnectTo = "connectto"
	cmdSendTo    = "sendto"
	cmdAlias     = "alias"
	cmdStat      = "stat"
	cmdAPI       = "api"
	cmdHelp      = "help"

	cmdTree = "tree"
)

var ErrWrongNumArguments = errors.New("wrong number of arguments")

// addCommands add commands
func (cli *Treecli) addCommands() {
	cli.commands = append(cli.commands,
		// cli.newCmdAlias(),
		// cli.newCmdConnectTo(),
		// cli.newCmdSendTo(),
		// cli.newCmdStat(),
		// cli.newCmdAPI(),

		cli.newCmdTree(),
	)
}
