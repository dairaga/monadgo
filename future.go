package monadgo

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"time"
)

var cancelFailure = FailureOf(fmt.Errorf("user cancel"))

// Future represents scala-like Future.
type Future interface {
	fmt.Stringer

	Completed() bool
	OnComplete(func(Try))

	Map(f interface{}) Future
	FlatMap(f interface{}) Future

	Recover(f interface{}) Future
	RecoverWith(f interface{}) Future

	Foreach(f interface{})
	Filter(f interface{}) Future
	Collect(pf PartialFunc) Future

	Value() Option
	Ready(atMost time.Duration) Option
	Result(atMost time.Duration) Option

	Cancel()
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

func (u *future) String() string {
	if u.completed {
		return fmt.Sprintf("Future(%v)", u.val)
	}

	return "Future(Not Yet)"
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

func (u *future) transform(f func(Try) Try) Future {
	e := newPromise(u.ctx, f)

	u.OnComplete(func(v Try) {
		e.Complete(v)
	})
	return e
}

func (u *future) transformWith(f func(Try) Future) Future {
	e := EmptyPromise(u.ctx)

	u.OnComplete(func(v Try) {
		n := f(v)
		if n != u {
			e.CompleteWith(n)
		} else {
			e.Complete(v)
		}
	})

	return e
}

func (u *future) Value() Option {
	if u.Completed() {
		return u.val.ToOption()
	}
	return None
}

func (u *future) Foreach(f interface{}) {
	u.Map(f)
}

func (u *future) Map(f interface{}) Future {
	if u.completed && u.val.Failed() {
		return u
	}
	ft := func(v Try) Try {
		if v.Failed() {
			return v
		}
		x := funcOf(f).invoke(v.Get())
		return tryFromX(x)
	}

	return u.transform(ft)
}

func (u *future) Recover(f interface{}) Future {
	if u.completed && u.val.OK() {
		return u
	}

	ft := func(v Try) Try {
		if v.OK() {
			return v
		}
		x := funcOf(f).invoke(v.Get())
		return tryFromX(x)
	}

	return u.transform(ft)
}

func (u *future) FlatMap(f interface{}) Future {
	if u.completed && u.val.Failed() {
		return u
	}

	ft := func(v Try) Future {
		if v.Failed() {
			return u
		}
		return funcOf(f).invoke(v.Get()).(Future)
	}

	return u.transformWith(ft)
}

func (u *future) RecoverWith(f interface{}) Future {
	if u.completed && u.val.OK() {
		return u
	}

	ft := func(v Try) Future {
		if v.OK() {
			return u
		}
		return funcOf(f).invoke(v.Get()).(Future)
	}

	return u.transformWith(ft)
}

func (u *future) Ready(atMost time.Duration) Option {
	if u.completed {
		return SomeOf(u)
	}

	ctx, cancel := context.WithTimeout(context.Background(), atMost)
	defer cancel()
	done := make(chan bool, 1)
	defer close(done)

	u.OnComplete(func(Try) {
		done <- true
	})

	select {
	case <-done:
		return SomeOf(u)
	case <-ctx.Done():
		return None
	}
}

func (u *future) Result(atMost time.Duration) Option {

	x := u.Ready(atMost)

	return x.FlatMap(func(f Future) Option {
		return f.Value()
	})
}

func (u *future) Filter(f interface{}) Future {
	return u.Map(f)
}

func (u *future) Collect(pf PartialFunc) Future {
	if u.completed && u.val.Failed() {
		return u
	}

	ft := func(v Try) Try {
		if v.Failed() {
			return v
		}

		result := pf.Call(reflect.ValueOf(v.Get()))
		if result == nothingValue {
			return FailureOf(false)
		}

		return SuccessOf(result.Interface())
	}

	return u.transform(ft)
}

func (u *future) Cancel() {
	u.cancel()
}

// ----------------------------------------------------------------------------

func initFuture(ctx context.Context) *future {
	ret := &future{mux: &sync.Mutex{}}
	ret.ctx, ret.cancel = context.WithCancel(ctx)
	ret.in = make(chan Try)
	return ret
}

func newFuture(ctx context.Context, f func(Try) Try) *future {
	ret := initFuture(ctx)

	go func() {
		defer func() {
			close(ret.in)
		}()

		select {
		case <-ret.ctx.Done():
			ret.val = cancelFailure
			ret.completed = true
			defer ret.mux.Unlock()
			ret.mux.Lock()
			ret.next = nil
			return
		case x, ok := <-ret.in:
			if !ok {
				return
			}
			if f != nil {
				ret.val = f(x)
			} else {
				ret.val = x
			}

			ret.completed = true

			defer ret.mux.Unlock()
			ret.mux.Lock()

			for _, callback := range ret.next {
				callback(ret.val)
			}
			ret.next = nil
		}
	}()

	return ret
}

// FutureOf returns a future.
func FutureOf(f interface{}) Future {
	return unitPromise.Map(f)
}

/*

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
*/
