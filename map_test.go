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

	m = MapOf([]Pair{PairOf("a", 11), PairOf("b", 22)})
	printGet(m.Get())

	m = MapOf(PairOf("a", 11))
	printGet(m.Get())

	m = MapOf(seqOf([]Pair{PairOf("a", 11), PairOf("b", 22)}))
	printGet(m.Get())

	// Output:
	// map[]
	// map[a:11 b:22]
	// map[a:11 b:22], map[string]int
	// map[a:11 b:22], map[string]int
	// map[a:11], map[string]int
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
	}).Get().([]int)
	for _, x := range s1 {
		fmt.Println(x)
	}

	s2 := MapOf(map[string][]Pair{
		"a": []Pair{PairOf(1, 11), PairOf(1, 111)},
		"b": []Pair{PairOf(2, 22), PairOf(2, 222)},
	}).FlatMap(func(_ string, x []Pair) []Pair {
		return x
	}).Get().(map[int]int)

	for k, v := range s2 {
		fmt.Println(k, v)
	}

	s3 := MapOf(map[string]int{
		"a": 1,
		"b": 2,
	}).FlatMap(func(k string, v int) map[string]int {
		return map[string]int{
			k:         v,
			k + k:     v + v,
			k + k + k: v + v + v,
		}
	}).Get().(map[string]int)

	for k, v := range s3 {
		fmt.Println(k, v)
	}

	// Unordered output:
	// 22
	// 222
	// 11
	// 111
	// 1 111
	// 2 222
	// a 1
	// aa 2
	// aaa 3
	// b 2
	// bb 4
	// bbb 6
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
	// (20,210), monadgo.Pair
}

func ExampleMap_Reduce() {
	m := MapOf(map[int]int{
		1: 11,
		2: 22,
		3: 33,
		4: 44,
	})

	printGet(m.Reduce(func(p1, p2 Pair) Pair {
		return PairOf(
			p1.Key().(int)+p2.Key().(int),
			p1.Value().(int)+p2.Value().(int),
		)
	}))

	printGet(m.Reduce(func(k1, v1, k2, v2 int) Pair {
		return PairOf(k1+k2, v1+v2)
	}))

	// Output:
	// (10,110), monadgo.Pair
	// (10,110), monadgo.Pair
}

func ExampleMap_GroupBy() {
	m := MapOf(map[int]int{
		1: 11,
		2: 22,
		3: 33,
		4: 44,
	})

	m2 := m.GroupBy(func(k, v int) int {
		return v % 2
	}).Get().(map[int]map[int]int)

	for k, v := range m2 {
		fmt.Println(k)

		for k1, v1 := range v {
			fmt.Println(k1, v1)
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

func ExampleMap_Exists() {
	m := MapOf(map[int]int{
		1: 11,
		2: 22,
		3: 33,
		4: 44,
	})

	printGet(m.Exists(func(k, v int) bool {
		return (k+v)&1 == 1
	}))

	printGet(m.Exists(func(k, v int) bool {
		return !((k+v)&1 == 1)
	}))

	// Output:
	// false, bool
	// true, bool
}

func ExampleMap_Find() {
	m := MapOf(map[int]int{
		1: 11,
		2: 22,
		3: 33,
		4: 44,
	})

	printGet(m.Find(func(k, v int) bool {
		return (k+v)&1 == 1
	}))

	printGet(m.Find(func(k, v int) bool {
		return !((k+v)&1 == 1)
	}).Defined())

	// Output:
	// None, *monadgo.traitOption
	// true, bool
}

func ExampleMap_Filter() {
	m := MapOf(map[int]int{
		1: 11,
		2: 22,
		3: 33,
		4: 44,
	})

	m1 := m.Filter(func(k, v int) bool {
		return (k+v)&1 == 1
	}).Get().(map[int]int)

	for k, v := range m1 {
		fmt.Println(k, v)
	}

	m1 = m.Filter(func(k, v int) bool {
		return !((k+v)&1 == 1)
	}).Get().(map[int]int)

	for k, v := range m1 {
		fmt.Println(k, v)
	}

	// Unordered Output:
	// 1 11
	// 2 22
	// 3 33
	// 4 44
}

func ExampleMap_Split() {
	m := MapOf(map[int]int{
		1: 11,
		2: 22,
		3: 33,
		4: 44,
	})

	t := m.Split(func(k, v int) bool {
		return (k+v)&1 == 1
	})

	fmt.Println("v1")
	for k, v := range t.V1().(map[int]int) {
		fmt.Println(k, v)
	}
	fmt.Println("v2")
	for k, v := range t.V2().(map[int]int) {
		fmt.Println(k, v)
	}

	// Unordered Output:
	// v1
	// 1 11
	// 2 22
	// 3 33
	// 4 44
	// v2
}

func ExampleMap_MkString() {
	m := MapOf(map[int]int{
		1: 11,
		2: 22,
		3: 33,
		4: 44,
	})

	fmt.Println(m.MkString("", "\n", "\n"))
	// Unordered Output:
	// (1,11)
	// (2,22)
	// (3,33)
	// (4,44)
}

func ExampleMap_Collect() {
	m := MapOf(map[int]int{
		1: 11,
		2: 22,
		3: 33,
		4: 44,
	})
	pf := PartialFuncOf(func(k, v int) bool {
		return (k+v)&1 == 1
	}, func(p Pair) Pair {
		return p
	})

	m1 := m.Collect(pf).Get().(map[int]int)

	for k, v := range m1 {
		fmt.Println(k, v)
	}

	pf = PartialFuncOf(func(k, v int) bool {
		return !((k+v)&1 == 1)
	}, func(k, v int) int {
		return k + v
	})

	s1 := m.Collect(pf).Get().([]int)

	for _, v := range s1 {
		fmt.Println(v)
	}

	pf = PartialFuncOf(func(k, v int) bool {
		return !((k+v)&1 == 1)
	}, func(p Pair) Pair {
		return p
	})

	m1 = m.Collect(pf).Get().(map[int]int)

	for k, v := range m1 {
		fmt.Println(k, v)
	}

	// Unordered Output:
	// 12
	// 24
	// 36
	// 48
	// 1 11
	// 2 22
	// 3 33
	// 4 44
}
