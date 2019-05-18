package monadgo

func ExampleTuple3() {
	t := Tuple3Of(1, "2", float64(3.14))

	printGet(t)
	printGet(t.V1().(int))
	printGet(t.V(0).(int))
	printGet(t.V2().(string))
	printGet(t.V(1).(string))
	printGet(t.V3().(float64))
	printGet(t.V(2).(float64))
	printGet(t.T1())
	printGet(t.T(0))
	printGet(t.T2())
	printGet(t.T(1))
	printGet(t.T3())
	printGet(t.T(2))
	printGet(t.Get())
	printGet(t.reduce())
	printGet(t.toValues())
	printGet(t.Dimension())

	// Output:
	// (1,2,3.14), monadgo.Tuple3
	// 1, int
	// 1, int
	// 2, string
	// 2, string
	// 3.14, float64
	// 3.14, float64
	// int, *reflect.rtype
	// int, *reflect.rtype
	// string, *reflect.rtype
	// string, *reflect.rtype
	// float64, *reflect.rtype
	// float64, *reflect.rtype
	// [1 2 3.14], [3]interface {}
	// (1,2), monadgo.Tuple2
	// [<int Value> 2 <float64 Value>], []reflect.Value
	// 3, int
}
