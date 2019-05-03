package monadgo

import (
	"fmt"
	"reflect"
)

// Try represents Scala Try
type Try interface {
	Any

	// Traversable methods start
	Forall(f interface{}) bool
	Foreach(f interface{})
	Fold(z, f interface{}) interface{}
	ToSlice() Slice
	// Traversable methods end

	OK() bool
	Failed() bool

	Map(f interface{}) Try
	FlatMap(f interface{}) Try

	GetOrElse(z interface{}) interface{}
	OrElse(TryOrElse) Try

	ToOption() Option
}

// TryOrElse represents a function applied on Try.OrElse.
type TryOrElse func() Try

// ----------------------------------------------------------------------------

type _failure struct {
	f interface{}
}

var _ Try = _failure{}

func (f _failure) Get() interface{} {
	return f.f
}

func (f _failure) String() string {
	return fmt.Sprintf("Failure(%v)", f.f)
}

func (f _failure) OK() bool {
	return false
}

func (f _failure) Failed() bool {
	return true
}

func (f _failure) Map(interface{}) Try {
	return f
}

func (f _failure) FlatMap(interface{}) Try {
	return f
}

func (f _failure) Forall(interface{}) bool {
	return false
}

func (f _failure) Foreach(interface{}) {

}

func (f _failure) Fold(z, _ interface{}) interface{} {
	ztyp := reflect.TypeOf(z)
	if ztyp.Kind() != reflect.Func {
		return z
	}
	zw := funcOf(z)
	return zw.invoke(f.Get())
}

func (f _failure) ToSlice() Slice {
	switch v := f.f.(type) {
	case error:
		return SliceOf([]error{v})
	case bool:
		return SliceOf([]bool{v})
	default:
		return nil
	}
}

func (f _failure) GetOrElse(z interface{}) interface{} {
	if x, ok := checkFuncAndInvoke(z); ok {
		return x
	}

	return z
}

func (f _failure) OrElse(t TryOrElse) Try {
	return t()
}

func (f _failure) ToOption() Option {
	return None
}

// FailureOf returns Failure of x if x is error or false, or returns Success of x.
func FailureOf(x interface{}) Try {
	if isErrorOrFalse(x) {
		return _failure{x}
	}
	return nil
}

// ----------------------------------------------------------------------------

type _success reflect.Value

var _ Try = _success{}

func (s _success) Get() interface{} {
	return reflect.Value(s).Interface()
}

func (s _success) String() string {
	return fmt.Sprintf("Success(%v)", s.Get())
}

func (s _success) OK() bool {
	return true
}

func (s _success) Failed() bool {
	return false
}

func (s _success) Map(f interface{}) Try {
	return SuccessOf(funcOf(f).invoke(s.Get()))
}

func (s _success) FlatMap(f interface{}) Try {
	return funcOf(f).invoke(s.Get()).(Try)
}

func (s _success) Forall(f interface{}) bool {
	return funcOf(f).invoke(s.Get()).(bool)
}

func (s _success) Foreach(f interface{}) {
	funcOf(f).invoke(s.Get())
}

func (s _success) Fold(z, f interface{}) interface{} {
	ret := funcOf(f).invoke(s.Get())

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

	return TryOf(x).Fold(z, nil)
}

func (s _success) ToSlice() Slice {
	return newSlice(oneToSlice(reflect.ValueOf(s.Get())))
}

func (s _success) GetOrElse(interface{}) interface{} {
	return s.Get()
}

func (s _success) OrElse(TryOrElse) Try {
	return s
}

func (s _success) ToOption() Option {
	return OptionOf(s.Get())
}

// SuccessOf returns Success of x.
func SuccessOf(x interface{}) Try {
	if x == nil {
		return _success(nullValue)
	}

	return _success(reflect.ValueOf(x))
}

// ----------------------------------------------------------------------------

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

// tryFromTuple return a Try from Tuple.
func tryFromTuple(t Tuple) Try {
	x := t.V(t.Dimension() - 1)
	if f := FailureOf(x); f != nil {
		return f
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

// TryOf returns a Try with one parameter.
// errOrFalse must be bool or error type.
// Return Failure if errOrFalse is false or error existing,
// or Success of Null.
func TryOf(errOrFalse interface{}) Try {
	if t := FailureOf(errOrFalse); t != nil {
		return t
	}
	return SuccessOf(errOrFalse)
}

// Try1Of returns a Try.
// errOrFalse must be bool or error type.
// Return Failure if errOrFalse is false or error existing,
// or Success of x.
func Try1Of(x, errOrFalse interface{}) Try {
	if f := FailureOf(errOrFalse); f != nil {
		return f
	}

	return SuccessOf(x)
}

// Try2Of returns a Try.
// errOrFalse must be bool or error type.
// Return Failure if errOrFalse is false or error existing,
// or Success of Tuple2(x1,x2).
func Try2Of(x1, x2, errOrFalse interface{}) Try {
	if f := FailureOf(errOrFalse); f != nil {
		return f
	}

	return SuccessOf(Tuple2Of(x1, x2))
}
