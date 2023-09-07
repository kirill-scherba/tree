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

	cmdTree = "tree"
)

var ErrWrongNumArguments = errors.New("wrong number of arguments")

// TreecliCommand common Treecli command structure
type TreecliCommand struct{ *Treecli }

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
