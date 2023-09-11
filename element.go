// Copyright 2023 Kirill Scherba <kirill@scherba.ru>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Element module of Tree package.

package tree

import (
	"fmt"
	"sync"

	"golang.org/x/exp/slices"
)

// Element is the multi-chields tree element
type Element[T TreeData] struct {
	value T
	ways  waysMap[T]
	*sync.RWMutex
}
type waysMap[T TreeData] map[*Element[T]]wayOptions

// WayOptions is way options
type WayOptions struct {
	// Cost (weight) of this way (way to this element)
	Cost float64

	// If true than The road available only oneway from this element to selected
	// element (to child), and not available back from selected element (from
	// child) to this element
	Oneway bool
}
type wayOptions struct {
	WayOptions
	wayAllowed bool // Key path is allowed if true
}

// Add adds way from e tree element to the c tree element
func (e *Element[T]) Add(c *Element[T], options ...WayOptions) (*Element[T],
	error) {

	e.Lock()
	c.Lock()
	defer e.Unlock()
	defer c.Unlock()

	// Check input options
	opt := wayOptions{
		WayOptions: WayOptions{
			Cost:   1.0,
			Oneway: false,
		},
		wayAllowed: true,
	}
	if len(options) > 0 {
		opt.WayOptions = options[0]
	}

	// Check the Child in the Elements ways
	if _, ok := e.ways[c]; ok {
		return nil, ErrChildAlreadyAdded
	}
	// Check the Elements in the Childs ways
	if _, ok := c.ways[e]; ok {
		return nil, ErrChildAlreadyAdded
	}

	// Add ways
	e.ways[c] = opt
	if opt.Oneway {
		opt.wayAllowed = false
	}
	c.ways[e] = opt

	return c, nil
}

// Del delete child from tree element
func (e *Element[T]) Del(c *Element[T]) (*Element[T], error) {
	e.Lock()
	c.Lock()
	defer e.Unlock()
	defer c.Unlock()

	delete(e.ways, c)
	delete(c.ways, e)
	return e, nil
}

// Get finds end returns tree elements by name, or nil if not found
func (e *Element[T]) Get(n string) *Element[T] {
	e.RLock()
	defer e.RUnlock()

	// Check current element
	if e.Value().String() == n {
		return e
	}

	// Check connected elements
	for c := range e.ways {
		if c.Value().String() == n {
			return c
		}
	}

	return nil
}

// Remove delete element from tree
func (e *Element[T]) Remove() (*Element[T], error) {
	e.Lock()
	defer e.Unlock()

	for c := range e.ways {
		c.Lock()
		delete(c.ways, e)
		c.Unlock()
	}
	e.ways = make(waysMap[T])
	return e, nil
}

// Value returns elements value
func (e *Element[T]) Value() T { return e.value }

// Cost returns elements way to child cost
func (e *Element[T]) Cost(c *Element[T]) (cost float64, ok bool) {
	e.RLock()
	defer e.RUnlock()
	return e.cost(c)
}

// WayAllowed return true if the path from e to c element is available
func (e *Element[T]) WayAllowed(c *Element[T]) bool {
	e.RLock()
	defer e.RUnlock()
	return e.wayAllowed(c)
}

// Ways returns elments ways maps
func (e *Element[T]) Ways() waysMap[T] {
	e.RLock()
	defer e.RUnlock()
	return e.ways
}

// List returns list of elements in tree
func (e *Element[T]) List() (list []*Element[T]) {
	e.Lock()
	defer e.Unlock()
	list = append(list, e)
	e.list(&list)
	return
}

// list returns list of elements in tree (unsafe)
func (e *Element[T]) list(l *[]*Element[T]) {
	// list = append(list, e)
	for c := range e.ways {
		idx := slices.IndexFunc(*l, func(l *Element[T]) bool { return l == c })
		if idx == -1 {
			*l = append(*l, c)
			c.list(l)
		}
	}

	return
}

