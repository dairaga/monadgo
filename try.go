package monadgo

import (
	"fmt"
	"reflect"
)

// TryOrElse is a function type for Try.OrElse.
type TryOrElse func() Try

// Try represents scala-like Try.
type Try interface {
	Any

	OK() bool
	Failed() bool

	// Traversable methods start
	Forall(f interface{}) bool
	Foreach(f interface{})
	Fold(z, f interface{}) interface{}
	// Traversable methods end

	Map(f interface{}) Try
	FlatMap(f interface{}) Try

	OrElse(z TryOrElse) Try
	GetOrElse(z interface{}) interface{}

	ToOption() Option
}

type traitTry struct {
	ok bool
	container
}

func (t *traitTry) String() string {
	if t.ok {
		return fmt.Sprintf("Success(%v)", t.Get())
	}
	return fmt.Sprintf("Failure(%v)", t.Get())
}

func (t *traitTry) OK() bool {
	return t.ok
}

func (t *traitTry) Failed() bool {
	return !t.ok
}

func (t *traitTry) Forall(f interface{}) bool {
	if !t.ok {
		return false
	}

	return t.container.Forall(f)
}

func (t *traitTry) Foreach(f interface{}) {
	if t.ok {
		t.container.Foreach(f)
	}
}

func (t *traitTry) Fold(z, f interface{}) interface{} {

	ztyp := reflect.TypeOf(z)

	if !t.ok {
		if ztyp.Kind() == reflect.Func {
			return t.invoke(z).Interface()
		}
		return z
	}

	ret := t.container.invoke(f).Interface()

	var x interface{}
	v, ok := ret.(Tuple)
	if ok {
		x = v.V(v.Dimension() - 1)
	} else {
		x = ret
	}

	if !isErrorOrFalse(x) {
		return ret
	}

	return FailureOf(x).Fold(z, nil)
}

func (t *traitTry) Map(f interface{}) Try {
	if t.ok {
		return tryFromContainer(t._map(f), true)
	}

	return t
}

func (t *traitTry) FlatMap(f interface{}) Try {
	if t.ok {
		return t._flatMap(f).(Try)

	}
	return t
}

func (t *traitTry) OrElse(z TryOrElse) Try {
	if !t.ok {
		return z()
	}

	return t
}

func (t *traitTry) GetOrElse(z interface{}) interface{} {
	if !t.ok {
		return checkAndInvoke(z)
	}

	return t.Get()
}

func (t *traitTry) ToOption() Option {
	if !t.ok {
		return None
	}

	return OptionOf(t.container)
}

func tryFromContainer(c container, ok bool) Try {
	return &traitTry{
		ok:        ok,
		container: c,
	}
}

/*func tryFromX(x interface{}) Try {
	return tryFromContainer(containerOf(x), !isErrorOrFalse(x))
}*/

// isErrorOrFalse checks x is an existing error or false.
func isErrorOrFalse(x interface{}) bool {
	switch v := x.(type) {
	case error:
		return v != nil
	case bool:
		return !v
	default:
		return false
	}
}

// ----------------------------------------------------------------------------

// FailureOf returns Failure of x if x is error or false, or returns Success of x.
func FailureOf(x interface{}) Try {
	return tryFromContainer(containerOf(x), false)
}

// ----------------------------------------------------------------------------

// SuccessOf returns Success of x.
func SuccessOf(x interface{}) Try {
	return tryFromContainer(containerOf(x), true)
}

/*
// tryFromTuple return a Try from Tuple.
func tryFromTuple(t Tuple) Try {
	x := t.V(t.Dimension() - 1)
	if isErrorOrFalse(x) {
		return FailureOf(x)
	}

	reduce := false
	switch x.(type) {
	case error, bool:
		reduce = true
	}

	if t.Dimension() > 2 {
		if reduce {
			return SuccessOf(t.reduce())
		}
		return SuccessOf(t)
	}

	return SuccessOf(t.V(0))
}
*/

// TryOf returns a Try with one parameter.
// errOrFalse must be bool or error type.
// Return Failure if errOrFalse is false or error existing,
// or Success of Null.
func TryOf(errOrFalse interface{}) Try {
	if isErrorOrFalse(errOrFalse) {
		return FailureOf(errOrFalse)
	}
	return SuccessOf(errOrFalse)
}

// Try1Of returns a Try.
// errOrFalse must be bool or error type.
// Return Failure if errOrFalse is false or error existing,
// or Success of x.
func Try1Of(x, errOrFalse interface{}) Try {
	if isErrorOrFalse(errOrFalse) {
		return FailureOf(errOrFalse)
	}

	return SuccessOf(x)
}

// Try2Of returns a Try.
// errOrFalse must be bool or error type.
// Return Failure if errOrFalse is false or error existing,
// or Success of Tuple2(x1,x2).
func Try2Of(x1, x2, errOrFalse interface{}) Try {
	if isErrorOrFalse(errOrFalse) {
		return FailureOf(errOrFalse)
	}

	return SuccessOf(Tuple2Of(x1, x2))
}
