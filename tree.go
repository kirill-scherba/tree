// tree is /thread protected/ any direction multi-child tree implementation
// with way cost child
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
func New[T TreeData](value T) *Tree[T] { return &Tree[T]{} }

// New creates newtree element
func (t *Tree[T]) New(value T) *Element[T] {
	return &Element[T]{
		value:   value,
		ways:    make(waysMap[T]),
		RWMutex: new(sync.RWMutex),
	}
}
