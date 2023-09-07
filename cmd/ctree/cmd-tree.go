// Copyright 2023 Kirill Scherba <kirill@scherba.ru>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Tree CLI application. Command Tree module.

package main

import (
	"fmt"
	"strings"

	"github.com/kirill-scherba/tree"
	"github.com/teonet-go/teonet/cmd/teonet/menu"
)

// Create CmdTree commands
func (cli *Treecli) newCmdTree() menu.Item {
	return CmdTree{TreecliCommand: TreecliCommand{cli}}
}

// CmdTree connect to peer command ----------------------------------------
type CmdTree struct {
	TreecliCommand
}

func (c CmdTree) Name() string  { return cmdTree }
func (c CmdTree) Usage() string { return "[flag] <name>" }
func (c CmdTree) Help() string  { return "select or create tree ('choose' flag runs by default)" }
func (c CmdTree) Exec(line string) (err error) {
	var list, save, new, choose bool
	flags := c.NewFlagSet(c.Name(), c.Usage(), c.Help())
	flags.BoolVar(&new, "new", list, "create new tree")
	flags.BoolVar(&save, "save", save, "save current tree")
	flags.BoolVar(&list, "list", list, "show list of trees")
	flags.BoolVar(&choose, "choose", list, "choose tree to use by id")
	err = flags.Parse(c.menu.SplitSpace(line))
	if err != nil {
		return
	}
	args := flags.Args()
	argc := len(args)

	switch {
	// Check help
	case argc > 0 && args[0] == cmdHelp:
		flags.Usage()
		return

	// Check length of arguments
	// case argc == 0 && new:
	// 	flags.Usage()
	// 	err = ErrWrongNumArguments
	// 	return

	// Check -new flag
	case new:
		if argc > 0 {
			name := strings.Join(args, " ")
			c.tree = tree.New[TreeData](name)
		} else {
			c.tree = tree.New[TreeData]()
		}
		c.treeList.add(c.tree)
		fmt.Printf("new tree `%s` created, id: %s\n", c.tree, c.tree.Id())
		return

	// Check -list flag
	case list:
		fmt.Printf("%s", c.treeList.String())
		return

	// Check -list flag
	// case list:
	// 	aliases := c.alias.list()
	// 	for i := range aliases {
	// 		fmt.Printf("%s\n", aliases[i])
	// 	}
	// 	return

	// Check -save flag
	// case save:
	// 	aliases := c.alias.list()
	// 	c.batch.Save(aliasBatchFile, CmdTree, aliases)
	// 	return

	default:
		fmt.Printf("current tree `%s`, id: %s\n", c.tree, c.tree.Id())
	}

	// Add alias
	// c.alias.add(args[0], args[1])

	return
}
func (c CmdTree) Compliter() (cmpl []menu.Compliter) {
	return c.menu.MakeCompliterFromString([]string{"-list", "-save"})
}
