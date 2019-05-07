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

type _unit struct{}

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

// ----------------------------------------------------------------------------

// PartialFunc represents scala-like PartailFunction.
type PartialFunc struct {
	condition funcTR
	action    funcTR
}

// PartialFuncOf returns a partail function for monandgo.
func PartialFuncOf(c, a interface{}) PartialFunc {
	return PartialFunc{
		condition: funcOf(c),
		action:    funcOf(a),
	}

}

// DefinedAt returns x is defined at p or not.
func (p PartialFunc) DefinedAt(v reflect.Value) bool {
	return p.condition.call(v).Bool()
}

// Call invokes action on x, returns Nothing if x is not defined in p.
func (p PartialFunc) Call(v reflect.Value) reflect.Value {
	if p.DefinedAt(v) {
		return p.action.call(v)
	}
	return nothingValue
}
