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

	// OK returns true if this is Success.
	OK() bool

	// Failed returns true if this is Failure.
	Failed() bool

	// Foreach applies f to Try's value if this is Success.
	// f: func(T)
	Foreach(f interface{})

	// Fold applies z if this is a Failure,
	// or f if this is a Success.
	// If f is initially applied and last element in results is false or error,
	// then z applied with this element value.
	// z: func(A) X. A can be error or bool.
	// f: func(B) X. B is the element type in Success.
	// returns value with type X.
	Fold(z, f interface{}) interface{}

	// Map applies f to the value from this Success,
	// or returns this if this is a Failure.
	// f: func(T) X
	Map(f interface{}) Try

	// FlatMap returns the f applied to the value from this Success,
	// or returns this if this is a Failure.
	// f: func(T) Try
	FlatMap(f interface{}) Try

	// OrElse returns this if it's a Success,
	// or z if this is a Failure.
	OrElse(z TryOrElse) Try

	// GetOrElse returns the value from this Success,
	// or z if this is a Failure.
	GetOrElse(z interface{}) interface{}

	// ToOption returns None if this is a Failure,
	// or a Some containing Success's value.
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

func tryFromX(x interface{}) Try {
	if isErrorOrFalse(x) {
		return FailureOf(x)
	}

	switch v := x.(type) {
	case Tuple:
		switch v.Dimension() {
		case 2:
			return try1Of(v.V(0), v.V(1))
		case 3:
			return try2Of(v.V(0), v.V(1), v.V(2))
		default:
			return tryNOf(v)
		}
	default:
		return SuccessOf(x)
	}

}

// TryOf returns a Try.
// The last argument must be bool or error type.
// Return Failure if errOrFalse is false or error existing,
// or Success of Null.
func TryOf(x ...interface{}) Try {
	switch len(x) {
	case 0:
		return SuccessOf(unit)
	case 1:
		if isErrorOrFalse(x[0]) {
			return FailureOf(x[0])
		}
		return SuccessOf(x[0])
	case 2:
		return try1Of(x[0], x[1])
	case 3:
		return try2Of(x[0], x[1], x[2])
	default:
		t := TupleOf(x)
		return tryNOf(t)
	}

}

// try1Of returns a Try.
// errOrFalse must be bool or error type.
// Return Failure if errOrFalse is false or error existing,
// or Success of x.
func try1Of(x, errOrFalse interface{}) Try {
	if isErrorOrFalse(errOrFalse) {
		return FailureOf(errOrFalse)
	}

	return SuccessOf(x)
}

// try2Of returns a Try.
// errOrFalse must be bool or error type.
// Return Failure if errOrFalse is false or error existing,
// or Success of Tuple2(x1,x2).
func try2Of(x1, x2, errOrFalse interface{}) Try {
	if isErrorOrFalse(errOrFalse) {
		return FailureOf(errOrFalse)
	}

	return SuccessOf(Tuple2Of(x1, x2))
}

func tryNOf(t Tuple) Try {
	last := t.V(t.Dimension() - 1)
	if isErrorOrFalse(last) {
		return FailureOf(last)
	}
	return SuccessOf(t.reduce())
}
