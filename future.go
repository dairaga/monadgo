package monadgo

import (
	"context"
	"sync"
)

// Future represents scala-like Future.
type Future interface {
	Completed() bool
	OnComplete(func(Try))
	Value() Option
	Foreach(f interface{})
	Map(f interface{}) Future
	FlatMap(f interface{}) Future
}

// ----------------------------------------------------------------------------

type future struct {
	completed bool
	in        chan Try
	ctx       context.Context
	cancel    context.CancelFunc
	val       Try
	mux       *sync.Mutex
	next      []func(Try)
}

var _ Future = &future{}

func (u *future) register(f func(Try)) {
	defer u.mux.Unlock()
	u.mux.Lock()
	u.next = append(u.next, f)
}

func (u *future) transform(f func(Try) Try) Future {
	e := EmptyPromise(u.ctx)

	u.OnComplete(func(v Try) {
		e.Complete(f(v))
	})
	return e
}

func (u *future) transformWith(f func(Try) Future) Future {
	e := EmptyPromise(u.ctx)

	u.OnComplete(func(v Try) {
		e.CompleteWith(f(v))
	})

	return e
}

func (u *future) Completed() bool {
	return u.completed
}

func (u *future) OnComplete(f func(Try)) {
	if u.completed {
		f(u.val)
		return
	}
	u.register(f)
}

func (u *future) Value() Option {
	if u.Completed() {
		return SomeOf(u.val)
	}
	return None
}

func (u *future) Foreach(f interface{}) {

}

func (u *future) Map(f interface{}) Future {
	ft := func(v Try) Try {
		if v.OK() {
			x := funcOf(f).invoke(v.Get())
			return tryFromX(x)
		}
		return v
	}

	return u.transform(ft)
}

func (u *future) FlatMap(f interface{}) Future {
	return nil
}

// ----------------------------------------------------------------------------

func emptyFuture(ctx context.Context) *future {
	ret := &future{mux: &sync.Mutex{}}
	ret.ctx, ret.cancel = context.WithCancel(ctx)
	return ret
}

func futureFromTry(ctx context.Context, result Try) *future {
	ret := emptyFuture(ctx)
	ret.val = result
	ret.completed = true
	return ret
}

// FutureOf returns a future.
func FutureOf(f interface{}) Future {
	return unitPromise.Map(f)
}
