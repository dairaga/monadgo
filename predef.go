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

type _unit struct{}

func (u _unit) Get() interface{} {
	return unit
}

func (u _unit) String() string {
	return "void"
}

var unit = _unit{}
var unitValue = reflect.ValueOf(unit)

// ----------------------------------------------------------------------------

// Null represents scala-like Null.
type Null interface {
	Any
}

type _null struct{}

func (n *_null) Get() interface{} {
	return nil
}

func (n *_null) String() string {
	return "null"
}

var null Null = &_null{}
var nullValue = reflect.ValueOf(null)

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
