// Copyright 2023 Kirill Scherba <kirill@scherba.ru>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Tree CLI application. Command Element.

package main

import (
	"flag"
	"fmt"
	"strconv"
	"strings"

	"github.com/kirill-scherba/tree"
	"github.com/teonet-go/teonet/cmd/teonet/menu"
)

// CmdTree command structure
type CmdElement struct {
	TreeCommand
	flagsSet *flag.FlagSet
	flags    struct {
		new, add, list, path, print, selct bool
	}
	// flagsM flagsMap
	args []string
	argc int
}

// type flagsMap = map[string]flagsMapData
// type flagsMapData struct {
// 	flag bool
// 	f    func() error
// }

// Create CmdElement command
func (cli *Tree) newCmdElement() menu.Item {
	item := &CmdElement{TreeCommand: TreeCommand{cli}}
	// item.flagsM = flagsMap {
	// 	"new": {},
	// }
	return item
}
func (c CmdElement) Name() string  { return cmdElement }
func (c CmdElement) Usage() string { return "[flag] [name][, cost]" }
func (c CmdElement) Help() string {
	return "" +
		"any operation with current tree elements depending of flag " +
		"(print current tree element if flag omitted)"
}
func (c *CmdElement) Parse(line string) (err error) {
	c.flagsSet = c.NewFlagSet(c.Name(), c.Usage(), c.Help())
	c.flagsSet.BoolVar(&c.flags.new, "new", false, "create new tree element")
	c.flagsSet.BoolVar(&c.flags.add, "add", false, "add way to existing or new element from current trees element")
	c.flagsSet.BoolVar(&c.flags.list, "list", false, "list all elements in this tree")
	c.flagsSet.BoolVar(&c.flags.path, "path", false, "prints path from current element to selected in this tree")
	c.flagsSet.BoolVar(&c.flags.print, "print", false, "prints the tree started from current element")
	c.flagsSet.BoolVar(&c.flags.selct, "select", false, "select element in current tree by name")

	err = c.flagsSet.Parse(c.menu.SplitSpace(line))
	if err != nil {
		return
	}
	c.args = c.flagsSet.Args()
	c.argc = len(c.args)

	return
}
func (c CmdElement) Compliter() (cmpl []menu.Compliter) {
	return c.menu.MakeCompliterFromString([]string{
		"-new", "-add", "-list", "-path", "-print", "-select", "-" + cmdHelp,
	})
}
func (c CmdElement) Exec(line string) (err error) {

	c.Parse(line)

	switch {

	// Print help
	case c.argc > 0 && c.args[0] == cmdHelp:
		c.flagsSet.Usage()

	// Create new tree element: -new flag
	case c.flags.new:
		err = c.new()

	// Adds new element to current trees element: -add flag
	case c.flags.add:
		err = c.add()

	// Prints all tree elements: -list flag
	case c.flags.list:
		err = c.list()

	// Prints the tree started from current element: -print flag
	case c.flags.print:
		err = c.print()

	// Prints path from current element to selected in this tree: -path flag
	case c.flags.path:
		err = c.path()

	// Select element in current tree by name: -selct flag
	case c.flags.selct:
		err = c.selct()

	// Print current tree element if flag omitted
	default:
		fmt.Printf("current element: '%s'\n", c.element.Value())
	}

	return
}

// checkArgs checks number of erguments
func (c CmdElement) checkArgs(n int) (err error) {

	if c.argc < n {
		c.flagsSet.Usage()
		err = ErrWrongNumArguments
	}

	return
}

// new creates new tree element: -new flag
func (c CmdElement) new() (err error) {

	if err = c.checkArgs(1); err != nil {
		return
	}

	name := strings.Join(c.args, " ")
	e := c.tree.New(TreeData(name))
	c.element = e
	fmt.Printf("element '%s' created\n", e.Value())

	return
}

// add adds new element to current trees element: -add flag
func (c CmdElement) add() (err error) {

	if err = c.checkArgs(1); err != nil {
		return
	}

	// Parse arguments
	name := strings.Join(c.args, " ")
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

	return
}

// list prints all tree elements: -list flag
func (c CmdElement) list() (err error) {
	fmt.Printf("elements in tree name: '%s', id: %s\n%s\n",
		c.tree, c.tree.Id(), c.element.List().Sort())
	return
}

// print prints the tree started from current element: -print flag
func (c CmdElement) print() (err error) {
	fmt.Printf("elements in tree name: '%s', id: %s\n%s\n",
		c.tree, c.tree.Id(), c.element)
	return
}

// path prints path from current element to selected in this tree: -path flag
func (c CmdElement) path() (err error) {

	if c.argc == 0 {
		c.flagsSet.Usage()
		err = ErrWrongNumArguments
		return
	}
	name := strings.Join(c.args, " ")
	e := c.element.Get(name)
	if e == nil {
		err = ErrElementNotFound
		return
	}
	p := c.element.PathTo(e).Sort()
	fmt.Printf("paths to element in tree name: '%s', id: %s\n%s\n",
		c.tree, c.tree.Id(), p)

	return
}

// selct selects element in current tree by name: -selct flag
func (c CmdElement) selct() (err error) {

	if c.argc == 0 {
		c.flagsSet.Usage()
		err = ErrWrongNumArguments
		return
	}
	name := strings.Join(c.args, " ")
	e := c.element.Get(name)
	if e == nil {
		err = ErrElementNotFound
		return
	}
	c.element = e
	fmt.Printf("element '%s' selected\n", e.Value())

	return
}
