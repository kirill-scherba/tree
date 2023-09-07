// Copyright 2023 Kirill Scherba <kirill@scherba.ru>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Tree is thread protected any direction multi-child tree implementation
// with way cost child.
package tree

import (
	"errors"
	"sync"

	"github.com/google/uuid"
)

var (
	ErrChildAlreadyAdded  = errors.New("child already added")
	ErrParentAlreadyAdded = errors.New("parent already added")
)

// Tree is the tree methods receiver
type Tree[T TreeData] struct {
	name string
	id   string
}

// TreeData interface represents TreeData required methods
type TreeData interface {
	String() string
}

// New creates new multi-chields tree
func New[T TreeData](opts ...string) *Tree[T] {
	name, id := "", uuid.New().String()
	if len(opts) > 0 {
		name = opts[0]
	} else {
		name = id
	}
	return &Tree[T]{name: name, id: id}
}

// New creates newtree element
func (t *Tree[T]) New(value T) *Element[T] {
	return &Element[T]{
		value:   value,
		ways:    make(waysMap[T]),
		RWMutex: new(sync.RWMutex),
	}
}

// Name returns Name of tree
func (t *Tree[T]) Name() string { return t.name }

// String returns name of string (if name was ommited when tree was created
// then name equal id)
func (t *Tree[T]) String() string { return t.Name() }

// Id return trees id
func (t *Tree[T]) Id() string { return t.id }
