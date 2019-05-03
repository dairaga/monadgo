package monadgo

import (
	"fmt"
	"reflect"
)

// Tuple2 represents scala-like Tuple2
type Tuple2 interface {
	Tuple
	T1() reflect.Type
	V1() interface{}

	T2() reflect.Type
	V2() interface{}
}

type _tuple2 struct {
	types  [2]reflect.Type
	values [2]interface{}
	vals   [2]reflect.Value
}

var _ Tuple2 = _tuple2{}

func (t _tuple2) Get() interface{} {
	return t.values
}

func (t _tuple2) String() string {
	return fmt.Sprintf("(%v,%v)", t.values[0], t.values[1])
}

func (t _tuple2) Dimension() int {
	return 2
}

func (t _tuple2) toValues() []reflect.Value {
	return t.vals[0:]
}

func (t _tuple2) reduce() Tuple {
	return t
}

func (t _tuple2) T(n int) reflect.Type {
	return t.types[n]
}

func (t _tuple2) V(n int) interface{} {
	return t.values[n]
}

func (t _tuple2) T1() reflect.Type {
	return t.types[0]
}

func (t _tuple2) T2() reflect.Type {
	return t.types[1]
}

func (t _tuple2) V1() interface{} {
	return t.values[0]
}

func (t _tuple2) V2() interface{} {
	return t.values[1]
}

// ----------------------------------------------------------------------------

// Tuple2Of returns a Tuple2.
func Tuple2Of(v1, v2 interface{}) Tuple2 {
	return formTuple2(reflect.TypeOf(v1), reflect.TypeOf(v2), v1, v2)
}

func formTuple2(t1, t2 reflect.Type, v1, v2 interface{}) Tuple2 {
	return _tuple2{
		types:  [2]reflect.Type{t1, t2},
		values: [2]interface{}{v1, v2},
		vals:   [2]reflect.Value{reflect.ValueOf(v1), reflect.ValueOf(v2)},
	}
}

// newTuple2 returns a Tuple2.
func newTuple2(t1, t2 reflect.Type, v1, v2 reflect.Value) _tuple2 {
	return _tuple2{
		types:  [2]reflect.Type{t1, t2},
		values: [2]interface{}{v1.Interface(), v2.Interface()},
		vals:   [2]reflect.Value{v1, v2},
	}
}
