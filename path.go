// Copyright 2023 Kirill Scherba <kirill@scherba.ru>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Path module of Tree package.

package tree

import (
	"fmt"
	"sort"
	"sync"
)

// PathArray contains array of Path and mrthods to process this array
type PathArray[T TreeData] struct {
	arr []Path[T]
	m   sync.RWMutex
}

// Path contains cost of way through path, and path - array of elements
type Path[T TreeData] struct {
	Cost float64
	Path []*Element[T]
}

// String stringify PathArray
func (p *PathArray[T]) String() (str string) {
	p.m.RLock()
	defer p.m.RUnlock()

	for i, path := range p.arr {

		if i > 0 {
			str += "\n"
		}
		str += fmt.Sprintf(
			"%d. From '%s' to '%s', nodes: %d , Cost: %.2f\n",
			i+1,
			path.Path[0].Value(),
			path.Path[len(path.Path)-1].Value(),
			len(path.Path),
			path.Cost)

		for j, e := range p.arr[i].Path {
			if j == 0 {
				str += fmt.Sprintf("  node: %s", e.Value())
				continue
			}
			str += "\n"
			cost, _ := p.arr[i].Path[j-1].Cost(e)
			str += fmt.Sprintf("  node: %s, cost: %.2f", e.Value(), cost)
		}
	}
	return
}

// Append appends path to PathArray
func (p *PathArray[T]) Append(path Path[T]) {
	p.m.Lock()
	defer p.m.Unlock()

	p.arr = append(p.arr, path)
}

// Len returns length of PathArray
func (p *PathArray[T]) Len() int {
	p.m.RLock()
	defer p.m.RUnlock()

	return len(p.arr)
}

// PathSortOptions is path sort options used in Sort function
type PathSortOptions struct {
	SortByCost      bool // default true
	SortByPeers     bool // default true
	SortByCostFirst bool // default true
}

// Sort sorts PathArray
func (p *PathArray[T]) Sort(options ...PathSortOptions) *PathArray[T] {

	p.m.Lock()
	defer p.m.Unlock()

	option := PathSortOptions{true, true, true}
	if len(options) > 0 {
		option = options[0]
	}

	sort.Slice(p.arr, func(i, j int) bool {

		sortByCostAndPeer := func() bool {
			if p.arr[i].Cost == p.arr[j].Cost {
				return len(p.arr[i].Path) < len(p.arr[j].Path)
			}
			return p.arr[i].Cost < p.arr[j].Cost
		}

		if option.SortByCost && option.SortByPeers && option.SortByCostFirst {
			return sortByCostAndPeer()
		}

		if option.SortByCost && option.SortByPeers && !option.SortByCostFirst {
			if len(p.arr[i].Path) == len(p.arr[j].Path) {
				return p.arr[i].Cost < p.arr[j].Cost
			}
			return len(p.arr[i].Path) < len(p.arr[j].Path)
		}

		if option.SortByCost {
			return p.arr[i].Cost < p.arr[j].Cost
		}

		if option.SortByPeers {
			return len(p.arr[i].Path) < len(p.arr[j].Path)
		}

		// Default
		return sortByCostAndPeer()
	})

	return p
}
