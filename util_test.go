package monadgo

import "reflect"

func ExampleMakeMap() {
	m := makeMap(reflect.TypeOf(int(0)), reflect.TypeOf(int64(0)), -1).Interface()
	printGet(m)
	printGet(len(m.(map[int]int64)))

	m = makeMap(reflect.TypeOf(int(0)), reflect.TypeOf(float64(0.0)), 2).Interface()
	printGet(m)
	printGet(len(m.(map[int]float64)))

	// Output:
	// map[], map[int]int64
	// 0, int
	// map[], map[int]float64
	// 0, int

}

func ExampleMakeSlice() {
	s := makeSlice(reflect.TypeOf(int(0))).Interface()
	printGet(s)
	printGet(len(s.([]int)))
	printGet(cap(s.([]int)))

	s = makeSlice(reflect.TypeOf(int(0)), 1).Interface()
	printGet(s)
	printGet(len(s.([]int)))
	printGet(cap(s.([]int)))

	s = makeSlice(reflect.TypeOf(int(0)), 1, 2).Interface()
	printGet(s)
	printGet(len(s.([]int)))
	printGet(cap(s.([]int)))

	// Output:
	// [], []int
	// 0, int
	// 0, int
	// [0], []int
	// 1, int
	// 1, int
	// [0], []int
	// 1, int
	// 2, int

}

func ExampleAppendSlice() {
	x := reflect.Value{}
	y := reflect.ValueOf(10)
	z := appendSlice(x, y).Interface()
	printGet(z)

	x = reflect.ValueOf([]int{1})
	y = reflect.ValueOf(100)
	z = appendSlice(x, y).Interface()
	printGet(z)

	// Output:
	// [10], []int
	// [1 100], []int
}

func ExampleMergeSlice() {
	x := reflect.Value{}
	y := reflect.ValueOf(10)
	z := mergeSlice(x, y).Interface()
	printGet(z)

	x = reflect.ValueOf(100)
	y = reflect.ValueOf(101)
	z = mergeSlice(x, y).Interface()
	printGet(z)

	x = reflect.ValueOf([]int{100, 101})
	y = reflect.ValueOf(102)
	z = mergeSlice(x, y).Interface()
	printGet(z)

	x = reflect.ValueOf([]int{100, 101})
	y = reflect.ValueOf([]int{102, 103})
	z = mergeSlice(x, y).Interface()
	printGet(z)

	// Output:
	// 10, int
	// [100 101], []int
	// [100 101 102], []int
	// [100 101 102 103], []int

}

func ExampleMergeMap() {
	x := reflect.Value{}
	y := reflect.ValueOf(PairOf(1, 2))
	z := mergeMap(x, y).Interface()
	printGet(z)

	x = reflect.ValueOf(PairOf(1, 2))
	y = reflect.ValueOf(PairOf(3, 4))
	z = mergeMap(x, y).Interface()
	printGet(z)

	x = reflect.ValueOf(map[int]int{1: 11})
	y = reflect.ValueOf(PairOf(2, 22))
	z = mergeMap(x, y).Interface()
	printGet(z)

	x = reflect.ValueOf(map[int]int{1: 11, 2: 22})
	y = reflect.ValueOf(map[int]int{3: 33, 4: 44})
	z = mergeMap(x, y).Interface()
	printGet(z)

	// Output:
	// map[1:2], map[int]int
	// map[1:2 3:4], map[int]int
	// map[1:11 2:22], map[int]int
	// map[1:11 2:22 3:33 4:44], map[int]int

}

func ExampleMergeKeyValue() {
	x := reflect.Value{}
	z := mergeKeyValue(x, reflect.ValueOf(1), reflect.ValueOf(11)).Interface()
	printGet(z)

	x = reflect.ValueOf(map[int]int{1: 11})
	z = mergeKeyValue(x, reflect.ValueOf(1), reflect.ValueOf(22)).Interface()
	printGet(z)

	// Output:
	// map[1:11], map[int]int
	// map[1:22], map[int]int

}
