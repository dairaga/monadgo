package monadgo

import "fmt"

func printGet(x interface{}) {
	fmt.Printf("%v, %T\n", x, x)
}
