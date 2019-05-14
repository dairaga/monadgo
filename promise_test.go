package monadgo

import (
	"context"
)

func ExampleEmptyPromise() {
	p := EmptyPromise(context.Background())
	printGet(p.Completed())
	p.OnComplete(func(v Try) {
		printGet(v)
	})

	p.Success(100)

	p = EmptyPromise(context.Background())
	printGet(p.Completed())
	p.OnComplete(func(v Try) {
		printGet(v)
	})

	p.Failure(false)
	sleep(1)

	// Unordered output:
	// false, bool
	// Success(100), *monadgo.traitTry
	// false, bool
	// Failure(false), *monadgo.traitTry

}
