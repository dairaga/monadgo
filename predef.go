package monadgo

import (
	"fmt"
	"reflect"
)

// Any respensts root type of monadgo.
type Any interface {
	// Get returns original go value.
	Get() interface{}

	// rv returns reflect.Value wrapping original value.
	rv() reflect.Value

	// force monadgo type must implement String()
	fmt.Stringer
}

// ----------------------------------------------------------------------------

// Unit represents scala-like Unit.
type Unit interface {
	Any
}

type _unit bool

func (u _unit) Get() interface{} {
	return unit
}

func (u _unit) rv() reflect.Value {
	return unitValue
}

func (u _unit) String() string {
	return "Void"
}

// ----------------------------------------------------------------------------

// Null represents scala-like Null.
type Null interface {
	Any
}

type _null struct{}

func (n *_null) Get() interface{} {
	return nil
}

func (n *_null) rv() reflect.Value {
	return nullValue
}

func (n *_null) String() string {
	return "Null"
}

// ----------------------------------------------------------------------------

// Nothing represents scala-like Nothing.
type Nothing interface {
	Null
}

type _nothing struct {
	Null
}

func (n *_nothing) rv() reflect.Value {
	return nothingValue
}

func (n *_nothing) String() string {
	return "Nothing"
}
