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

func ExamplePartialFuncOf() {
	p := PartialFuncOf(
		func(x int) bool {
			return x > 100
		},
		func(x int) string {
			return fmt.Sprintf("%d", x*x)
		},
	)

	printGet(p.DefinedAt(101))
	printGet(p.Call(101))
	printGet(p.DefinedAt(10))
	printGet(p.Call(100))

	// output:
	// true, bool
	// 10201, string
	// false, bool
	// Nothing, *monadgo._nothing
}
