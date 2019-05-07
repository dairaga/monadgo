package monadgo

import (
	"fmt"
)

func ExampleSliceOf() {
	s := SliceOf([]int{1, 2, 3, 4, 5})
	fmt.Println(s)
	printGet(s.Get())
	fmt.Println(s.Len())
	fmt.Println(s.Cap())

	s = SliceOf([6]int{1, 2, 3, 4, 5, 6})
	fmt.Println(s)
	printGet(s.Get())
	fmt.Println(s.Len())
	fmt.Println(s.Cap())

	s = SliceOf(101)
	printGet(s.Get())
	fmt.Println(s.Len())
	fmt.Println(s.Cap())
	// Output:
	// [1 2 3 4 5]
	// [1 2 3 4 5], []int
	// 5
	// 5
	// [1 2 3 4 5 6]
	// [1 2 3 4 5 6], []int
	// 6
	// 6
	// [101], []int
	// 1
	// 1
}

func ExampleSlice_Foreach() {
	SliceOf([]int{1, 2, 3, 4, 5}).Foreach(func(x int) {
		fmt.Println(x)
	})

	SliceOf([5]int{1, 2, 3, 4, 5}).Foreach(func(x int) {
		fmt.Println(x)
	})

	// Output:
	// 1
	// 2
	// 3
	// 4
	// 5
	// 1
	// 2
	// 3
	// 4
	// 5
}

func ExampleSlice_Forall() {
	a := SliceOf([]int{1, 3, 5, 7, 9}).Forall(func(x int) bool {
		return x&1 == 1
	})
	fmt.Println(a)

	a = SliceOf([5]int{1, 3, 5, 7, 9}).Forall(func(x int) bool {
		return !(x&1 == 1)
	})
	fmt.Println(a)

	// Output:
	// true
	// false

}

func ExampleSlice_Fold() {
	sum := SliceOf([]int{1, 2, 3, 4, 5}).Fold(1, func(z, x int) int {
		return z + x
	})
	fmt.Println(sum)

	sum = SliceOf([5]int{1, 2, 3, 4, 5}).Fold(1, func(z, x int) int {
		return z + x
	})

	fmt.Println(sum)

	// Output:
	// 16
	// 16
}
func ExampleSlice_Map() {

	s1 := SliceOf([]Pair{PairOf(1, 11), PairOf(2, 22), PairOf(1, 111), PairOf(2, 222)})
	fmt.Println(s1)
	printGet(s1.Get())

	s2 := s1.Map(func(p Pair) string {
		return p.String()
	})

	fmt.Println(s2)
	fmt.Printf("%v, %T\n", s2.Get(), s2.Get())

	s2 = s1.Map(func(k, v int) (string, int) {
		return fmt.Sprintf("%d", k+v), k * v
	})

	fmt.Println(s2)
	fmt.Printf("%v, %T\n", s2.Get(), s2.Get())

	// Output:
	// [(1,11) (2,22) (1,111) (2,222)]
	// [(1,11) (2,22) (1,111) (2,222)], []monadgo.Pair
	// [(1,11) (2,22) (1,111) (2,222)]
	// [(1,11) (2,22) (1,111) (2,222)], []string
	// [(12,11) (24,44) (112,111) (224,444)]
	// [(12,11) (24,44) (112,111) (224,444)], []monadgo.Pair
}

