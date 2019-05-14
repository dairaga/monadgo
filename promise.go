package monadgo

import (
	"context"
)

// Promise represents scala-like Promise.
type Promise struct {
	Future
}

// Completed return true if the promise is completed.
func (p *Promise) Completed() bool {
	return p.Future != nil && p.Future.Completed()
}

// Complete completes the promise with result.
func (p *Promise) Complete(result Try) *Promise {
	if p.Completed() {
		return p
	}

	//p.Future = futureFromTry(context.Background(), result)

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

	p.Future = f
	return p
}

// ----------------------------------------------------------------------------

var unitPromise = (&Promise{}).Success(unit)

// EmptyPromise ...
func EmptyPromise(ctx context.Context) *Promise {
	return &Promise{
		Future: emptyFuture(ctx),
	}
}
