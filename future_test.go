package monadgo

import (
	"fmt"
	"testing"
	"time"
)

func sleep(n int) {
	time.Sleep(time.Duration(n) * time.Second)
}

const wait = 5 * time.Second

func TestFuture(t *testing.T) {
	a := 1000
	b := 2000

	f1 := FutureOf(func() int {
		return a
	})

	f2 := FutureOf(func() int {
		return b
	})

	f3 := f1.FlatMap(func(x int) Future {
		return f2.Map(func(y int) int {
			return x * y
		})
	})

	if f1.Result(wait).Get().(int) != 1000 {
		t.Errorf("f1 failure")
	}

	if f2.Result(wait).Get().(int) != 2000 {
		t.Errorf("f2 failure")
	}

	if f3.Result(wait).Get().(int) != a*b {
		t.Errorf("f3 failure")
	}
}

func TestFuture2(t *testing.T) {

	a := 1000
	f1 := FutureOf(func() int {
		sleep(3)
		return a
	})

	f2 := f1.Map(func(x int) string {
		return fmt.Sprintf("%d", x*x)
	})

	if f1.Result(wait).Get().(int) != a {
		t.Errorf("f1 failure")
	}

	if f2.Result(wait).Get().(string) != fmt.Sprintf("%d", a*a) {
		t.Errorf("f2 failure")
	}
}

func TestFuture_Cancel(t *testing.T) {
	f1 := FutureOf(func() int {
		sleep(10)
		return 10
	})

	f2 := f1.Map(func(x int) string {
		return fmt.Sprintf("%d", x*x)
	})

	f1.Cancel()
	if f2.Result(wait).Defined() {
		t.Errorf("f2 should be None")
	}

	if !f2.Completed() {
		t.Errorf("f2 should be comleted even it is canceled")
	}
}

func ExampleFuture_Foreach() {
	f1 := FutureOf(func() (int, bool) {
		return 100, true
	})

	f1.Foreach(func(x int) {
		fmt.Println(x)
	})

	printGet(f1.Result(wait).Get())
	sleep(1)

	f1 = FutureOf(func() (int, bool) {
		return 0, false
	})
	f1.Foreach(func(x int) {
		fmt.Println(x)
	})
	printGet(f1.Result(wait).Get())
	sleep(1)

	// Unordered output:
	// 100
	// 100, int
	// Nothing, *monadgo._nothing

}

func TestFuture_Filter(t *testing.T) {
	f1 := FutureOf(func() (int, bool) {
		return 100, true
	})

	f2 := f1.Filter(func(x int) bool {
		return x > 101
	})

	flag := true
	f2.OnComplete(func(t Try) {
		flag = t.OK()
	})

	if f2.Result(wait).Defined() {
		t.Errorf("f2 should be None")
	}

	if flag {
		t.Errorf("f2 value should be Failure")
	}
	//t.Log(f2)

}

func TestFuture_Collect(t *testing.T) {
	a := 100
	b := 200
	f1 := FutureOf(func() (int, int, bool) {
		sleep(3)
		return a, b, true
	})

	f2 := f1.Collect(PartialFuncOf(
		func(x, y int) bool {
			return x+y >= 100
		},
		func(x, y int) (string, string) {

			return fmt.Sprintf("%d", x+y), fmt.Sprintf("%d", x*y)
		},
	))

	result := f2.Result(wait).Get().(Tuple)
	if result.V(0).(string) != fmt.Sprintf("%d", a+b) || result.V(1).(string) != fmt.Sprintf("%d", a*b) {
		t.Errorf("result of f2 not match %d, %d", a+b, a*b)
	}

	f3 := f1.Collect(PartialFuncOf(
		func(x, y int) bool {
			return x+y < 100
		},
		func(x, y int) (string, string) {

			return fmt.Sprintf("%d", x+y), fmt.Sprintf("%d", x*y)
		},
	))

	r3 := f3.Result(wait)
	if !f3.Completed() {
		t.Errorf("f3 should be comleted")
	}

	if r3.Defined() {
		t.Errorf("r3 should be None")
	}

}

func TestFuture_Result(t *testing.T) {
	f1 := FutureOf(func() int {
		sleep(10)
		return 10
	})

	r := f1.Result(1 * time.Second)
	if r.Defined() || f1.Completed() {
		t.Error("r must be None and f1 must not be completed")
	}
}

func TestFuture_Recover(t *testing.T) {
	f1 := FutureOf(func() bool {
		return false
	})

	f2 := f1.Recover(func(x bool) string {
		return "false -> true"
	})

	result := f2.Result(wait)
	if !result.Defined() || result.Get().(string) != "false -> true" {
		t.Errorf("recover failture")
	}
}

func TestFuture_RecoverWith(t *testing.T) {
	f1 := FutureOf(func() bool {
		return false
	})

	f2 := f1.RecoverWith(func(x bool) Future {
		return FutureOf(func() string {
			return "false -> true"
		})
	})

	result := f2.Result(wait)
	if !result.Defined() || result.Get().(string) != "false -> true" {
		t.Errorf("recover failture")
	}

}