func ExampleSlice_FlatMap() {
	s1 := SliceOf([]Pair{PairOf(1, 11), PairOf(2, 22), PairOf(1, 111), PairOf(2, 222)})
	fmt.Println(s1)
	printGet(s1.Get())

	s2 := s1.FlatMap(func(p Pair) []int {
		return []int{p.Key().(int), p.Value().(int)}
	})
	fmt.Println(s2)
	printGet(s2.Get())

	s2 = s1.FlatMap(func(p Pair) map[int]int {
		return map[int]int{
			p.Key().(int): p.Value().(int),
		}
	})
	fmt.Println(s2)
	printGet(s2.Get())

	SliceOf([]int{1, 2, 3}).FlatMap(func(x int) {
		SliceOf([]int{1, 2, 3}).Map(func(y int) {
			fmt.Printf("%dx%d=%d\n", x, y, x*y)
		})
	})

	// Output:
	// [(1,11) (2,22) (1,111) (2,222)]
	// [(1,11) (2,22) (1,111) (2,222)], []monadgo.Pair
	// [1 11 2 22 1 111 2 222]
	// [1 11 2 22 1 111 2 222], []int
	// [(1,11) (2,22) (1,111) (2,222)]
	// [(1,11) (2,22) (1,111) (2,222)], []monadgo.Pair
	// 1x1=1
	// 1x2=2
	// 1x3=3
	// 2x1=2
	// 2x2=4
	// 2x3=6
	// 3x1=3
	// 3x2=6
	// 3x3=9
}

func ExampleSlice_Tail() {
	s := SliceOf([]int{1, 2, 3, 4, 5}).Tail()
	printGet(s.Get())

	s = SliceOf([5]int{1, 2, 3, 4, 5}).Tail()
	printGet(s.Get())

	// Output:
	// [2 3 4 5], []int
	// [2 3 4 5], []int

}

func ExampleSlice_Reduce() {
	sum := SliceOf([]int{1, 2, 3, 4, 5}).Reduce(func(x1, x2 int) int {
		return x1 + x2
	})
	fmt.Println(sum)

	sum = SliceOf([5]int{1, 2, 3, 4, 5}).Reduce(func(x1, x2 int) int {
		return x1 * x2
	})

	fmt.Println(sum)

	// Output:
	// 15
	// 120
}

func ExampleSlice_Scan() {
	s := SliceOf([]int{1, 2, 3, 4, 5}).Scan(10, func(a, b int) int {
		return a * b
	})

	printGet(s.Get())

	// Output:
	// [10 10 20 60 240 1200], []int
}

func ExampleSlice_GroupBy() {
	m := SliceOf([]int{1, 2, 3, 4, 5}).GroupBy(func(x int) int {
		return x % 2
	}).Get().(map[int][]int)

	for k, v := range m {
		printGet(k)
		printGet(v)
	}

	// Unordered Output:
	// 1, int
	// [1 3 5], []int
	// 0, int
	// [2 4], []int
}

func ExampleSlice_Take() {
	s := SliceOf([]int{1, 2, 3, 4, 5}).Take(3)
	printGet(s.Get())

	// Output:
	// [1 2 3], []int
}

func ExampleSlice_TakeWhile() {
	s := SliceOf([]int{1, 2, 3, 4, 5}).TakeWhile(func(x int) bool {
		return x <= 3
	})
	printGet(s.Get())

	s = SliceOf([]int{1, 2, 3, 4, 5}).TakeWhile(func(x int) bool {
		return x > 0
	})
	printGet(s.Get())

	s = SliceOf([]int{1, 2, 3, 4, 5}).TakeWhile(func(x int) bool {
		return x > 10
	})
	printGet(s.Get())

	// Output:
	// [1 2 3], []int
	// [1 2 3 4 5], []int
	// [], []int
}

func ExampleSlice_Head() {
	s := SliceOf([]int{})
	printGet(s.Head())

	s = SliceOf([]int{1, 2, 3, 4, 5})
	printGet(s.Head())

	// Output:
	// <nil>, <nil>
	// 1, int
}

func ExampleSlice_HeadOption() {
	s := SliceOf([]int{})
	printGet(s.HeadOption())

	s = SliceOf([]int{1, 2, 3, 4, 5})
	printGet(s.HeadOption())

	// Output:
	// None, *monadgo.traitOption
	// Some(1), *monadgo.traitOption
}

