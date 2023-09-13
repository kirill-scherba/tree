// Copyright 2023 Kirill Scherba <kirill@scherba.ru>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Tree CLI application. Command Element.

package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/kirill-scherba/tree"
	"github.com/teonet-go/teonet/cmd/teonet/menu"
)

// CmdTree command structure
type CmdElement struct{ TreeCommand }

// Create CmdElement command
func (cli *Tree) newCmdElement() menu.Item {
	return CmdElement{TreeCommand: TreeCommand{cli}}
}

func (c CmdElement) Name() string  { return cmdElement }
func (c CmdElement) Usage() string { return "[flag] [name][, cost]" }
func (c CmdElement) Help() string {
	return "" +
		"any operation with current tree elements depending of flag " +
		"(print current tree element if flag omitted)"
}
func (c CmdElement) Compliter() (cmpl []menu.Compliter) {
	return c.menu.MakeCompliterFromString([]string{
		"-new", "-add", "-list", "-path", "-print", "-select", "-" + cmdHelp,
	})
}
func (c CmdElement) Exec(line string) (err error) {

	// Define and parse flags and get arguments
	var new, add, list, path, print, selct /* , save */ bool
	flags := c.NewFlagSet(c.Name(), c.Usage(), c.Help())
	flags.BoolVar(&new, "new", new, "create new tree element")
	flags.BoolVar(&add, "add", add, "add way to existing or new element from current trees element")
	flags.BoolVar(&list, "list", list, "list all elements in this tree")
	flags.BoolVar(&path, "path", path, "prints path from current element to selected in this tree")
	flags.BoolVar(&print, "print", print, "prints the tree started from current element")
	flags.BoolVar(&selct, "select", selct, "select element in current tree by name")
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

		// Parse arguments
		name := strings.Join(args, " ")
		opt := tree.WayOptions{Cost: 1.0}
		if par := strings.Split(name, ","); len(par) > 1 {
			name = strings.TrimSpace(par[0])
			if s, err := strconv.ParseFloat(strings.TrimSpace(par[1]), 64); err == nil {
				opt.Cost = s
			}
		}

		// Get element by name and create new if not exists
		e := c.element.Get(name)
		if e == nil {
			e = c.tree.New(TreeData(name))
		}

		// Add way to element
		c.element.Add(e, opt)
		fmt.Printf("element '%s' created and added to %s\n",
			e.Value(), c.element.Value())

	// Prints all tree elements: -list flag
	case list:
		fmt.Printf("elements in tree name: '%s', id: %s\n%s\n",
			c.tree, c.tree.Id(), c.element.List().Sort())

	// Prints the tree started from current element: -print flag
	case print:
		fmt.Printf("elements in tree name: '%s', id: %s\n%s\n",
			c.tree, c.tree.Id(), c.element)

	// Prints path from current element to selected in this tree: -path flag
	case path:
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
		p := c.element.PathTo(e).Sort()
		fmt.Printf("paths to element in tree name: '%s', id: %s\n%s\n",
			c.tree, c.tree.Id(), p)

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

	// Print current tree element if flag omitted
	default:
		fmt.Printf("current element: '%s'\n", c.element.Value())
	}

	return
}
