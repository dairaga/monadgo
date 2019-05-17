package monadgo

import (
	"fmt"
)

func ExampleEitherOf() {
	e := RightOf(nil)
	printGet(e.IsLeft())
	printGet(e.IsRight())
	printGet(e)
	printGet(e.Get())

	e = LeftOf(fmt.Errorf("error"))
	printGet(e.IsLeft())
	printGet(e.IsRight())
	printGet(e)
	printGet(e.Get())

	e = RightOf(true)
	printGet(e.IsLeft())
	printGet(e.IsRight())
	printGet(e)
	printGet(e.Get())

	e = LeftOf(false)
	printGet(e.IsLeft())
	printGet(e.IsRight())
	printGet(e)
	printGet(e.Get())

	// Output:
	// false, bool
	// true, bool
	// Right(Null), *monadgo.traitEither
	// Null, *monadgo._null
	// true, bool
	// false, bool
	// Left(error), *monadgo.traitEither
	// error, *errors.errorString
	// false, bool
	// true, bool
	// Right(true), *monadgo.traitEither
	// true, bool
	// true, bool
	// false, bool
	// Left(false), *monadgo.traitEither
	// false, bool
}

func ExampleEither1Of() {
	e := RightOf(100)
	printGet(e.IsLeft())
	printGet(e.IsRight())
	printGet(e)
	printGet(e.Get())

	e = RightOf(10)
	printGet(e.IsLeft())
	printGet(e.IsRight())
	printGet(e)
	printGet(e.Get())

	e = LeftOf(fmt.Errorf("error"))
	printGet(e.IsLeft())
	printGet(e.IsRight())
	printGet(e)
	printGet(e.Get())

	e = LeftOf(false)
	printGet(e.IsLeft())
	printGet(e.IsRight())
	printGet(e)
	printGet(e.Get())

	e = RightOf(nil)
	printGet(e.IsLeft())
	printGet(e.IsRight())
	printGet(e)
	printGet(e.Get())

	// Output:
	// false, bool
	// true, bool
	// Right(100), *monadgo.traitEither
	// 100, int
	// false, bool
	// true, bool
	// Right(10), *monadgo.traitEither
	// 10, int
	// true, bool
	// false, bool
	// Left(error), *monadgo.traitEither
	// error, *errors.errorString
	// true, bool
	// false, bool
	// Left(false), *monadgo.traitEither
	// false, bool
	// false, bool
	// true, bool
	// Right(Null), *monadgo.traitEither
	// Null, *monadgo._null
}

func ExampleEither2Of() {
	e := RightOf(100, "ABC")
	printGet(e.IsLeft())
	printGet(e.IsRight())
	printGet(e)
	printGet(e.Get())
	printGet(e.Get().(Tuple2).V1())
	printGet(e.Get().(Tuple2).V2())

	e = RightOf(10, "AB")
	printGet(e.IsLeft())
	printGet(e.IsRight())
	printGet(e)
	printGet(e.Get())

	e = LeftOf(fmt.Errorf("error"))
	printGet(e.IsLeft())
	printGet(e.IsRight())
	printGet(e)
	printGet(e.Get())

	e = LeftOf(false)
	printGet(e.IsLeft())
	printGet(e.IsRight())
	printGet(e)
	printGet(e.Get())

	// Output:
	// false, bool
	// true, bool
	// Right((100,ABC)), *monadgo.traitEither
	// (100,ABC), monadgo.Tuple2
	// 100, int
	// ABC, string
	// false, bool
	// true, bool
	// Right((10,AB)), *monadgo.traitEither
	// (10,AB), monadgo.Tuple2
	// true, bool
	// false, bool
	// Left(error), *monadgo.traitEither
	// error, *errors.errorString
	// true, bool
	// false, bool
	// Left(false), *monadgo.traitEither
	// false, bool
}

func ExampleEither_Foreach() {
	LeftOf(false).Foreach(func() {
		fmt.Println("false")
	})

	RightOf(true).Foreach(func(x bool) {
		printGet(x)
	})

	RightOf("ok").Foreach(func(x string) {
		printGet(x)
	})

	RightOf("ok", 100).Foreach(func(t Tuple2) {
		printGet(t.V1())
		printGet(t.V2())
	})

	RightOf("ok", 100).Foreach(func(t Tuple) {
		printGet(t)
	})

	RightOf("ok", 100).Foreach(func(x1 string, x2 int) {
		printGet(x1)
		printGet(x2)
	})

	// Output:
	// true, bool
	// ok, string
	// ok, string
	// 100, int
	// (ok,100), monadgo.Tuple2
	// ok, string
	// 100, int
}

