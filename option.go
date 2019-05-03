package monadgo

import (
	"fmt"
	"reflect"
)

// Option represents Scala Option.
type Option interface {
	Traversable

	Defined() bool
	Map(f interface{}) Option
	FlatMap(f interface{}) Option

	OrElse(z OptionOrElse) Option
	GetOrElse(z interface{}) interface{}
}

// OptionOrElse represents a alternative function for Option.OrElse.
type OptionOrElse func() Option

// ----------------------------------------------------------------------------

type _none reflect.Value

// None respresents Scala None in Option.
var None Option = &_none{}

func (n *_none) Get() interface{} {
	return nothing
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

func (n *_none) ToSeq() interface{} {
	return []Nothing{nothing}
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

type _some reflect.Value

func (s _some) Get() interface{} {
	return reflect.Value(s).Interface()
}

func (s _some) String() string {
	return fmt.Sprintf("Some(%v)", s.Get())
}

func (s _some) Defined() bool {
	return true
}

func (s _some) Map(f interface{}) Option {
	fw := funcOf(f)
	return OptionOf(fw.invoke(s.Get()))
}

func (s _some) FlatMap(f interface{}) Option {
	fw := funcOf(f)
	return fw.invoke(s.Get()).(Option)
}

func (s _some) Forall(f interface{}) bool {
	fw := funcOf(f)
	return fw.invoke(s.Get()).(bool)
}

func (s _some) Foreach(f interface{}) {
	fw := funcOf(f)
	fw.invoke(s.Get())
}

func (s _some) Fold(_, f interface{}) interface{} {
	fw := funcOf(f)
	return fw.invoke(s.Get())
}

func (s _some) ToSeq() interface{} {
	seq := makeSlice(reflect.TypeOf(s.Get()), 1)
	seq.Index(0).Set(reflect.ValueOf(s.Get()))
	return seq.Interface()
}

func (s _some) GetOrElse(interface{}) interface{} {
	return s.Get()
}

func (s _some) OrElse(OptionOrElse) Option {
	return s
}

// SomeOf returns Some of x.
func SomeOf(x interface{}) Option {
	if x == nil {
		return _some(nullValue)
	}

	return _some(reflect.ValueOf(x))
}

// ----------------------------------------------------------------------------

// OptionOf returns an Option.
func OptionOf(x interface{}) Option {
	return SomeOf(x)
}
