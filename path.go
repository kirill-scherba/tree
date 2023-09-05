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

		str += fmt.Sprintf(
			"%d. From '%s' to '%s', nodes: %d , Cost: %.2f\n",
			i+1,
			path.Path[0].Value(),
			path.Path[len(path.Path)-1].Value(),
			len(path.Path),
			path.Cost)

		for j, e := range p.arr[i].Path {
			var cost float64
			if j > 0 {
				cost, _ = p.arr[i].Path[j-1].Cost(e)
				str += fmt.Sprintf("  node: %s, cost: %.2f\n", e.Value(), cost)
			} else {
				str += fmt.Sprintf("  node: %s\n", e.Value())
			}
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

type PathSortOptions struct {
	SortByCost      bool // default true
	SortByPeers     bool // default true
	SortByCostFirst bool // default true
}

// Sort sorts PathArray
func (p *PathArray[T]) Sort(options ...PathSortOptions) *PathArray[T] {

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
