package monadgo

func ExampleTuple4() {
	t := Tuple4Of(1, "2", float64(3.14), complex(5, 7))

	printGet(t)
	printGet(t.V1())
	printGet(t.V(0))

	printGet(t.V2())
	printGet(t.V(1))

	printGet(t.V3())
	printGet(t.V(2))

	printGet(t.V4())
	printGet(t.V(3))

	printGet(t.T1())
	printGet(t.T(0))

	printGet(t.T2())
	printGet(t.T(1))

	printGet(t.T3())
	printGet(t.T(2))

	printGet(t.T4())
	printGet(t.T(3))

	printGet(t.Get())
	printGet(t.reduce())
	printGet(t.toValues())
	printGet(t.Dimension())

	// Output:
	// (1,2,3.14,(5+7i)), monadgo._tuple4
	// 1, int
	// 1, int
	// 2, string
	// 2, string
	// 3.14, float64
	// 3.14, float64
	// (5+7i), complex128
	// (5+7i), complex128
	// int, *reflect.rtype
	// int, *reflect.rtype
	// string, *reflect.rtype
	// string, *reflect.rtype
	// float64, *reflect.rtype
	// float64, *reflect.rtype
	// complex128, *reflect.rtype
	// complex128, *reflect.rtype
	// [1 2 3.14 (5+7i)], [4]interface {}
	// (1,2,3.14), monadgo._tuple3
	// [<int Value> 2 <float64 Value> <complex128 Value>], []reflect.Value
	// 4, int
}
