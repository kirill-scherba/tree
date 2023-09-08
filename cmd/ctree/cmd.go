// Copyright 2023 Kirill Scherba <kirill@scherba.ru>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Tree CLI application. Command module.

package main

import (
	"errors"
)

// Command name
const (
	// cmdConnectTo = "connectto"
	// cmdSendTo    = "sendto"
	// cmdAlias     = "alias"
	// cmdStat      = "stat"
	// cmdAPI       = "api"
	cmdHelp = "help"

	cmdTree    = "tree"
	cmdElement = "element"
)

var ErrWrongNumArguments = errors.New("wrong number of arguments")
var ErrWrongIdArgument = errors.New("wrong id in arguments")
var ErrNoFlags = errors.New("no flags, one of flag should be specified")

// TreeCommand common Tree CLI command structure
type TreeCommand struct{ *Tree }

// addCommands add commands
func (cli *Tree) addCommands() {
	cli.commands = append(cli.commands,
		// cli.newCmdAlias(),
		// cli.newCmdConnectTo(),
		// cli.newCmdSendTo(),
		// cli.newCmdStat(),
		// cli.newCmdAPI(),

		cli.newCmdTree(),
		cli.newCmdElement(),
	)
}
