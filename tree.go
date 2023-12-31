// Copyright 2023 Kirill Scherba <kirill@scherba.ru>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Tree is thread protected any direction multi-child tree implementation
// with way cost child.
package tree

import (
	"errors"
	"sync"
)

var (
	ErrChildAlreadyAdded  = errors.New("child already added")
	ErrParentAlreadyAdded = errors.New("parent already added")
)

// Tree is the tree methods receiver
type Tree[T TreeData] struct{}

// TreeData interface represents TreeData required methods
type TreeData interface {
	String() string
}

// New creates new multi-chields tree
func New[T TreeData]() *Tree[T] { return &Tree[T]{} }

// New creates newtree element
func (t *Tree[T]) New(value T) *Element[T] {
	return &Element[T]{
		value:   value,
		ways:    make(waysMap[T]),
		RWMutex: new(sync.RWMutex),
	}
}
