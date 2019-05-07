package monadgo

import "reflect"

var (
	typeValues = []reflect.Type{reflect.TypeOf(reflect.Value{})}
	typeError  = reflect.TypeOf((*error)(nil)).Elem()
)

var (
	unit      = _unit(false)
	unitValue = reflect.ValueOf(unit)
	typeUnit  = reflect.TypeOf((*Unit)(nil)).Elem()

	null      = &_null{}
	nullValue = reflect.ValueOf(null)
	typeNull  = reflect.TypeOf((*Null)(nil)).Elem()

	nothing      = &_nothing{null}
	nothingValue = reflect.ValueOf(nothing)
	typeNothing  = reflect.TypeOf((*Nothing)(nil)).Elem()

	nothings      = []Nothing{}
	typeNothings  = reflect.TypeOf(nothings)
	nothingsValue = reflect.ValueOf(nothings)

	typeSeq = reflect.TypeOf((*sequence)(nil)).Elem()
)

var (
	typeTuple  = reflect.TypeOf((*Tuple)(nil)).Elem()
	typeTuple2 = reflect.TypeOf((*Tuple2)(nil)).Elem()
	typeTuple3 = reflect.TypeOf((*Tuple3)(nil)).Elem()
	typeTuple4 = reflect.TypeOf((*Tuple4)(nil)).Elem()

	typePair = reflect.TypeOf((*Pair)(nil)).Elem()
)
