package monadgo

import (
	"context"
	"fmt"
	"sync"
)

// Promise represents scala-like Promise.
type Promise struct {
	*future
}

func (p *Promise) String() string {
	if p.Completed() {
		return fmt.Sprintf("Promise(%v)", p.future.val)
	}

	return "Promise(Not Yet)"
}

// Completed return true if the promise is completed.
func (p *Promise) Completed() bool {
	return p.future != nil && p.future.Completed()
}

// Complete completes the promise with try result.
// Have no effect on p if p is completed.
func (p *Promise) Complete(result Try) *Promise {
	if p.Completed() {
		return p
	}

	p.future.in <- result

	return p
}

// Success completes the promise with a Success value v.
// Have no effect on p if p is completed.
func (p *Promise) Success(v interface{}) *Promise {
	if p.Completed() {
		return p
	}

	return p.Complete(SuccessOf(v))
}

// Failure completes the promise p with a Failure value v.
// v must be error or false, otherwise p is completed with a Success v.
// Have no effect on p if p is completed.
func (p *Promise) Failure(v interface{}) *Promise {
	if p.Completed() {
		return p
	}

	if yes, _ := isErrorOrFalse(v); yes {
		return p.Complete(FailureOf(v))
	}
	return p.Success(v)
}

// CompleteWith completes the promise with future.
// Have no effect on p if p is completed.
func (p *Promise) CompleteWith(f Future) *Promise {
	if p.Completed() {
		return p
	}

	f.OnComplete(func(v Try) {
		p.future.in <- v
	})

	return p
}

// ----------------------------------------------------------------------------
// newPromise returns a new Promise from parent's context and assign a function to finish.
func newPromise(ctx context.Context, f func(Try) Try) *Promise {
	return &Promise{
		newFuture(ctx, f),
	}
}

// DefaultPromise returns a new Promise from parent's context,
// and wait for inputing value to complete the promise.
func DefaultPromise(ctx context.Context) *Promise {
	return newPromise(ctx, nil)
}

// unitPromise is a Promise with unit value.
// It is a root promise.
var unitPromise *Promise

func init() {
	unitPromise = &Promise{
		&future{
			completed: true,
			val:       SuccessOf(unit),
			mux:       &sync.Mutex{},
		},
	}

	unitPromise.ctx, unitPromise.cancel = context.WithCancel(context.Background())
}
