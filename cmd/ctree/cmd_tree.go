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
	"golang.org/x/exp/slices"
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
		"-list", "-save", "-load", "-new", "-select", "-" + cmdHelp,
	})
}
func (c CmdTree) Exec(line string) (err error) {
	// Define and parse flags and get arguments
	var new, save, load, list, selct bool
	flags := c.NewFlagSet(c.Name(), c.Usage(), c.Help())
	flags.BoolVar(&new, "new", new, "create new tree")
	flags.BoolVar(&save, "save", save, "save current tree")
	flags.BoolVar(&load, "load", load, "load current tree")
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
		c.element = c.tree.New(TreeData("My first node"))
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
		c.batch.Run(appShort, c.tree.Name()+".conf")
		return

	// Save current tree: -save flag
	case save:
		err = c.save()

	// Load current tree: -load flag
	case load:
		name := c.tree.Name()
		if argc > 0 {
			name = strings.Join(args, " ")
		}
		c.batch.Run(appShort, name+".conf")

	// Print current tre name and id
	default:
		fmt.Printf("current tree name: '%s', id: %s\n", c.tree, c.tree.Id())
	}

	return
}

// save saves tree elements to the batch file: -save flag
func (c CmdTree) save() (err error) {

	// Get children func
	var getChildren func(path *ListPairs, e *tree.Element[TreeData]) string
	getChildren = func(path *ListPairs, e *tree.Element[TreeData]) (str string) {
		var i int

		type childChData struct {
			path *ListPairs
			e    *tree.Element[TreeData]
		}
		childCh := make(chan childChData, len(e.Ways()))
		for child := range e.Ways() {
			// Skip if way not allowed or already exists in path
			if !e.WayAllowed(child) || path.Contains(e, child) {
				continue
			}
			path.Add(e, child)

			// Print Element Select
			if i == 0 {
				str += fmt.Sprintf("\nelement -select %s\n", e.Value())
				i++
			}

			// Check cost and oneway and pront Element Add
			cost, _ := e.Cost(child)
			oneway, _ := e.Oneway(child)
			costStr := func() (costStr string) {
				if !(cost == 1.0 && !oneway) {
					costStr = fmt.Sprintf(", %f", cost)
				}
				return
			}
			onewayStr := func() (onewayStr string) {
				if oneway {
					onewayStr = fmt.Sprintf(", oneway")
				}
				return
			}
			str += fmt.Sprintf("element -add %s%s%s\n", child.Value(), costStr(), onewayStr())

			// Send path and child to channel to process it after all
			// children have been processed
			childCh <- childChData{path, child}
		}
		close(childCh)

		for d := range childCh {
			str += getChildren(d.path, d.e)
		}

		return str
	}

	// Create string with tree commands
	str := fmt.Sprintf("element -new %s\n", c.element.Value())
	str += getChildren(&ListPairs{}, c.element)
	str += fmt.Sprintf("\nelement -select %s\n", c.element.Value())
	fmt.Print(str)

	// Save tree commands to batch file
	batch := strings.Split(str, "\n")
	c.batch.Save(appShort, c.tree.Name()+".conf", "", batch)

	return
}

// ListPairs is array of pair of elements
type ListPairs []Pair

// Pair is array of pair of elements
type Pair [2]*tree.Element[TreeData]

// Add adds input pair of elements to the ListPairs slice
func (l *ListPairs) Add(e, c *tree.Element[TreeData]) {
	*l = append(*l, Pair{e, c})
}

// Contains returns true if input pair of elements exists in the ListPairs slice
func (l *ListPairs) Contains(e, c *tree.Element[TreeData]) bool {
	return slices.ContainsFunc(*l, func(p Pair) bool {
		if p[0] == e && p[1] == c || p[0] == c && p[1] == e {
			return true
		}
		return false
	})
}
