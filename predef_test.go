package monadgo

import (
	"fmt"
)

func ExampleUnit() {
	fmt.Println(unit)
	printGet(unit)
	printGet(unit.Get())

	// Output:
	// Void
	// Void, monadgo._unit
	// Void, monadgo._unit
}

func ExampleNull() {
	fmt.Println(null)
	printGet(null)
	printGet(null.Get())

	// Output:
	// Null
	// Null, *monadgo._null
	// <nil>, <nil>
}

func ExampleNothing() {
	fmt.Println(nothing)
	printGet(nothing)

	// Output:
	// Nothing
	// Nothing, *monadgo._nothing
}