func ExampleEither_Forall() {
	v := LeftOf(false).Forall(func(x bool) bool {
		return false
	})
	printGet(v)

	v = RightOf(nil).Forall(func() bool {
		return false
	})
	printGet(v)

	v = RightOf("ABC").Forall(func(x string) bool {
		return len(x) < 2
	})
	printGet(v)

	v = RightOf("ABC", 10).Forall(func(x1 string, x2 int) bool {
		return len(x1) > 2 && x2 > 5
	})
	printGet(v)

	// Output:
	// true, bool
	// false, bool
	// false, bool
	// true, bool
}

func ExampleEither_Fold() {
	v := LeftOf(false).Fold(func(bool) int {
		return 100
	}, func(x int, b bool) int {
		return x + 100
	})
	printGet(v)

	v = LeftOf(fmt.Errorf("error")).Fold(func(err error) int {
		return 20
	}, func(x int) int {
		return x + 10
	})
	printGet(v)

	v = RightOf(nil).Fold(func() string {
		return "ABC"
	}, func() string {
		return "DEF"
	})
	printGet(v)

	v = RightOf(true).Fold(func(x bool) (string, bool) {
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
	// (failure,false), monadgo.Pair

}

func ExampleEither_GetOrElse() {
	v := LeftOf(false).GetOrElse(func() int {
		return 10
	})
	printGet(v)

	v = LeftOf(false).GetOrElse(100)
	printGet(v)

	v = RightOf(true).GetOrElse(200)
	printGet(v)

	// Output:
	// 10, int
	// 100, int
	// true, bool
}

func ExampleEither_Map() {
	v := LeftOf(false).Map(func(x bool) string {
		return fmt.Sprintf("ok:%v", x)
	})
	printGet(v)

	v = LeftOf(fmt.Errorf("error")).Map(func(x Null) string {
		return fmt.Sprintf("ok:%v", x)
	})
	printGet(v)

	v = RightOf(true).Map(func(x bool) string {
		return fmt.Sprintf("ok:%v", x)
	})
	printGet(v)
	printGet(v.Get())

	v = RightOf((error)(nil)).Map(func(x Null) string {
		return fmt.Sprintf("ok:%v", x)
	})
	printGet(v)
	printGet(v.Get())

	v = RightOf(true).Map(func(x bool) (string, int) {
		return fmt.Sprintf("ok:%v", x), 1
	})
	printGet(v)
	printGet(v.Get())

	v = RightOf(true).Map(func(x bool) fmt.Stringer {
		return nil
	})
	printGet(v)
	printGet(v.Get())

	// Output:
	// Left(false), *monadgo.traitEither
	// Left(error), *monadgo.traitEither
	// Right(ok:true), *monadgo.traitEither
	// ok:true, string
	// Right(ok:Null), *monadgo.traitEither
	// ok:Null, string
	// Right((ok:true,1)), *monadgo.traitEither
	// (ok:true,1), monadgo.Pair
	// Right(Null), *monadgo.traitEither
	// Null, *monadgo._null

}

func ExampleEither_FlatMap() {
	v := LeftOf(false).FlatMap(func(x bool) Either {
		return RightOf(100)
	})
	printGet(v)
	printGet(v.Get())

	v = RightOf(true).FlatMap(func(x bool) Either {
		return RightOf(100)
	})
	printGet(v)
	printGet(v.Get())

	// Output:
	// Left(false), *monadgo.traitEither
	// false, bool
	// Right(100), *monadgo.traitEither
	// 100, int
}

func ExampleEither_ToOption() {
	v := LeftOf(false).ToOption()
	printGet(v)

	v = LeftOf(fmt.Errorf("error")).ToOption()
	printGet(v)

	v = RightOf(nil).ToOption()
	printGet(v)

	v = RightOf(10).ToOption()
	printGet(v)

	// Output:
	// None, *monadgo.traitOption
	// None, *monadgo.traitOption
	// Some(Null), *monadgo.traitOption
	// Some(10), *monadgo.traitOption

}
func ExampleEither_FilterOrElse() {
	p := func(x int) bool {
		return x >= 50
	}

	z := func() string {
		return "ok"
	}

	x := RightOf(100).FilterOrElse(p, z)
	printGet(x)

	x = RightOf(30).FilterOrElse(p, z)
	printGet(x)

	x = LeftOf(7).FilterOrElse(p, z)
	printGet(x)

	// Output:
	// Right(100), *monadgo.traitEither
	// Left(ok), *monadgo.traitEither
	// Left(7), *monadgo.traitEither

}

func ExampleEither_Exists() {
	v := RightOf(12).Exists(func(x int) bool {
		return x > 10
	}) // true
	printGet(v)

	v = RightOf(7).Exists(func(x int) bool {
		return x > 10
	}) // false
	printGet(v)

	v = LeftOf(12).Exists(func(int) bool {
		return true
	}) // false
	printGet(v)

	// Output:
	// true, bool
	// false, bool
	// false, bool
}
