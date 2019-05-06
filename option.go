package monadgo

import (
	"fmt"
	"reflect"
)

// Option represents Scala Option.
type Option interface {
	Any
	// Traversable methods start
	Forall(f interface{}) bool
	Foreach(f interface{})
	Fold(z, f interface{}) interface{}
	// Traversable methods end

	Defined() bool
	Map(f interface{}) Option
	FlatMap(f interface{}) Option

	OrElse(z OptionOrElse) Option
	GetOrElse(z interface{}) interface{}
}

// OptionOrElse represents a alternative function for Option.OrElse.
type OptionOrElse func() Option

var optionCBF CanBuildFrom = func(v reflect.Value) reflect.Value {
	if v == nothingValue {
		return None
	}

	return someFromValue(v)
}

// OptionOf ...
func OptionOf(x interface{}) (result Option) {
	defer func() {
		if r := recover(); r != nil {
			result = None
		}
	}()

	if v, ok := checkFuncAndInvoke(x); ok {
		return SomeOf(v)
	}

	return SomeOf(x)
}

// ----------------------------------------------------------------------------

type _none struct {
	container
}

// None respresents Scala None in Option.
var None Option = &_none{containerFromValue(nothingValue)}
var noneValue = reflect.ValueOf(None)

func (n *_none) rv() reflect.Value {
	return noneValue
}

func (n *_none) String() string {
	return "None"
}

func (n *_none) Defined() bool {
	return false
}

func (n *_none) Map(interface{}) Option {
	return None
}

func (n *_none) FlatMap(interface{}) Option {
	return None
}

func (n *_none) Forall(interface{}) bool {
	return false
}

func (n *_none) Foreach(interface{}) {

}

func (n *_none) Fold(z, _ interface{}) interface{} {
	if x, ok := checkFuncAndInvoke(z); ok {
		return x
	}
	return z
}

func (n *_none) GetOrElse(z interface{}) interface{} {
	if x, ok := checkFuncAndInvoke(z); ok {
		return x
	}

	return z
}

func (n *_none) OrElse(z OptionOrElse) Option {
	return z()
}

// ----------------------------------------------------------------------------

type _some struct {
	container
}

var _ Option = _some{}

func someFromValue(v reflect.Value) _some {
	if !v.IsValid() {
		v = nullValue
	}

	return _some{containerFromValue(v)}
}

// SomeOf returns Some of x.
func SomeOf(x interface{}) Option {
	if x == null {
		return someFromValue(nullValue)
	}

	switch v := x.(type) {
	case reflect.Value:
		return someFromValue(v)
	default:
		return someFromValue(reflect.ValueOf(x))
	}
}

func (s _some) String() string {
	return fmt.Sprintf("Some(%v)", s.Get())
}

func (s _some) Defined() bool {
	return true
}

func (s _some) Map(f interface{}) Option {
	return s.container._map(f, optionCBF).(Option)
}

func (s _some) FlatMap(f interface{}) Option {
	return s.container.flatMap(f).(Option)
}

func (s _some) Fold(_, f interface{}) interface{} {
	return s.container.invoke(f).Interface()
}

func (s _some) GetOrElse(interface{}) interface{} {
	return s.Get()
}

func (s _some) OrElse(OptionOrElse) Option {
	return s
}
