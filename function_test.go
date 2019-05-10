package monadgo

import "fmt"

func func0() {
	fmt.Println("test")
}
func func00() fmt.Stringer {
	return nil
}

func func01() int {
	return 1
}
func func02() (int, string) {
	return 2, "AB"
}
func func03() (int, string, float64) {
	return 3, "ABC", 3.33
}

func func04() (int, string, float64, complex128) {
	return 4, "ABCD", 4.444, complex(4, 4)
}

func func05() (int, string, float64, complex128, byte) {
	return 5, "ABCDE", 5.555, complex(5, 5), '5'
}

func ExampleFuncOf0() {

	f := funcOf(func0)
	v := f.invoke(unit)
	printGet(v)

	f = funcOf(func00)
	v = f.invoke(unit)
	printGet(v)

	f = funcOf(func01)
	v = f.invoke(unit)
	printGet(v)

	f = funcOf(func02)
	v = f.invoke(unit)
	printGet(v)

	f = funcOf(func03)
	v = f.invoke(unit)
	printGet(v)

	f = funcOf(func04)
	v = f.invoke(unit)
	printGet(v)

	f = funcOf(func05)
	v = f.invoke(unit)
	printGet(v)

	// Output:
	// test
	// Void, monadgo._unit
	// Null, *monadgo._null
	// 1, int
	// (2,AB), monadgo.Pair
	// (3,ABC,3.33), monadgo.Tuple3
	// (4,ABCD,4.444,(4+4i)), monadgo.Tuple4
	// (5,ABCDE,5.555,(5+5i),53), monadgo.TupleN
}

// ----------------------------------------------------------------------------

func func1(int) {
	fmt.Println("test")
}
func func11(string) int {
	return 1
}
func func12(float64) (int, string) {
	return 2, "AB"
}
func func13(int64) (int, string, float64) {
	return 3, "ABC", 3.33
}

func func14(bool) (int, string, float64, complex128) {
	return 4, "ABCD", 4.444, complex(4, 4)
}

func func15(Tuple) (int, string, float64, complex128, byte) {
	return 5, "ABCDE", 5.555, complex(5, 5), '5'
}

func ExampleFuncOf1() {

	f := funcOf(func1)
	v := f.invoke(1)
	printGet(v)

	f = funcOf(func11)
	v = f.invoke("")
	printGet(v)

	f = funcOf(func12)
	v = f.invoke(float64(1))
	printGet(v)

	f = funcOf(func13)
	v = f.invoke(int64(1))
	printGet(v)

	f = funcOf(func14)
	v = f.invoke(true)
	printGet(v)

	f = funcOf(func15)
	v = f.invoke(Tuple2Of(1, 2))
	printGet(v)

	// Output:
	// test
	// Void, monadgo._unit
	// 1, int
	// (2,AB), monadgo.Pair
	// (3,ABC,3.33), monadgo.Tuple3
	// (4,ABCD,4.444,(4+4i)), monadgo.Tuple4
	// (5,ABCDE,5.555,(5+5i),53), monadgo.TupleN
}

// ----------------------------------------------------------------------------

func func2(int, int) {
	fmt.Println("test")
}
func func21(string, string) int {
	return 1
}
func func22(float64, float64) (int, string) {
	return 2, "AB"
}
func func23(int64, int64) (int, string, float64) {
	return 3, "ABC", 3.33
}

func func24(bool, bool) (int, string, float64, complex128) {
	return 4, "ABCD", 4.444, complex(4, 4)
}

func func25(Tuple, Tuple) (int, string, float64, complex128, byte) {
	return 5, "ABCDE", 5.555, complex(5, 5), '5'
}

func ExampleFuncOf2() {

	f := funcOf(func2)
	v := f.invoke(Tuple2Of(1, 1))
	printGet(v)

	f = funcOf(func21)
	v = f.invoke(Tuple2Of("", ""))
	printGet(v)

	f = funcOf(func22)
	v = f.invoke(Tuple2Of(float64(1), float64(1)))
	printGet(v)

	f = funcOf(func23)
	v = f.invoke(Tuple2Of(int64(1), int64(1)))
	printGet(v)

	f = funcOf(func24)
	v = f.invoke(Tuple2Of(true, true))
	printGet(v)

	f = funcOf(func25)
	v = f.invoke(Tuple2Of(Tuple2Of(1, 2), Tuple2Of(1, 2)))
	printGet(v)

	// Output:
	// test
	// Void, monadgo._unit
	// 1, int
	// (2,AB), monadgo.Pair
	// (3,ABC,3.33), monadgo.Tuple3
	// (4,ABCD,4.444,(4+4i)), monadgo.Tuple4
	// (5,ABCDE,5.555,(5+5i),53), monadgo.TupleN
}

// ----------------------------------------------------------------------------

func fold1(z int) int {
	return z + 1
}

func fold2(z []Pair, x Pair) []Pair {
	return append(z, x)
}

func fold3(z []Pair, k, v int) []Pair {
	return append(z, PairOf(k+1, v+1))
}

func fold4(k1, v1, k2, v2 int) []Pair {
	return []Pair{PairOf(k1, v1), PairOf(k2, v2)}
}

func ExampleFoldOf() {

	fw := foldOf(fold1)
	v := fw.fold(10, unit)
	printGet(v)

	fw = foldOf(fold2)
	v = fw.fold([]Pair{}, PairOf(1, 2))
	printGet(v)

	fw = foldOf(fold3)
	v = fw.fold([]Pair{}, PairOf(2, 4))
	printGet(v)

	fw = foldOf(fold4)
	v = fw.fold(PairOf(1, 2), PairOf(3, 4))
	printGet(v)

	// Output:
	// 11, int
	// [(1,2)], []monadgo.Pair
	// [(3,5)], []monadgo.Pair
	// [(1,2) (3,4)], []monadgo.Pair
}
