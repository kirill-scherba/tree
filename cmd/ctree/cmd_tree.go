// Copyright 2023 Kirill Scherba <kirill@scherba.ru>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Tree CLI application. Command Tree.

package main

import (
	"fmt"
	"strings"

	"github.com/kirill-scherba/tree"
	"github.com/teonet-go/teonet/cmd/teonet/menu"
)

// CmdTree command structure
type CmdTree struct{ TreeCommand }

// Create CmdTree command
func (cli *Tree) newCmdTree() menu.Item {
	return CmdTree{TreeCommand: TreeCommand{cli}}
}

func (c CmdTree) Name() string  { return cmdTree }
func (c CmdTree) Usage() string { return "[flag] [name||id]" }
func (c CmdTree) Help() string {
	return "" +
		"any operation with tree depending of flag " +
		"(print current tree if flag omitted)"
}
func (c CmdTree) Compliter() (cmpl []menu.Compliter) {
	return c.menu.MakeCompliterFromString([]string{
		"-list", "-save", "-new", "-select", "-" + cmdHelp,
	})
}
func (c CmdTree) Exec(line string) (err error) {

	// Define and parse flags and get arguments
	var new, save, list, selct bool
	flags := c.NewFlagSet(c.Name(), c.Usage(), c.Help())
	flags.BoolVar(&new, "new", new, "create new tree")
	flags.BoolVar(&save, "save", save, "save current tree")
	flags.BoolVar(&list, "list", list, "print list of trees")
	flags.BoolVar(&selct, "select", selct, "select tree from list of trees by id")
	err = flags.Parse(c.menu.SplitSpace(line))
	if err != nil {
		return
	}
	args := flags.Args()
	argc := len(args)

	switch {

	// Print help
	case argc > 0 && args[0] == cmdHelp:
		flags.Usage()
		return

	// Create new tree: -new flag
	case new:
		if argc > 0 {
			name := strings.Join(args, " ")
			c.tree = tree.New[TreeData](name)
		} else {
			c.tree = tree.New[TreeData]()
		}
		c.treeList.add(c.tree)
		fmt.Printf("new tree '%s' created, id: %s\n", c.tree, c.tree.Id())
		return

	// Print list of trees: -list flag
	case list:
		fmt.Printf("%s", c.treeList.String())
		return

	// Select tree from list of trees by id: -select flag
	case selct:
		if argc == 0 {
			flags.Usage()
			err = ErrWrongNumArguments
			return
		}
		tree := c.treeList.get(args[0])
		if tree == nil {
			err = ErrWrongIdArgument
			return
		}
		c.tree = tree
		fmt.Printf("tree '%s' selected, id: %s\n", c.tree, c.tree.Id())
		return

	// Print current tre name and id
	default:
		fmt.Printf("current tree bame: '%s', id: %s\n", c.tree, c.tree.Id())
	}

	// Add alias
	// c.alias.add(args[0], args[1])

	return
}
