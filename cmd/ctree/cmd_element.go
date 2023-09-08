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
		"-list", "-save", "-new", "-select", "-" + cmdHelp,
	})
}
func (c CmdElement) Exec(line string) (err error) {

	// Define and parse flags and get arguments
	var new /* , save, list, selct */ bool
	flags := c.NewFlagSet(c.Name(), c.Usage(), c.Help())
	flags.BoolVar(&new, "new", new, "create new tree element")
	// flags.BoolVar(&save, "save", save, "save current tree")
	// flags.BoolVar(&list, "list", list, "print list of trees")
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
		fmt.Printf("element '%s' created\n", e.Value())

	// Wrong flag selected or flags is empty
	default:
		flags.Usage()
		err = ErrNoFlags
		return
	}

	return
}
