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

		// Get children func
		var getChildren func(path *tree.List[TreeData], e *tree.Element[TreeData]) string
		getChildren = func(path *tree.List[TreeData], e *tree.Element[TreeData]) (str string) {
			var i int

			type childChData struct {
				path *tree.List[TreeData]
				e    *tree.Element[TreeData]
			}
			childCh := make(chan childChData, len(e.Ways()))
			for child := range e.Ways() {

				if slices.Contains(*path, child) {
					continue
				}
				*path = append(*path, child)

				if i == 0 {
					str += fmt.Sprintf("\nelement -select %s\n", e.Value())
					i++
				}

				cost, _ := c.element.Cost(child)
				oneway, _ := c.element.Oneway(child)
				str += fmt.Sprintf("element -add %s, %f, %v\n", child.Value(), cost, oneway)

				p := slices.Clone(*path)
				childCh <- childChData{&p, child}
			}
			close(childCh)

			for d := range childCh {
				str += getChildren(d.path, d.e)
			}

			return str
		}

		// Create string with tree commands
		str := fmt.Sprintf("element -new %s\n", c.element.Value())
		// var path = tree.List[TreeData]{c.element}
		str += getChildren(&tree.List[TreeData]{c.element}, c.element)
		str += fmt.Sprintf("\nelement -select %s\n", c.element.Value())
		fmt.Print(str)

		// Save tree commands to batch file
		batch := strings.Split(str, "\n")
		c.batch.Save(appShort, c.tree.Name()+".conf", "", batch)

	// Load current tree: -load flag
	case load:
		c.batch.Run(appShort, c.tree.Name()+".conf")

	// Print current tre name and id
	default:
		fmt.Printf("current tree name: '%s', id: %s\n", c.tree, c.tree.Id())
	}

	return
}
