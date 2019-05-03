package monadgo

import (
	"fmt"
)

func ExampleMapOf() {
	m := MapOf((map[string]int)(nil))
	fmt.Println(m)

	m = MapOf(map[string]int{
		"a": 11,
		"b": 22,
	})
	fmt.Println(m)
	printGet(m.Get())
	// Output:
	// map[]
	// map[a:11 b:22]
	// map[a:11 b:22], map[string]int
}

func ExampleMap_Foreach() {
	MapOf(map[string]int{
		"a": 11,
		"b": 22,
		"c": 33,
	}).Foreach(func(k string, v int) {
		fmt.Printf("%s->%d\n", k+k, v*v)
	})

	// Unordered output:
	// aa->121
	// bb->484
	// cc->1089
}

func ExampleMap_Forall() {
	a := MapOf(map[string]int{
		"a": 11,
		"b": 22,
		"c": 33,
	}).Forall(func(_ string, v int) bool {
		return v%11 == 0
	})
	fmt.Println(a)

	a = MapOf(map[string]int{
		"a": 11,
		"b": 22,
		"c": 33,
	}).Forall(func(_ string, v int) bool {
		return v < 20
	})
	fmt.Println(a)

	// Output:
	// true
	// false
}

func ExampleMap_Map() {
	s1 := MapOf(map[string][]int{
		"a": []int{11, 111},
		"b": []int{22, 222},
	}).Map(func(_ string, x []int) []int {
		return x
	}).Get().([][]int)

	for _, x := range s1 {
		for _, y := range x {
			fmt.Println(y)
		}
	}

	s2 := MapOf(map[string][]Pair{
		"a": []Pair{PairOf(1, 11), PairOf(1, 111)},
		"b": []Pair{PairOf(2, 22), PairOf(2, 222)},
	}).Map(func(_ string, x []Pair) []Pair {
		return x
	}).Get().([][]Pair)
	for _, x := range s2 {
		for _, y := range x {
			fmt.Println(y)
		}
	}

	s3 := MapOf(map[string]int{
		"a": 1,
		"b": 2,
	}).Map(func(k string, v int) (string, int) {
		return k + k, v + v
	}).Get().(map[string]int)

	for k, v := range s3 {
		fmt.Println(k, v)
	}

	s4 := MapOf(map[string]int{
		"a": 1,
		"b": 2,
	}).Map(func(k string, v int) Pair {
		return PairOf(k+k+k, v+v+v)
	}).Get().(map[string]int)

	for k, v := range s4 {
		fmt.Println(k, v)
	}

	// Unordered output:
	// 11
	// 111
	// 22
	// 222
	// (1,11)
	// (1,111)
	// (2,22)
	// (2,222)
	// aa 2
	// bb 4
	// aaa 3
	// bbb 6
}

func ExampleMap_FlatMap() {
	s1 := MapOf(map[string][]int{
		"a": []int{11, 111},
		"b": []int{22, 222},
	}).FlatMap(func(_ string, x []int) []int {
		return x
	})
	fmt.Println(s1)
	printGet(s1.Get())

	s2 := MapOf(map[string][]Pair{
		"a": []Pair{PairOf(1, 11), PairOf(1, 111)},
		"b": []Pair{PairOf(2, 22), PairOf(2, 222)},
	}).FlatMap(func(_ string, x []Pair) []Pair {
		return x
	})
	fmt.Println(s2)
	fmt.Printf("%v, %T\n", s2.Get(), s2.Get())

	s3 := MapOf(map[string]int{
		"a": 1,
		"b": 2,
	}).FlatMap(func(k string, v int) map[string]int {
		return map[string]int{
			k:         v,
			k + k:     v + v,
			k + k + k: v + v + v,
		}
	})
	fmt.Println(s3)
	fmt.Printf("%v, %T\n", s3.Get(), s3.Get())

	// Output:
	// [11 111 22 222]
	// [11 111 22 222], []int
	// map[1:111 2:222]
	// map[1:111 2:222], map[int]int
	// map[a:1 aa:2 aaa:3 b:2 bb:4 bbb:6]
	// map[a:1 aa:2 aaa:3 b:2 bb:4 bbb:6], map[string]int
}

func ExampleMap_ToSlice() {
	s := MapOf(map[string][]int{
		"a": []int{11, 111},
		"b": []int{22, 222},
	}).ToSlice().Get().([]Pair)

	for _, x := range s {
		printGet(x)
	}

	// Unordered Output:
	// (a,[11 111]), monadgo._pair
	// (b,[22 222]), monadgo._pair
}

func ExampleMap_Fold() {
	m := MapOf(map[int]int{
		1: 11,
		2: 22,
		3: 33,
		4: 44,
	})

	var z []Pair

	z = m.Fold(z, func(s []Pair, k, v int) []Pair {
		return append(s, PairOf(k+1, v+1))
	}).([]Pair)

	for _, x := range z {
		fmt.Println(x)
	}

	x := MapOf(map[int]int{
		1: 11,
		2: 22,
		3: 33,
		4: 44,
	}).Fold(PairOf(10, 100), func(p Pair, k, v int) Pair {
		return PairOf(p.Key().(int)+k, p.Value().(int)+v)
	})
	printGet(x)

	// Unordered Output:
	// (2,12)
	// (3,23)
	// (4,34)
	// (5,45)
	// (20,210), monadgo._pair
}

func ExampleMap_GroupBy() {
	m := MapOf(map[int]int{
		1: 11,
		2: 22,
		3: 33,
		4: 44,
	}).GroupBy(func(_, v int) string {
		return fmt.Sprintf("%d", v%2)
	}).Get().(map[string]map[int]int)

	for k1, v1 := range m {
		fmt.Println(k1)

		for k2, v2 := range v1 {
			fmt.Println(k2, v2)
		}
	}

	// Unordered Output:
	// 0
	// 2 22
	// 4 44
	// 1
	// 1 11
	// 3 33
}