func ExampleSlice_Drop() {
	s := SliceOf([]int{})
	printGet(s.Drop(10).Get())

	s = SliceOf([]int{1, 2, 3, 4, 5})
	printGet(s.Drop(2).Get())
	printGet(s.Drop(s.Len() + 1).Get())

	// Output:
	// [], []int
	// [3 4 5], []int
	// [], []int
}

func ExampleSlice_Exists() {
	s := SliceOf([]int{})
	printGet(s.Exists(func(x int) bool {
		return x > 3
	}))

	s = SliceOf([]int{1, 2, 3, 4, 5})
	printGet(s.Exists(func(x int) bool {
		return x > 3
	}))

	printGet(s.Exists(func(x int) bool {
		return x > 10
	}))

	// Output:
	// false, bool
	// true, bool
	// false, bool
}

func ExampleSlice_Find() {
	s := SliceOf([]int{})
	printGet(s.Find(func(x int) bool {
		return x > 3
	}))

	s = SliceOf([]int{1, 2, 3, 4, 5})
	printGet(s.Find(func(x int) bool {
		return x > 3
	}))

	printGet(s.Find(func(x int) bool {
		return x > 10
	}))

	// Output:
	// None, *monadgo.traitOption
	// Some(4), *monadgo.traitOption
	// None, *monadgo.traitOption
}

func ExampleSlice_Filter() {
	s := SliceOf([]int{})
	printGet(s.Filter(func(x int) bool {
		return x > 3
	}))

	s = SliceOf([]int{1, 2, 3, 4, 5})
	printGet(s.Filter(func(x int) bool {
		return x > 3
	}))

	printGet(s.Filter(func(x int) bool {
		return x > 10
	}))

	// Output:
	// [], monadgo.seq
	// [4 5], monadgo.seq
	// [], monadgo.seq
}

func ExampleSlice_IndexWhere() {
	s := SliceOf([]int{})
	printGet(s.IndexWhere(func(x int) bool {
		return x > 3
	}, 100))

	s = SliceOf([]int{1, 2, 3, 4, 5})
	printGet(s.IndexWhere(func(x int) bool {
		return x > 2
	}, 3))

	printGet(s.IndexWhere(func(x int) bool {
		return x > 2
	}, -1))

	printGet(s.IndexWhere(func(x int) bool {
		return x > 2
	}, 100))

	// Output:
	// -1, int
	// 3, int
	// 2, int
	// -1, int
}

func ExampleSlice_LastIndexWhere() {
	s := SliceOf([]int{})
	printGet(s.LastIndexWhere(func(x int) bool {
		return x > 3
	}, 100))

	s = SliceOf([]int{1, 2, 3, 4, 5})
	printGet(s.LastIndexWhere(func(x int) bool {
		return x > 2
	}, 3))

	printGet(s.LastIndexWhere(func(x int) bool {
		return x > 3
	}, 2))

	printGet(s.LastIndexWhere(func(x int) bool {
		return x > 2
	}, 100))

	// Output:
	// -1, int
	// 3, int
	// -1, int
	// 4, int
}

func ExampleSlice_MkString() {
	s := SliceOf([]int{})
	printGet(s.MkString("<", ":", ">"))

	s = SliceOf([]int{1, 2, 3, 4, 5})
	printGet(s.MkString("<", ":", ">"))

	// Output:
	// <>, string
	// <1:2:3:4:5:>, string
}

func ExampleSlice_Reverse() {
	s := SliceOf([]int{})
	printGet(s.Reverse().Get())

	s = SliceOf([]int{1, 2, 3, 4, 5})
	printGet(s.Reverse().Get())

	// Output:
	// [], []int
	// [5 4 3 2 1], []int
}

func ExampleSlice_Span() {
	s := SliceOf([]int{})
	printGet(s.Span(func(x int) bool {
		return x >= 2
	}))

	s = SliceOf([]int{1, 2, 3, 4, 5})
	printGet(s.Span(func(x int) bool {
		return x >= 2
	}))

	// Output:
	// ([],[]), monadgo._tuple2
	// ([1],[2 3 4 5]), monadgo._tuple2
}
