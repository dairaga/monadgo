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

// Complete completes the promise with result.
func (p *Promise) Complete(result Try) *Promise {
	if p.Completed() {
		return p
	}

	p.future.in <- result

	return p
}

// Success ...
func (p *Promise) Success(v interface{}) *Promise {
	if p.Completed() {
		return p
	}

	return p.Complete(SuccessOf(v))
}

// Failure ...
func (p *Promise) Failure(v interface{}) *Promise {
	if p.Completed() {
		return p
	}

	if isErrorOrFalse(v) {
		return p.Complete(FailureOf(v))
	}
	return p.Success(v)
}

// CompleteWith completes the promise with future.
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

func newPromise(ctx context.Context, f func(Try) Try) *Promise {
	return &Promise{
		newFuture(ctx, f),
	}
}

// EmptyPromise ...
func EmptyPromise(ctx context.Context) *Promise {
	return newPromise(ctx, nil)
}

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
