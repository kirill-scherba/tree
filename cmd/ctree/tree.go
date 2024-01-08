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
func NewTreeCli(appShort string) (t *Tree, err error) {
	t = &Tree{}

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
	t.addCommands()

	// Create readline based cli menu and add menu items (commands)
	t.menu, err = menu.New(appShort)
	if err != nil {
		err = fmt.Errorf("can't create menu, %s", err)
		return
	}
	t.menu.Add(t.commands...)
	t.batch = menu.NewBatch(t.menu)

	// Create default tree and add it to default tree
	t.tree = tree.New[TreeData]("Default tree")
	t.treeList.add(t.tree)

	return
}

// Command get command by name or nil if not found
func (t Tree) Command(name string) interface{} {
	for i := range t.commands {
		if t.commands[i].Name() == name {
			return t.commands[i]
		}
	}
	return nil
}

// Run command line interface menu
func (t Tree) Run() {
	t.menu.Run()
}

// BatchRun run batch file
func (t Tree) BatchRun(appShort, name string) error {
	return t.batch.Run(appShort, name)
}

const selectedName = "selected.tmp"

// saveSelectedTree saves the selected tree to the config.
//
// Parameters:
// - appShort: The short name of the application.
//
// Returns:
// - err: The error, if any.
func (t *Tree) saveSelectedTree(appShort string) (err error) {
	// Get config folder
	f, err := t.getSelectedConfig(appShort)
	if err != nil {
		return
	}

	// Write selected tree name to the file
	err = os.WriteFile(f, []byte(t.tree.Name()), 0644)

	return
}

// loadSelectedTree loads the selected tree from the config folder.
//
// Parameters:
// - appShort: The short name of the application.
//
// It returns the name of the selected tree and any error encountered.
func (t *Tree) loadSelectedTree(appShort string) (name string, err error) {
	// Get config folder
	f, err := t.getSelectedConfig(appShort)
	if err != nil {
		return
	}

	// Read selected tree name from the file
	data, err := os.ReadFile(f)
	if err != nil {
		return
	}
	name = string(data)

	return
}

// getSelectedConfig returns the path to the selected tree config file.
func (t *Tree) getSelectedConfig(appShort string) (f string, err error) {
	f, err = os.UserConfigDir()
	if err != nil {
		return
	}
	f += "/" + teonet.ConfigDir + "/" + appShort + "/" + selectedName
	return
}

// setUsage set flags usage helper function
func (t Tree) setUsage(usage string, flags *flag.FlagSet, help ...string) {
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
func (t Tree) NewFlagSet(name, usage string, help ...string) (flags *flag.FlagSet) {
	flags = flag.NewFlagSet(name, flag.ContinueOnError)
	t.setUsage(name+" "+usage, flags, help...)
	return
}
