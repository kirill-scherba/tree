package main

import (
	"fmt"

	"github.com/kirill-scherba/tree"
)

// TreeData is multi-chields tree data structure
type TreeData string

// String is mandatory TreeData method which return element name
func (t TreeData) String() string {
	return string(t)
}

func main() {

	// Create new multi-chields t
	t := tree.New[TreeData]("My first element")

	// Create new tree elements (first element and end point element)
	e := t.New("My first element")
	ep := t.New("End point")

	// Create children of e (first element) element
	ch1, _ := e.Add(t.New("My first child"))
	ch2, _ := e.Add(t.New("My second child"), tree.WayOptions{Cost: 3.0})
	ch4, _ := e.Add(t.New("My fourth child"), tree.WayOptions{Cost: 3.0})

	// Create sub children elements
	ch3, _ := ch2.Add(t.New("Some third child"))
	ch3.Add(ch4)

	// Add children to ep (end point) element
	ch1.Add(ch4)
	ch2.Add(ep)
	ch4.Add(ep)

	// Set oneway path from ep (end point) to e (first element) element
	ep.Add(e, tree.WayOptions{Cost: 5.0, Oneway: true})

	// Print tree started from e (first element) element
	fmt.Printf("\nPrint tree:\n%v\n\n", e)

	// return

	testsLen := 1 // Set testsLen to 3 to execute additional test
	for i := 0; i < testsLen; i++ {

		// Print parents ways
		fmt.Printf("\n%s ways:\n", e.Value())
		for c := range e.Ways() {
			cost, _ := e.Cost(c)
			fmt.Printf("  way to '%s' cost: %.2f, way allowed: %v\n",
				c.Value(), cost, e.WayAllowed(c))
		}

		// Print child ways
		for child := range e.Ways() {
			fmt.Printf("\n%s ways:\n", child.Value())
			for c := range child.Ways() {
				cost, _ := child.Cost(c)
				fmt.Printf("  way to '%s' cost: %.2f, way allowed: %v\n",
					c.Value(), cost, child.WayAllowed(c))
			}
		}

		// Find path
		p := e.PathTo(ep).Sort()
		if p.Len() > 0 {
			fmt.Printf("\nPath array:\n%s\n", p.String())
		}

		if testsLen > 1 {
			// Delete element elements. To execute this test set testsLen to 3
			switch i {
			case 0:
				// Delete child
				e.Del(ch1)
				fmt.Printf("\n----\nDelete child '%v', from '%v'\n",
					ch1.Value(), e.Value())

			case 1:
				// Delete element
				ch2.Remove()
				fmt.Printf("\n----\nDelete element '%v'\n", ch2.Value())
			}
		}
	}
}
