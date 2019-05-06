package monadgo

import (
	"reflect"
)

// Buttom Values
var (
	Unit      = _unit{}
	unitValue = reflect.ValueOf(Unit)

	Null      = &_null{}
	nullValue = reflect.ValueOf(nullValue)

	nothing      = &_nothing{}
	nothingValue = reflect.ValueOf(nothing)
	nothings     = []Nothing{}
)

// interface types
var (
	typeValue = reflect.TypeOf(reflect.Value{})
	typeUnit  = reflect.TypeOf(Unit)

	typeTuple  = reflect.TypeOf((*Tuple)(nil)).Elem()
	typeTuple2 = reflect.TypeOf((*Tuple2)(nil)).Elem()
	typeTuple3 = reflect.TypeOf((*Tuple3)(nil)).Elem()
	typeTuple4 = reflect.TypeOf((*Tuple4)(nil)).Elem()

	typePair  = reflect.TypeOf((*Pair)(nil)).Elem()
	typeError = reflect.TypeOf((*error)(nil)).Elem()

	typeSeq = reflect.TypeOf((*sequence)(nil)).Elem()

	typeMap   = reflect.TypeOf((*Map)(nil)).Elem()
	typeSlice = reflect.TypeOf((*Slice)(nil)).Elem()
)
