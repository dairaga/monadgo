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

// Unit represents scala-like Unit.
type _unit struct{}

func (u _unit) Get() interface{} {
	return Unit
}

func (u _unit) rv() reflect.Value {
	return unitValue
}

func (u _unit) String() string {
	return "void"
}

// ----------------------------------------------------------------------------

// Null represents scala-like Null.

type _null struct{}

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
	Any
}

type _nothing struct{}

func (n *_nothing) Get() interface{} {
	return nothing
}

func (n *_nothing) String() string {
	return "Nothing"
}

func (n *_nothing) rv() reflect.Value {
	return nothingValue
}

// ----------------------------------------------------------------------------

// CanBuildFrom constructs object.
type CanBuildFrom func(reflect.Value) reflect.Value

// Build ...
func (b CanBuildFrom) Build(v reflect.Value) reflect.Value {
	return b(v)
}

// ----------------------------------------------------------------------------

// PartialFunc is a scala-like PartialFunction.
type PartialFunc struct {
	condition reflect.Value
	action    reflect.Value
}

// PartialFuncOf ...
func PartialFuncOf(c, a interface{}) PartialFunc {
	return PartialFunc{
		condition: reflect.ValueOf(c),
		action:    reflect.ValueOf(a),
	}
}

// IsDefinedAt ...
func (p PartialFunc) IsDefinedAt(x interface{}) bool {
	switch p.condition.Kind() {
	case reflect.Bool:
		return p.condition.Bool()
	case reflect.Func:
		return p.condition.Call([]reflect.Value{reflect.ValueOf(x)})[0].Bool()
	}
}

// Call ...
func (p PartialFunc) Call(x interface{}) interface{} {
	if p.IsDefinedAt(x) {
		if p.action.Kind() == reflect.Func {
			return p.action.Call([]reflect.Value{reflect.ValueOf(x)})[0].Interface()
		}
		return p.action.Interface()
	}
	return nothing
}
