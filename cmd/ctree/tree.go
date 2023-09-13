// Copyright 2023 Kirill Scherba <kirill@scherba.ru>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Tree CLI application. Tree client.

package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/kirill-scherba/tree"
	"github.com/teonet-go/teonet"
	"github.com/teonet-go/teonet/cmd/teonet/menu"
)

// Tree represents Tree CLI data structure and methods receiver.
type Tree struct {
	treeList TreesList               // Trees array
	tree     *tree.Tree[TreeData]    // Current tree
	element  *tree.Element[TreeData] // Current element
	menu     *menu.Menu              // Commands menu
	commands []menu.Item             // Commands menu items
	batch    *menu.Batch             // Batch files menu object
}

// TreeData is tree elements data structure
type TreeData string

// String is mandatory TreeData method which return element name
func (t TreeData) String() string {
	return string(t)
}

// NewTreeCli create new Tree CLI client
func NewTreeCli(appShort string) (cli *Tree, err error) {
	cli = &Tree{}

	// Create config directory if does not exists
	dir, err := os.UserConfigDir()
	if err != nil {
		dir = os.TempDir()
	}
	path := dir + "/" + teonet.ConfigDir + "/" + appShort
	if _, err = os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err = os.Mkdir(path, os.ModePerm)
		if err != nil {
			err = fmt.Errorf("can't create config directory: %s", err)
			return
		}
	}

	// Add commands
	cli.addCommands()

	// Create readline based cli menu and add menu items (commands)
	cli.menu, err = menu.New(appShort)
	if err != nil {
		err = fmt.Errorf("can't create menu, %s", err)
		return
	}
	cli.menu.Add(cli.commands...)
	cli.batch = menu.NewBatch(cli.menu)

	// Create default tree and add it to default tree
	cli.tree = tree.New[TreeData]()
	cli.treeList.add(cli.tree)

	return
}

// Command get command by name or nil if not found
func (cli Tree) Command(name string) interface{} {
	for i := range cli.commands {
		if cli.commands[i].Name() == name {
			return cli.commands[i]
		}
	}
	return nil
}

// Run command line interface menu
func (cli Tree) Run() {
	cli.menu.Run()
}

// setUsage set flags usage helper function
func (cli Tree) setUsage(usage string, flags *flag.FlagSet, help ...string) {
	savUsage := flags.Usage
	flags.Usage = func() {
		fmt.Print("usage: " + usage + "\n\n")
		if len(help) > 0 && len(help[0]) > 0 {
			fmt.Print(strings.ToUpper(help[0][0:1]) + help[0][1:] + "\n\n")
		}
		savUsage()
		fmt.Println()
	}
}

// NewFlagSet
func (cli Tree) NewFlagSet(name, usage string, help ...string) (flags *flag.FlagSet) {
	flags = flag.NewFlagSet(name, flag.ContinueOnError)
	cli.setUsage(name+" "+usage, flags, help...)
	return
}
