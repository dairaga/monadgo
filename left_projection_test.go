package monadgo

import (
	"fmt"
	"testing"
)

func TestLeftProjection(t *testing.T) {
	e := EitherOf(nil)

	if e.Left().E() != e {
		t.Errorf("prjection faiure for Right")
	}

	e = EitherOf(fmt.Errorf("error"))
	if e.Left().E() != e {
		t.Errorf("prjection faiure for Left")
	}

	e = EitherOf(true)
	if e.Left().E() != e {
		t.Errorf("prjection faiure for Right")
	}

	e = EitherOf(false)
	if e.Left().E() != e {
		t.Errorf("prjection faiure for Left")
	}
}

func ExampleLeftProjection_Foreach() {
	EitherOf(false).Left().Foreach(func(bool) {
		fmt.Println("false")
	})

	EitherOf(fmt.Errorf("x")).Left().Foreach(func(err error) {
		fmt.Println(err)
	})

	EitherOf(true).Left().Foreach(func(x bool) {
		printGet(x)
	})

	Either1Of("ok", true).Left().Foreach(func(x string) {
		printGet(x)
	})

	Either2Of("ok", 100, nil).Left().Foreach(func(t Tuple2) {
		printGet(t.V1())
		printGet(t.V2())
	})

	Either2Of("ok", 100, nil).Left().Foreach(func(t Tuple) {
		printGet(t)
	})

	Either2Of("ok", 100, nil).Left().Foreach(func(x1 string, x2 int) {
		printGet(x1)
		printGet(x2)
	})

	// Output:
	// false
	// x
}

func ExampleLeftProjection_Forall() {
	v := EitherOf(false).Left().Forall(func(x bool) bool {
		// Left and apply returning true.
		return true
	})
	printGet(v)

	v = EitherOf(nil).Left().Forall(func() bool {
		// Right and always return true
		return false
	})
	printGet(v)

	v = Either1Of("ABC", true).Left().Forall(func(x string) bool {
		// Right and always return true
		return len(x) < 2
	})
	printGet(v)

	v = Either2Of("ABC", 10, true).Left().Forall(func(x1 string, x2 int) bool {
		// Right and always return true
		return len(x1) < 2 || x2 < 5
	})
	printGet(v)

	// Output:
	// true, bool
	// true, bool
	// true, bool
	// true, bool
}

func ExampleLeftProjection_GetOrElse() {
	v := EitherOf(false).Left().GetOrElse(func() int {
		return 10
	})
	printGet(v)

	v = EitherOf(false).Left().GetOrElse(100)
	printGet(v)

	v = EitherOf(true).Left().GetOrElse(200)
	printGet(v)

	// Output:
	// false, bool
	// false, bool
	// 200, int
}

func ExampleLeftProjection_Map() {
	v := EitherOf(false).Left().Map(func(x bool) string {
		return fmt.Sprintf("ok:%v", x)
	})
	printGet(v)
	printGet(v.Get())

	v = EitherOf(fmt.Errorf("error")).Left().Map(func(x error) string {
		return fmt.Sprintf("ok:%v", x)
	})
	printGet(v)
	printGet(v.Get())

	v = EitherOf(true).Left().Map(func(x bool) string {
		return fmt.Sprintf("ok:%v", x)
	})
	printGet(v)

	v = EitherOf((error)(nil)).Left().Map(func(x Null) string {
		return fmt.Sprintf("ok:%v", x)
	})
	printGet(v)

	// Output:
	// Left(ok:false), *monadgo.traitEither
	// ok:false, string
	// Left(ok:error), *monadgo.traitEither
	// ok:error, string
	// Right(true), *monadgo.traitEither
	// Right(Null), *monadgo.traitEither

}

func ExampleLeftProjection_FlatMap() {
	v := EitherOf(false).Left().FlatMap(func(x bool) Either {
		return RightOf(100, true)
	})
	printGet(v)
	printGet(v.Get())

	v = EitherOf(true).Left().FlatMap(func(x bool) Either {
		return Either1Of(100, true)
	})
	printGet(v)
	printGet(v.Get())

	// Output:
	// Right((100,true)), *monadgo.traitEither
	// (100,true), monadgo._tuple2
	// Right(true), *monadgo.traitEither
	// true, bool
}

func ExampleLeftProjection_ToOption() {
	v := EitherOf(false).Left().ToOption()
	printGet(v)

	v = EitherOf(fmt.Errorf("error")).Left().ToOption()
	printGet(v)

	v = EitherOf(nil).Left().ToOption()
	printGet(v)

	v = Either1Of(10, true).Left().ToOption()
	printGet(v)

	// Output:
	// Some(false), *monadgo.traitOption
	// Some(error), *monadgo.traitOption
	// None, *monadgo.traitOption
	// None, *monadgo.traitOption

}

func ExampleLeftProjection_Exists() {
	p := func(x int) bool {
		return x > 10
	}
	v := LeftOf(12).Left().Exists(p)
	printGet(v) // true

	v = LeftOf(7).Left().Exists(p)
	printGet(v) // false

	v = RightOf(12).Left().Exists(p)
	printGet(v) // false

	// Output:
	// true, bool
	// false, bool
	// false, bool
}

func ExampleLeftProjection_Filter() {
	p := func(x int) bool {
		return x > 10
	}

	v := LeftOf(12).Left().Filter(p)
	printGet(v) // Some(Left(12))

	v = LeftOf(7).Left().Filter(p)
	printGet(v) // None

	v = RightOf(12).Left().Filter(p)
	printGet(v) // None

	// Output:
	// Some(Left(12)), *monadgo.traitOption
	// None, *monadgo.traitOption
	// None, *monadgo.traitOption
}