// PathTo finds pathes from current element to dst element in the tree
func (e *Element[T]) PathTo(dst *Element[T]) (p *PathArray[T]) {
	e.RLock()
	defer e.RUnlock()

	p = new(PathArray[T])
	e.pathTo(Path[T]{}, p, e, dst)
	return
}

// String returns string with print of tree started from e element
func (e *Element[T]) String() (str string) {
	e.RLock()
	defer e.RUnlock()

	var path []*Element[T]
	str = fmt.Sprintf(". %s", e.Value())
	str += e.string(nil, &path, 0, "")
	return
}

// cost returns elements way to child cost (Unsafe)
func (e *Element[T]) cost(c *Element[T]) (cost float64, ok bool) {
	opt, ok := e.ways[c]
	if ok {
		cost = opt.Cost
	}
	return
}

// wayAllowed return true if the path from e to c element is available (Unsafe)
func (e *Element[T]) wayAllowed(c *Element[T]) bool {
	opt, ok := e.ways[c]
	if ok && opt.wayAllowed {
		return true
	}
	return false
}

// PathTo finds pathes from current element to dst element in the tree (Unsafe)
func (e *Element[T]) pathTo(path Path[T], parr *PathArray[T], next, dst *Element[T]) {

	// Check that next element already exists in path and return error if so
	for i := range path.Path {
		if path.Path[i] == next {
			// Error: path not found because next element already exists
			// fmt.Printf("reject\n")
			return
		}
	}

	// Make sure the path to the next element is available
	if e != next && !e.wayAllowed(next) {
		// Error: path not found because next elements path not allowed
		return
	}

	// Add next element to the path
	cost, _ := e.cost(next)
	path.Cost += cost
	path.Path = append(path.Path, next)

	// Check that path completed
	if next == dst {
		// Done: path complited addit to path array
		// *parr = append(*parr, path)
		parr.Append(path)
		return
	}

	// Create(copy) new pathes to find dst element in childs of next
	var wg sync.WaitGroup
	for child := range next.ways {

		var dstPath Path[T]
		dstPath.Cost = path.Cost
		dstPath.Path = make([]*Element[T], len(path.Path))
		copy(dstPath.Path, path.Path)

		wg.Add(1)
		go func(child *Element[T]) {
			next.pathTo(dstPath, parr, child, dst)
			wg.Done()
		}(child)
	}
	wg.Wait()
}

// String returns string with print of tree started from e element (Unsafe)
func (e *Element[T]) string(parent *Element[T], path *[]*Element[T], level int,
	sline string) (str string) {

	// Check that element is already in path
	for i := range *path {
		if (*path)[i] == e {
			str += " 🡡" // " ⮉"
			return
		}
	}
	*path = append(*path, e)

	const (
		doesNotShowNotallowed = false
		showLevel             = false
	)

	var i int
	lenWays := len(e.ways)
	for c, options := range e.ways {

		i++

		// Skip parent
		if c == parent {
			continue
		}

		// Create wayAllowed text
		var wayAllowed string
		if !e.wayAllowed(c) {
			if doesNotShowNotallowed {
				continue
			}
			wayAllowed = " (way not allowed)"
		} else if options.Oneway {
			wayAllowed = " (one way road)"
		}

		// Create current and next level line
		ch1, spc, ln, ver := "├", "   ", "──", "│"
		if i == lenWays {
			ch1 = "└"
			ver = " "
		}
		line := sline + ch1 + ln
		nextSline := sline + ver + spc

		// Print tree branch
		cost, _ := e.cost(c)
		var levelStr string
		if showLevel {
			levelStr = fmt.Sprintf(" (%d)", level)
		}
		str += fmt.Sprintf("\n%s%s %s, cost: %.2f%s", line, levelStr, c.Value(),
			cost, wayAllowed)

		// Process children
		str += c.string(e, path, level+1, nextSline)
	}
	return
}
