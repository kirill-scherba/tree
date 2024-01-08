// Copyright 2023 Kirill Scherba <kirill@scherba.ru>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Tree CLI application. List of Trees.

package main

import (
	"fmt"
	"os"
	"path"

	"github.com/kirill-scherba/tree"
	"github.com/teonet-go/teonet"
)

const confExt = ".conf"

// TreesList is a list of trees
type TreesList []*tree.Tree[TreeData]

// String returns a string representation of the list of trees
func (t TreesList) String() string {
	var s string
	for i := range t {
		s += fmt.Sprintf("%s\n", t[i].Name())
	}
	return s
}

// add adds tree to the list of trees
func (t *TreesList) add(tree *tree.Tree[TreeData]) {
	*t = append(*t, tree)
}

// get returns tree from the list of trees by name or nil if not found
func (t TreesList) get(name string) *tree.Tree[TreeData] {
	for i := range t {
		if t[i].Name() == name {
			return t[i]
		}
	}
	return nil
}

// load loads tree from the config folder
func (t *TreesList) load(appShort string) {

	// Get config folder
	f, err := os.UserConfigDir()
	if err != nil {
		return
	}
	f += "/" + teonet.ConfigDir + "/" + appShort + "/"

	// Get list of files in config folder
	files, err := os.ReadDir(f)
	if err != nil {
		return
	}

	fmt.Printf("load trees name from config\n\n")

	// Get files with .conf extension and add it to the tree list
	for _, file := range files {
		if path.Ext(file.Name()) == confExt {
			name := path.Base(file.Name())
			name = name[:len(name)-len(confExt)]

			if t.get(name) == nil {
				tree := tree.New[TreeData](name)
				t.add(tree)
				fmt.Printf("tree -new %s\n", name)
				fmt.Printf("> added tree '%s'\n", name)
			}
		}
	}
}
