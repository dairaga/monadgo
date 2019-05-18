package monadgo

func ExampleTuple() {
	t := TupleOf([]interface{}{
		1, "2", 3.14, complex(5, 7), true,
	})

	printGet(t)
	printGet(t.V(0).(int))
	printGet(t.V(1).(string))
	printGet(t.V(2).(float64))
	printGet(t.V(3).(complex128))
	printGet(t.V(4).(bool))
	printGet(t.T(0))
	printGet(t.T(1))
	printGet(t.T(2))
	printGet(t.T(3))
	printGet(t.T(4))

	printGet(t.Get())
	printGet(t.reduce())
	printGet(t.toValues())
	printGet(t.Dimension())

	// Output:
	// (1,2,3.14,(5+7i),true), *monadgo.TupleN
	// 1, int
	// 2, string
	// 3.14, float64
	// (5+7i), complex128
	// true, bool
	// int, *reflect.rtype
	// string, *reflect.rtype
	// float64, *reflect.rtype
	// complex128, *reflect.rtype
	// bool, *reflect.rtype
	// [1 2 3.14 (5+7i) true], []interface {}
	// (1,2,3.14,(5+7i)), monadgo.Tuple4
	// [<int Value> 2 <float64 Value> <complex128 Value> <bool Value>], []reflect.Value
	// 5, int
}
