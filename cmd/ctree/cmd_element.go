// Copyright 2023 Kirill Scherba <kirill@scherba.ru>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Tree CLI application. Command Element.

package main

import (
	"fmt"
	"strings"

	"github.com/teonet-go/teonet/cmd/teonet/menu"
)

// CmdTree command structure
type CmdElement struct{ TreeCommand }

// Create CmdElement command
func (cli *Tree) newCmdElement() menu.Item {
	return CmdElement{TreeCommand: TreeCommand{cli}}
}

func (c CmdElement) Name() string  { return cmdElement }
func (c CmdElement) Usage() string { return "<flag> [name||id]" }
func (c CmdElement) Help() string {
	return "" +
		"any operation with current tree elements depending of flag " +
		"(print current tree element if flag omitted)"
}
func (c CmdElement) Compliter() (cmpl []menu.Compliter) {
	return c.menu.MakeCompliterFromString([]string{
		"-new", "-add", "-list", "-print", "-select", "-" + cmdHelp,
	})
}
func (c CmdElement) Exec(line string) (err error) {

	// Define and parse flags and get arguments
	var new, add, list, print, selct /* , save */ bool
	flags := c.NewFlagSet(c.Name(), c.Usage(), c.Help())
	flags.BoolVar(&new, "new", new, "create new tree element")
	flags.BoolVar(&add, "add", add, "add new element to current trees element")
	flags.BoolVar(&list, "list", print, "prints list of element in this tree")
	flags.BoolVar(&print, "print", print, "prints the tree started from current element")
	flags.BoolVar(&selct, "select", selct, "select element in current tree by name")
	// flags.BoolVar(&save, "save", save, "save current tree")
	// flags.BoolVar(&selct, "select", list, "select tree from list of trees by id")
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

	// Create new tree element: -new flag
	case new:
		if argc == 0 {
			flags.Usage()
			err = ErrWrongNumArguments
			return
		}

		name := strings.Join(args, " ")
		e := c.tree.New(TreeData(name))
		c.element = e
		fmt.Printf("element '%s' created\n", e.Value())

	// Adds new element to current trees element: -add flag
	case add:
		if argc == 0 {
			flags.Usage()
			err = ErrWrongNumArguments
			return
		}
		name := strings.Join(args, " ")
		e := c.tree.New(TreeData(name))
		c.element.Add(e)
		fmt.Printf("element '%s' created and added to %s\n",
			e.Value(), c.element.Value())

	// Prints the tree started from current element: -print flag
	case print:
		fmt.Printf("list elements in tree name: '%s', id: %s\n%s\n",
			c.tree, c.tree.Id(), c.element)

	// Select element in current tree by name: -selct flag
	case selct:
		if argc == 0 {
			flags.Usage()
			err = ErrWrongNumArguments
			return
		}
		name := strings.Join(args, " ")
		e := c.element.Get(name)
		if e == nil {
			err = ErrElementNotFound
			return
		}
		c.element = e
		fmt.Printf("element '%s' selected\n", e.Value())

	// Wrong flag selected or flags is empty
	default:
		flags.Usage()
		err = ErrNoFlags
		return
	}

	return
}
