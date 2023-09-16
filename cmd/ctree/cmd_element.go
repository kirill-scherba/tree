// Copyright 2023 Kirill Scherba <kirill@scherba.ru>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Tree CLI application. Command Element.

package main

import (
	"flag"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/kirill-scherba/tree"
	"github.com/teonet-go/teonet/cmd/teonet/menu"
)

// CmdTree command structure
type CmdElement struct {
	TreeCommand
	flagsSet *flag.FlagSet
	flagsM   flagsMap
	args     []string
	argc     int
}
type flagsMap = map[string]*flagsMapData
type flagsMapData struct {
	flag  bool
	usage string
	f     func() error
}

// Create CmdElement command
func (cli *Tree) newCmdElement() menu.Item {
	item := &CmdElement{TreeCommand: TreeCommand{cli}}
	item.Flags()
	return item
}
func (c CmdElement) Name() string  { return cmdElement }
func (c CmdElement) Usage() string { return "[flag] [name][, cost]" }
func (c CmdElement) Help() string {
	return "" +
		"any operation with current tree elements depending of flag " +
		"(print current tree element if flag omitted)"
}
func (c *CmdElement) Flags() (err error) {
	c.flagsM = flagsMap{
		"new":    {usage: "create new tree element", f: c.new},
		"add":    {usage: "add way to existing or new element from current trees element", f: c.add},
		"list":   {usage: "list all elements in this tree", f: c.list},
		"path":   {usage: "prints path from current element to selected in this tree", f: c.path},
		"ways":   {usage: "prints current element and his children ways", f: c.ways},
		"remove": {usage: "remove current element", f: c.remove},
		"del":    {usage: "delete way from current element to selected elements splitted by comma", f: c.del},
		"print":  {usage: "prints the tree started from current element", f: c.print},
		"select": {usage: "select element in current tree by name", f: c.selectFlag},
	}
	return
}
func (c *CmdElement) Parse(line string) (err error) {
	c.flagsSet = c.NewFlagSet(c.Name(), c.Usage(), c.Help())
	for f, d := range c.flagsM {
		c.flagsSet.BoolVar(&d.flag, f, false, d.usage)
	}
	err = c.flagsSet.Parse(c.menu.SplitSpace(line))
	if err != nil {
		return
	}
	c.args = c.flagsSet.Args()
	c.argc = len(c.args)

	return
}
func (c CmdElement) Compliter() (cmpl []menu.Compliter) {
	var str []string
	for flag := range c.flagsM {
		str = append(str, "-"+flag)
	}
	sort.Slice(str, func(i, j int) bool {
		if str[i] < str[j] {
			return true
		}
		return false
	})
	str = append(str, "-"+cmdHelp)
	return c.menu.MakeCompliterFromString(str)
}
func (c *CmdElement) Exec(line string) (err error) {

	// Parse arguments and flags
	c.Parse(line)

	// Print help
	if c.argc > 0 && c.args[0] == cmdHelp {
		c.flagsSet.Usage()
		return
	}

	// Find flag sets to true in flags map and execute its function
	for _, d := range c.flagsM {
		if d.flag {
			if d.f != nil {
				return d.f()
			}
			return ErrFunctionNotDefined
		}
	}

	// Default action when flag no sets
	fmt.Printf("current element: '%s'\n", c.element.Value())

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
func (c *CmdElement) new() (err error) {

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
func (c *CmdElement) add() (err error) {

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
		if len(par) > 2 {
			parOneway := strings.TrimSpace(par[2])
			if parOneway == "true" || parOneway == "oneway" {
				opt.Oneway = true
			}
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
func (c *CmdElement) list() (err error) {
	fmt.Printf("elements in tree name: '%s', id: %s\n%s\n",
		c.tree, c.tree.Id(), c.element.List().Sort())
	return
}

// ways prints current element and his children ways: -ways flag
func (c *CmdElement) ways() (err error) {
	// Print parents ways
	fmt.Printf("%s ways:\n", c.element.Value())
	for child := range c.element.Ways() {
		cost, _ := c.element.Cost(child)
		fmt.Printf("  way to '%s' cost: %.2f, way allowed: %v\n",
			child.Value(), cost, c.element.WayAllowed(child))
	}
	return
}

// remove removes current element: -remove flag
func (c *CmdElement) remove() (err error) {
	var child *tree.Element[TreeData]
	for child = range c.element.Ways() {
		break
	}
	_, err = c.element.Remove()
	if child != nil {
		c.element = child
	}
	fmt.Printf("element '%s' removed\n", c.element.Value())
	return
}

// del deletes way from current element to selected: -del flag
func (c *CmdElement) del() (err error) {

	if err = c.checkArgs(1); err != nil {
		return
	}

	name := strings.Join(c.args, " ")
	names := strings.Split(name, ",")
	for i := range names {
		name = strings.TrimSpace(names[i])
		e := c.element.Get(name)
		if e == nil {
			err = ErrElementNotFound
			fmt.Printf("error: '%s' %s\n", name, ErrElementNotFound)
			continue
		}
		_, err = c.element.Del(e)
		if err != nil {
			fmt.Printf("error: '%s' %s\n", name, ErrElementNotFound)
			continue
		}
		fmt.Printf("way to element '%s' deleted\n", name)
	}

	return nil
}

// print prints the tree started from current element: -print flag
func (c *CmdElement) print() (err error) {
	fmt.Printf("elements in tree name: '%s', id: %s\n%s\n",
		c.tree, c.tree.Id(), c.element)
	return
}

// path prints path from current element to selected in this tree: -path flag
func (c *CmdElement) path() (err error) {

	if err = c.checkArgs(1); err != nil {
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

// selct selects element in current tree by name: -select flag
func (c *CmdElement) selectFlag() (err error) {

	if err = c.checkArgs(1); err != nil {
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
