package monadgo

import (
	"fmt"
)

func ExampleTryOf() {
	t := TryOf(nil)
	fmt.Println(t.OK())
	fmt.Println(t.Failed())
	printGet(t)
	printGet(t.Get())

	t = TryOf(fmt.Errorf("error"))
	fmt.Println(t.OK())
	fmt.Println(t.Failed())
	printGet(t)
	printGet(t.Get())

	// Output:
	// true
	// false
	// Success(Null), *monadgo.traitTry
	// Null, *monadgo._null
	// false
	// true
	// Failure(error), *monadgo.traitTry
	// error, *errors.errorString
}

func ExampleTry1Of() {
	t := Try1Of(100, nil)
	fmt.Println(t.OK())
	fmt.Println(t.Failed())
	printGet(t)
	printGet(t.Get())

	t = Try1Of(10, true)
	fmt.Println(t.OK())
	fmt.Println(t.Failed())
	printGet(t)
	printGet(t.Get())

	t = Try1Of(0, fmt.Errorf("error"))
	fmt.Println(t.OK())
	fmt.Println(t.Failed())
	printGet(t)
	printGet(t.Get())

	t = Try1Of(0, false)
	fmt.Println(t.OK())
	fmt.Println(t.Failed())
	printGet(t)
	printGet(t.Get())

	t = Try1Of(nil, true)
	fmt.Println(t.OK())
	fmt.Println(t.Failed())
	printGet(t)
	printGet(t.Get())

	// Output:
	// true
	// false
	// Success(100), *monadgo.traitTry
	// 100, int
	// true
	// false
	// Success(10), *monadgo.traitTry
	// 10, int
	// false
	// true
	// Failure(error), *monadgo.traitTry
	// error, *errors.errorString
	// false
	// true
	// Failure(false), *monadgo.traitTry
	// false, bool
	// true
	// false
	// Success(Null), *monadgo.traitTry
	// Null, *monadgo._null
}

func ExampleTry2Of() {
	t := Try2Of(100, "ABC", nil)
	fmt.Println(t.OK())
	fmt.Println(t.Failed())
	printGet(t)
	printGet(t.Get())
	printGet(t.Get().(Tuple2).V1())
	printGet(t.Get().(Tuple2).V2())

	t = Try2Of(10, "AB", true)
	fmt.Println(t.OK())
	fmt.Println(t.Failed())
	printGet(t)
	printGet(t.Get())

	t = Try2Of(0, "", fmt.Errorf("error"))
	fmt.Println(t.OK())
	fmt.Println(t.Failed())
	printGet(t)
	printGet(t.Get())

	t = Try2Of(0, "", false)
	fmt.Println(t.OK())
	fmt.Println(t.Failed())
	printGet(t)
	printGet(t.Get())

	// Output:
	// true
	// false
	// Success((100,ABC)), *monadgo.traitTry
	// (100,ABC), monadgo._tuple2
	// 100, int
	// ABC, string
	// true
	// false
	// Success((10,AB)), *monadgo.traitTry
	// (10,AB), monadgo._tuple2
	// false
	// true
	// Failure(error), *monadgo.traitTry
	// error, *errors.errorString
	// false
	// true
	// Failure(false), *monadgo.traitTry
	// false, bool
}

func ExampleTry_Foreach() {
	TryOf(false).Foreach(func() {
		fmt.Println("false")
	})

	TryOf(true).Foreach(func(x bool) {
		printGet(x)
	})

	Try1Of("ok", true).Foreach(func(x string) {
		printGet(x)
	})

	Try2Of("ok", 100, nil).Foreach(func(t Tuple2) {
		printGet(t.V1())
		printGet(t.V2())
	})

	Try2Of("ok", 100, nil).Foreach(func(t Tuple) {
		printGet(t)
	})

	Try2Of("ok", 100, nil).Foreach(func(x1 string, x2 int) {
		printGet(x1)
		printGet(x2)
	})

	// Output:
	// true, bool
	// ok, string
	// ok, string
	// 100, int
	// (ok,100), monadgo._tuple2
	// ok, string
	// 100, int
}

