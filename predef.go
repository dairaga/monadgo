package monadgo

import (
	"fmt"
	"reflect"
)

// Unit represents scala-like Unit.
type Unit interface {
	Any
}

var _ Unit = _unit{}
var unit = _unit{}
var unitValue = reflect.ValueOf(unit)

type _unit struct{}

func (u _unit) Get() interface{} {
	return unit
}

func (u _unit) rv() reflect.Value {
	return unitValue
}

func (u _unit) String() string {
	return "void"
}

// ----------------------------------------------------------------------------

// Null represents scala-like Null.
type Null interface {
	Any
}

type _null struct{}

var null Null = &_null{}
var nullValue = reflect.ValueOf(null)

func (n *_null) Get() interface{} {
	return nil
}

func (n *_null) rv() reflect.Value {
	return nullValue
}

func (n *_null) String() string {
	return "null"
}

// ----------------------------------------------------------------------------

// Nothing represents scala Nothing.
type Nothing interface {
	fmt.Stringer
}

type _nothing struct{}

func (n *_nothing) String() string {
	return "Nothing"
}

var nothing Nothing = &_nothing{}
var nothings []Nothing

// ----------------------------------------------------------------------------

// CanBuildFrom constructs object.
type CanBuildFrom func(reflect.Value) reflect.Value

// Build ...
func (b CanBuildFrom) Build(v reflect.Value) reflect.Value {
	return b(v)
}
