package monadgo

func ExampleTuple2() {
	t := Tuple2Of("2", 1)

	printGet(t)
	printGet(t.V1().(string))
	printGet(t.V(0).(string))
	printGet(t.V2().(int))
	printGet(t.V(1).(int))
	printGet(t.T1())
	printGet(t.T(0))
	printGet(t.T2())
	printGet(t.T(1))
	printGet(t.Get())
	printGet(t.reduce())
	printGet(t.toValues())
	printGet(t.Dimension())

	var x []int
	t = Tuple2Of(x, int64(100))

	printGet(t)
	printGet(t.V1())
	printGet(t.V(0))
	printGet(t.V2())
	printGet(t.V(1))
	printGet(t.T1())
	printGet(t.T(0))
	printGet(t.T2())
	printGet(t.T(1))

	// Output:
	// (2,1), monadgo.Tuple2
	// 2, string
	// 2, string
	// 1, int
	// 1, int
	// string, *reflect.rtype
	// string, *reflect.rtype
	// int, *reflect.rtype
	// int, *reflect.rtype
	// [2 1], [2]interface {}
	// (2,1), monadgo.Tuple2
	// [2 <int Value>], []reflect.Value
	// 2, int
	// ([],100), monadgo.Tuple2
	// [], []int
	// [], []int
	// 100, int64
	// 100, int64
	// []int, *reflect.rtype
	// []int, *reflect.rtype
	// int64, *reflect.rtype
	// int64, *reflect.rtype

}
