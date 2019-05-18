package monadgo

func ExamplePairOf() {
	p := PairOf(1, "100")

	printGet(p.Key().(int))
	printGet(p.Value().(string))

	// Ouput:
	// 1, int
	// 100, string
}
