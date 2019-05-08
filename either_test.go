package monadgo

import (
	"fmt"
)

func ExampleEitherOf() {
	e := EitherOf(nil)
	printGet(e.IsLeft())
	printGet(e.IsRight())
	printGet(e)
	printGet(e.Get())

	e = EitherOf(fmt.Errorf("error"))
	printGet(e.IsLeft())
	printGet(e.IsRight())
	printGet(e)
	printGet(e.Get())

	e = EitherOf(true)
	printGet(e.IsLeft())
	printGet(e.IsRight())
	printGet(e)
	printGet(e.Get())

	e = EitherOf(false)
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
	e := Either1Of(100, nil)
	printGet(e.IsLeft())
	printGet(e.IsRight())
	printGet(e)
	printGet(e.Get())

	e = Either1Of(10, true)
	printGet(e.IsLeft())
	printGet(e.IsRight())
	printGet(e)
	printGet(e.Get())

	e = Either1Of(0, fmt.Errorf("error"))
	printGet(e.IsLeft())
	printGet(e.IsRight())
	printGet(e)
	printGet(e.Get())

	e = Either1Of(0, false)
	printGet(e.IsLeft())
	printGet(e.IsRight())
	printGet(e)
	printGet(e.Get())

	e = Either1Of(nil, true)
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
	e := Either2Of(100, "ABC", nil)
	printGet(e.IsLeft())
	printGet(e.IsRight())
	printGet(e)
	printGet(e.Get())
	printGet(e.Get().(Tuple2).V1())
	printGet(e.Get().(Tuple2).V2())

	e = Either2Of(10, "AB", true)
	printGet(e.IsLeft())
	printGet(e.IsRight())
	printGet(e)
	printGet(e.Get())

	e = Either2Of(0, "", fmt.Errorf("error"))
	printGet(e.IsLeft())
	printGet(e.IsRight())
	printGet(e)
	printGet(e.Get())

	e = Either2Of(0, "", false)
	printGet(e.IsLeft())
	printGet(e.IsRight())
	printGet(e)
	printGet(e.Get())

	// Output:
	// false, bool
	// true, bool
	// Right((100,ABC)), *monadgo.traitEither
	// (100,ABC), monadgo._tuple2
	// 100, int
	// ABC, string
	// false, bool
	// true, bool
	// Right((10,AB)), *monadgo.traitEither
	// (10,AB), monadgo._tuple2
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
	EitherOf(false).Foreach(func() {
		fmt.Println("false")
	})

	EitherOf(true).Foreach(func(x bool) {
		printGet(x)
	})

	Either1Of("ok", true).Foreach(func(x string) {
		printGet(x)
	})

	Either2Of("ok", 100, nil).Foreach(func(t Tuple2) {
		printGet(t.V1())
		printGet(t.V2())
	})

	Either2Of("ok", 100, nil).Foreach(func(t Tuple) {
		printGet(t)
	})

	Either2Of("ok", 100, nil).Foreach(func(x1 string, x2 int) {
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

func ExampleEither_Forall() {
	v := EitherOf(false).Forall(func(x bool) bool {
		return false
	})
	printGet(v)

	v = EitherOf(nil).Forall(func() bool {
		return false
	})
	printGet(v)

	v = Either1Of("ABC", true).Forall(func(x string) bool {
		return len(x) < 2
	})
	printGet(v)

	v = Either2Of("ABC", 10, true).Forall(func(x1 string, x2 int) bool {
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
	v := EitherOf(false).Fold(func(bool) int {
		return 100
	}, func(x int, b bool) int {
		return x + 100
	})
	printGet(v)

	v = EitherOf(fmt.Errorf("error")).Fold(func(err error) int {
		return 20
	}, func(x int) int {
		return x + 10
	})
	printGet(v)

	v = EitherOf(nil).Fold(func() string {
		return "ABC"
	}, func() string {
		return "DEF"
	})
	printGet(v)

	v = EitherOf(true).Fold(func(x bool) (string, bool) {
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
	// (failure,false), monadgo._pair

}

func ExampleEither_GetOrElse() {
	v := EitherOf(false).GetOrElse(func() int {
		return 10
	})
	printGet(v)

	v = EitherOf(false).GetOrElse(100)
	printGet(v)

	v = EitherOf(true).GetOrElse(200)
	printGet(v)

	// Output:
	// 10, int
	// 100, int
	// true, bool
}

func ExampleEither_Map() {
	v := EitherOf(false).Map(func(x bool) string {
		return fmt.Sprintf("ok:%v", x)
	})
	printGet(v)

	v = EitherOf(fmt.Errorf("error")).Map(func(x Null) string {
		return fmt.Sprintf("ok:%v", x)
	})
	printGet(v)

	v = EitherOf(true).Map(func(x bool) string {
		return fmt.Sprintf("ok:%v", x)
	})
	printGet(v)
	printGet(v.Get())

	v = EitherOf((error)(nil)).Map(func(x Null) string {
		return fmt.Sprintf("ok:%v", x)
	})
	printGet(v)
	printGet(v.Get())

	v = EitherOf(true).Map(func(x bool) (string, int) {
		return fmt.Sprintf("ok:%v", x), 1
	})
	printGet(v)
	printGet(v.Get())

	v = EitherOf(true).Map(func(x bool) fmt.Stringer {
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
	// (ok:true,1), monadgo._pair
	// Right(Null), *monadgo.traitEither
	// Null, *monadgo._null

}

func ExampleEither_FlatMap() {
	v := EitherOf(false).FlatMap(func(x bool) Either {
		return Either1Of(100, true)
	})
	printGet(v)
	printGet(v.Get())

	v = EitherOf(true).FlatMap(func(x bool) Either {
		return Either1Of(100, true)
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
	v := EitherOf(false).ToOption()
	printGet(v)

	v = EitherOf(fmt.Errorf("error")).ToOption()
	printGet(v)

	v = EitherOf(nil).ToOption()
	printGet(v)

	v = Either1Of(10, true).ToOption()
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
