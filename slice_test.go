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

	// Output:
	// [1 2 3 4 5]
	// [1 2 3 4 5], []int
	// 5
	// 5
	// [1 2 3 4 5 6]
	// [1 2 3 4 5 6], [6]int
	// 6
	// 6
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

	// Output:
	// [(1,11) (2,22) (1,111) (2,222)]
	// [(1,11) (2,22) (1,111) (2,222)], []monadgo.Pair
	// [(1,11) (2,22) (1,111) (2,222)]
	// [(1,11) (2,22) (1,111) (2,222)], []string
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

	// Output:
	// [(1,11) (2,22) (1,111) (2,222)]
	// [(1,11) (2,22) (1,111) (2,222)], []monadgo.Pair
	// [1 11 2 22 1 111 2 222]
	// [1 11 2 22 1 111 2 222], []int
	// [(1,11) (2,22) (1,111) (2,222)]
	// [(1,11) (2,22) (1,111) (2,222)], []monadgo.Pair
}

func ExampleSlice_ToSeq() {
	s1 := SliceOf([]Pair{PairOf(1, 11), PairOf(2, 22), PairOf(1, 111), PairOf(2, 222)})
	fmt.Println(s1)
	fmt.Printf("%v, %T\n", s1.Get(), s1.Get())

	s2 := s1.ToSeq().([]Pair)
	fmt.Println(s2)
	fmt.Printf("%v, %T\n", s2, s2)

	// Output:
	// [(1,11) (2,22) (1,111) (2,222)]
	// [(1,11) (2,22) (1,111) (2,222)], []monadgo.Pair
	// [(1,11) (2,22) (1,111) (2,222)]
	// [(1,11) (2,22) (1,111) (2,222)], []monadgo.Pair
}