func ExampleTry_Fold() {
	v := TryOf(false).Fold(func(bool) int {
		return 100
	}, func(x int, b bool) int {
		return x + 100
	})
	printGet(v)

	v = TryOf(fmt.Errorf("error")).Fold(20, func(x int) int {
		return x + 10
	})
	printGet(v)

	v = TryOf(nil).Fold(func() string {
		return "ABC"
	}, func(x Null) string {
		return "DEF"
	})
	printGet(v)

	v = TryOf(true).Fold(func(x bool) (string, bool) {
		if x {
			return "ok", true
		}
		return "not ok", false
	}, func(x bool) (string, bool) {
		return "failure", false
	})
	printGet(v)

	// Output:
	// 100, int
	// 20, int
	// DEF, string
	// (not ok,false), monadgo._pair

}

func ExampleTry_GetOrElse() {
	v := TryOf(false).GetOrElse(func() int {
		return 10
	})
	printGet(v)

	v = TryOf(false).GetOrElse(100)
	printGet(v)

	v = TryOf(true).GetOrElse(200)
	printGet(v)

	// Output:
	// 10, int
	// 100, int
	// true, bool
}

func ExampleTry_OrElse() {
	v := TryOf(false).OrElse(func() Try {
		return Try1Of(100, true)
	})
	printGet(v)

	v = TryOf(true).OrElse(func() Try {
		return Try1Of(100, true)
	})
	printGet(v)

	// Output:
	// Success(100), *monadgo.traitTry
	// Success(true), *monadgo.traitTry
}

func ExampleTry_Map() {
	v := TryOf(false).Map(func(x bool) string {
		return fmt.Sprintf("ok:%v", x)
	})
	printGet(v)

	v = TryOf(fmt.Errorf("error")).Map(func(x Null) string {
		return fmt.Sprintf("ok:%v", x)
	})
	printGet(v)

	v = TryOf(true).Map(func(x bool) string {
		return fmt.Sprintf("ok:%v", x)
	})
	printGet(v)
	printGet(v.Get())

	v = TryOf((error)(nil)).Map(func(x Null) string {
		return fmt.Sprintf("ok:%v", x)
	})
	printGet(v)
	printGet(v.Get())

	v = TryOf(true).Map(func(x bool) (string, int) {
		return fmt.Sprintf("ok:%v", x), 1
	})
	printGet(v)
	printGet(v.Get())

	v = TryOf(true).Map(func(x bool) fmt.Stringer {
		return nil
	})
	printGet(v)
	printGet(v.Get())

	// Output:
	// Failure(false), *monadgo.traitTry
	// Failure(error), *monadgo.traitTry
	// Success(ok:true), *monadgo.traitTry
	// ok:true, string
	// Success(ok:Null), *monadgo.traitTry
	// ok:Null, string
	// Success((ok:true,1)), *monadgo.traitTry
	// (ok:true,1), monadgo._pair
	// Success(Null), *monadgo.traitTry
	// Null, *monadgo._null

}

func ExampleTry_FlatMap() {
	v := TryOf(false).FlatMap(func(x bool) Try {
		return Try1Of(100, true)
	})
	printGet(v)
	printGet(v.Get())

	v = TryOf(true).FlatMap(func(x bool) Try {
		return Try1Of(100, true)
	})
	printGet(v)
	printGet(v.Get())

	// Output:
	// Failure(false), *monadgo.traitTry
	// false, bool
	// Success(100), *monadgo.traitTry
	// 100, int
}

func ExampleTry_ToOption() {
	v := TryOf(false).ToOption()
	printGet(v)

	v = TryOf(fmt.Errorf("error")).ToOption()
	printGet(v)

	v = TryOf(nil).ToOption()
	printGet(v)

	v = Try1Of(10, true).ToOption()
	printGet(v)

	// Output:
	// None, *monadgo.traitOption
	// None, *monadgo.traitOption
	// Some(Null), *monadgo.traitOption
	// Some(10), *monadgo.traitOption

}
