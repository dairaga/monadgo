package monadgo

import "reflect"

var (
	typeValues = []reflect.Type{reflect.TypeOf(reflect.Value{})}
	typeError  = reflect.TypeOf((*error)(nil)).Elem()
)

var (
	typeAny = reflect.TypeOf((*Any)(nil)).Elem()

	unit      Unit = _unit{}
	unitValue      = reflect.ValueOf(unit)
	typeUnit       = reflect.TypeOf(unit)

	null      Null = &_null{}
	nullValue      = reflect.ValueOf(null)
	typeNull       = reflect.TypeOf(null)

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
	typeTuple2 = reflect.TypeOf(Tuple2{})
	typeTuple3 = reflect.TypeOf(Tuple3{})
	typeTuple4 = reflect.TypeOf(Tuple4{})

	typePair = reflect.TypeOf(Pair{})

	typeTuples = []reflect.Type{typeTuple2, typeTuple3, typeTuple4, typeTuple}
)
