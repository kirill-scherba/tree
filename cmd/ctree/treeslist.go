// Copyright 2023 Kirill Scherba <kirill@scherba.ru>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Tree CLI application. Trees list module.

package main

import (
	"fmt"

	"github.com/kirill-scherba/tree"
)

// TreesList is a list of loaded trees
type TreesList []*tree.Tree[TreeData]

// String returns a string representation of the tree list
func (t TreesList) String() string {
	var s string
	for i := range t {
		s += fmt.Sprintf("%s - %s\n", t[i].Id(), t[i].Name())
	}
	return s
}

// add adds tree to the tree liat
func (t *TreesList) add(tree *tree.Tree[TreeData]) {
	*t = append(*t, tree)
}
