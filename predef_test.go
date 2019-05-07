package monadgo

import (
	"fmt"
	"reflect"
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

	printGet(p.DefinedAt(reflect.ValueOf(101)))
	printGet(p.Call(reflect.ValueOf(101)).Interface())
	printGet(p.DefinedAt(reflect.ValueOf(10)))
	printGet(p.Call(reflect.ValueOf(100)).Interface())

	// output:
	// true, bool
	// 10201, string
	// false, bool
	// Nothing, *monadgo._nothing
}
