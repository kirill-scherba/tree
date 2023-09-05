package tree

import (
	"fmt"
	"sync"
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

	// If true than The road avalable only oneway from this element to selected
	// element (to child), and not avalable back from selected element (from
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

// cost returns elements way to child cost (Unsafe)
func (e *Element[T]) cost(c *Element[T]) (cost float64, ok bool) {
	opt, ok := e.ways[c]
	if ok {
		cost = opt.Cost
	}
	return
}

// WayAllowed return true if the path from e to c element is available
func (e *Element[T]) WayAllowed(c *Element[T]) bool {
	e.RLock()
	defer e.RUnlock()
	return e.wayAllowed(c)
}

// wayAllowed return true if the path from e to c element is available (Unsafe)
func (e *Element[T]) wayAllowed(c *Element[T]) bool {
	opt, ok := e.ways[c]
	if ok && opt.wayAllowed {
		return true
	}
	return false
}

// Ways returns elments ways maps
func (e *Element[T]) Ways() waysMap[T] {
	e.RLock()
	defer e.RUnlock()
	return e.ways
}

// PathTo finds pathes from current element to dst element in the tree
func (e *Element[T]) PathTo(dst *Element[T]) (p *PathArray[T]) {
	e.RLock()
	defer e.RUnlock()

	p = new(PathArray[T])
	e.pathTo(Path[T]{}, p, e, dst)
	return
}
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

// String prints the tree started from e element
func (e *Element[T]) String() (str string) {
	e.RLock()
	defer e.RUnlock()

	var path []*Element[T]
	str = fmt.Sprintf(". %s", e.Value())
	str += e.string(&path, 0, "")
	return
}
func (e *Element[T]) string(path *[]*Element[T], level int, sline string) (str string) {

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
		showLevel             = true
	)

	var i int
	lenWays := len(e.ways)
	for c, options := range e.ways {

		i++

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
		str += c.string(path, level+1, nextSline)
	}
	return
}
