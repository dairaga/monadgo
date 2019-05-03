package monadgo

import (
	"fmt"
)

func ExampleUnit() {
	fmt.Println(unit)
	printGet(unit)
	printGet(unit.Get())

	// Output:
	// void
	// void, monadgo._unit
	// void, monadgo._unit
}

func ExampleNull() {
	fmt.Println(null)
	printGet(null)
	printGet(null.Get())

	// Output:
	// null
	// null, *monadgo._null
	// <nil>, <nil>
}

func ExampleNothing() {
	fmt.Println(nothing)
	printGet(nothing)

	// Output:
	// Nothing
	// Nothing, *monadgo._nothing
}
