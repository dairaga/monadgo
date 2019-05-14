package monadgo

import (
	"fmt"
	"log"
	"testing"
	"time"
)

func sleep(n int) {
	time.Sleep(time.Duration(n) * time.Second)
}
func TestFuture(t *testing.T) {

	f1 := FutureOf(func() int {
		//sleep(2)
		return 1000
	})

	f2 := FutureOf(func() int {
		//sleep(2)
		return 2000
	})

	f3 := f1.FlatMap(func(x int) Future {
		return f2.Map(func(y int) int {
			//sleep(2)
			return x * y
		})
	})

	printGet(f1.Result(5 * time.Second))
	printGet(f2.Result(5 * time.Second))
	printGet(f3.Result(10 * time.Second))
}

func TestFuture2(t *testing.T) {

	f1 := FutureOf(func() int {
		log.Println("f1 called and wait")
		sleep(3)
		return 1000
	})

	f2 := f1.Map(func(x int) string {
		log.Println("f2 called")
		return fmt.Sprintf("%d", x*x)
	})

	printGet(f1.Result(5 * time.Second))
	printGet(f2.Result(5 * time.Second))

	//printGet(f1.Result(5 * time.Second))
	//printGet(f2.Result(5 * time.Second))
	//printGet(f3.Result(10 * time.Second))
}
