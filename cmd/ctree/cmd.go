// Copyright 2023 Kirill Scherba <kirill@scherba.ru>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Tree CLI application. Command module.

package main

import (
	"errors"
)

// Commands name
const (
	cmdHelp    = "help"
	cmdTree    = "tree"
	cmdElement = "element"
)

// Names of batch files
const (
	defaultTreeBatchFile = "def_tree.conf"
)

var (
	ErrWrongNumArguments = errors.New("wrong number of arguments")
	ErrWrongIdArgument   = errors.New("wrong id in arguments")
	ErrNoFlags           = errors.New("no flags, one of flag should be specified")
	ErrElementNotFound   = errors.New("element not found")
)

// TreeCommand common Tree CLI command structure
type TreeCommand struct{ *Tree }

// addCommands add commands
func (cli *Tree) addCommands() {
	cli.commands = append(cli.commands,
		cli.newCmdTree(),
		cli.newCmdElement(),
	)
}
