package monadgo

import (
	"fmt"
)

type test struct {
	ID   int
	Name string
}

func ExampleOptionOf() {
	n := OptionOf(nil)
	fmt.Println(n)

	n = OptionOf(fmt.Errorf("error"))
	fmt.Println(n)

	n = OptionOf("100")
	fmt.Println(n)

	// Output:
	// Some(Null)
	// Some(error)
	// Some(100)
}

func ExampleOption_Get() {
	x := OptionOf(100)
	printGet(x.Get())

	x = OptionOf("ABC")
	printGet(x.Get())

	x = OptionOf(nil)
	printGet(x.Get())

	printGet(None)
	printGet(None.Get())

	// Output:
	// 100, int
	// ABC, string
	// Null, *monadgo._null
	// None, *monadgo.traitOption
	// Nothing, *monadgo._nothing
}

func ExampleOption_Defined() {
	x := OptionOf(nil)
	fmt.Println(x.Defined())

	x = OptionOf(100)
	fmt.Println(x.Defined())

	fmt.Println(None.Defined())

	// Output:
	// true
	// true
	// false
}

func ExampleOption_Foreach() {
	OptionOf(nil).Foreach(func() {
		fmt.Println("test")
	})

	xx := OptionOf("ABC")
	xx.Foreach(func(x string) {
		fmt.Printf("value is %q\n", x)
	})

	None.Foreach(func() {
		fmt.Println("none")
	})

	// Output:
	// test
	// value is "ABC"
}

func ExampleOption_Forall() {
	v := OptionOf(nil).Forall(func(Null) bool {
		return true
	})
	fmt.Println(v)

	v = OptionOf(100).Forall(func(x int) bool {
		return x >= 100
	})
	fmt.Println(v)

	v = OptionOf(100).Forall(func(x int) bool {
		return x < 10
	})
	fmt.Println(v)

	v = None.Forall(func() bool {
		return false
	})
	fmt.Println(v)

	// Output:
	// true
	// true
	// false
	// true
}

func ExampleOption_Fold() {
	v := OptionOf(nil).Fold(func() int64 {
		return 10
	}, func(Null) int64 {
		return 100
	})
	printGet(v)

	v = OptionOf(int64(100)).Fold("ABC", func(x int64) string {
		return fmt.Sprintf("%d", x*11)
	})
	printGet(v)

	v = None.Fold("ABCD", func() string {
		return "None"
	})
	printGet(v)

	v = None.Fold(func() int { return 1000 }, func() string {
		return "None"
	})
	printGet(v)

	// Output:
	// 100, int64
	// 1100, string
	// ABCD, string
	// 1000, int
}

func ExampleOption_GetOrElse() {
	v := OptionOf(nil).GetOrElse(func() string {
		return "ABC"
	})
	printGet(v)

	x := OptionOf(1000).GetOrElse("ABC")
	printGet(x)

	x = None.GetOrElse(100.0)
	printGet(x)

	x = None.GetOrElse(func() int64 {
		return 101
	})
	printGet(x)

	// Output:
	// Null, *monadgo._null
	// 1000, int
	// 100, float64
	// 101, int64

}

func ExampleOption_Map() {
	v := OptionOf(nil)

	v = v.Map(func(Null) {})
	printGet(v)

	v = OptionOf(100)
	v = v.Map(func(x int) string {
		return fmt.Sprintf("%d", x*x)
	})
	printGet(v)

	v = None.Map(func() int {
		return 100
	})
	printGet(v)

	// Output:
	// Some(Void), *monadgo.traitOption
	// Some(10000), *monadgo.traitOption
	// None, *monadgo.traitOption

}

func ExampleOption_FlatMap() {
	v := OptionOf(nil)

	v = v.FlatMap(func(Null) Option {
		return None
	})
	printGet(v)

	v = OptionOf(100)
	v = v.FlatMap(func(x int) Option {
		return OptionOf(fmt.Sprintf("%d", x*x))
	})
	printGet(v)

	v = None.FlatMap(func() Option {
		return OptionOf(1000)
	})
	printGet(v)

	// Output:
	// None, *monadgo.traitOption
	// Some(10000), *monadgo.traitOption
	// None, *monadgo.traitOption
}

func ExampleOption_OrElse() {

	v := OptionOf(nil).OrElse(func() Option {
		return OptionOf(1000)
	})
	fmt.Println(v)
	printGet(v.Get())

	v = None.OrElse(func() Option {
		return OptionOf(1000)
	})
	fmt.Println(v)
	printGet(v.Get())

	// Output:
	// Some(Null)
	// Null, *monadgo._null
	// Some(1000)
	// 1000, int
}
