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
	v  reflect.Value
}

func (t *traitTry) Get() interface{} {
	return t.v.Interface()
}

func (t *traitTry) rv() reflect.Value {
	return t.v
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
		funcOf(f).call(t.v)
	}
}

func (t *traitTry) Fold(z, f interface{}) interface{} {

	ztyp := reflect.TypeOf(z)

	if !t.ok {
		if ztyp.Kind() == reflect.Func {
			return funcOf(z).call(t.v).Interface()
		}
		return z
	}

	result := tryCBF(funcOf(f).call(t.v))
	if result.OK() {
		return result.Get()
	}

	return result.Fold(z, nil)
}

func (t *traitTry) Map(f interface{}) Try {
	if t.ok {
		return tryCBF(funcOf(f).call(t.v))
	}

	return t
}

func (t *traitTry) FlatMap(f interface{}) Try {
	if t.ok {
		return funcOf(f).call(t.v).Interface().(Try)
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

	return OptionOf(t.Get())
}

/*func tryFromX(x interface{}) Try {
	return tryFromContainer(containerOf(x), !isErrorOrFalse(x))
}*/

// isErrorOrFalse checks x is an existing error or false.
func isErrorOrFalse(x interface{}) (bool, bool) {
	switch v := x.(type) {
	case error:
		return v != nil, true
	case bool:
		return !v, true
	default:
		return false, x == nil
	}
}

// ----------------------------------------------------------------------------

// FailureOf returns Failure of x if x is error or false, or returns Success of x.
func FailureOf(x interface{}) Try {
	return tryCBF(x)
}

// ----------------------------------------------------------------------------

// SuccessOf returns Success of x.
func SuccessOf(x ...interface{}) Try {
	return tryCBF(x...)
}

var successNull = &traitTry{ok: true, v: nullValue}
var successUnit = &traitTry{ok: true, v: unitValue}

func newTraitTry(ok bool, x interface{}) Try {
	return &traitTry{
		ok: ok,
		v:  reflect.ValueOf(x),
	}
}
func tryCBF(x ...interface{}) (ret Try) {
	len := len(x)

	if len == 0 {
		return successUnit
	}

	switch len {
	case 1:
		switch v := x[0].(type) {
		case Try:
			return v
		case reflect.Value:
			return tryCBF(v.Interface())
		case bool:
			return newTraitTry(v, v)
		case error:
			if v != nil {
				return newTraitTry(false, v)
			}
			return successNull
		case Tuple:
			return tryNOf(v)
		default:
			if v == nil {
				return successNull
			}

			if reflect.TypeOf(v).Kind() == reflect.Func {
				defer func() {
					if r := recover(); r != nil {
						ret = newTraitTry(false, reflect.ValueOf(fmt.Errorf("%v", r)))

					}
				}()

				return tryCBF(funcOf(v).call(unitValue).Interface())
			}
			return newTraitTry(true, v)
		}
	case 2:
		return try1Of(x[0], x[1])
	case 3:
		return try2Of(x[0], x[1], x[2])
	default:
		return tryNOf(TupleOf(x))
	}
}

// TryOf returns a Try.
// The last argument must be bool or error type.
// Return Failure if errOrFalse is false or error existing,
// or Success of Null.
func TryOf(x ...interface{}) (ret Try) {
	return tryCBF(x...)
}

// try1Of returns a Try.
// errOrFalse must be bool or error type.
// Return Failure if errOrFalse is false or error existing,
// or Success of x.
func try1Of(x, errOrFalse interface{}) Try {
	yes, reduce := isErrorOrFalse(errOrFalse)
	if yes {
		return newTraitTry(false, errOrFalse)
	}

	if reduce {
		if x == nil {
			return successNull
		}
		return newTraitTry(true, x)
	}

	return newTraitTry(true, Tuple2Of(x, errOrFalse))
}

// try2Of returns a Try.
// errOrFalse must be bool or error type.
// Return Failure if errOrFalse is false or error existing,
// or Success of Tuple2(x1,x2).
func try2Of(x1, x2, errOrFalse interface{}) Try {
	yes, reduce := isErrorOrFalse(errOrFalse)
	if yes {
		return newTraitTry(false, errOrFalse)
	}

	if reduce {
		return newTraitTry(true, Tuple2Of(x1, x2))
	}

	return newTraitTry(true, Tuple3Of(x1, x2, errOrFalse))
}

func tryNOf(t Tuple) Try {
	last := t.V(t.Dimension() - 1)
	yes, reduce := isErrorOrFalse(last)
	if yes {
		return newTraitTry(false, last)
	}

	if !reduce {
		return newTraitTry(true, t)
	}

	if t.Dimension() == 2 {
		if t.V(0) == nil {
			return successNull
		}
		return newTraitTry(true, t.V(0))
	}

	return newTraitTry(true, t.reduce())
}
