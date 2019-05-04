package monadgo

import (
	"fmt"
	"reflect"
)

// Tuple4 represents scala-like Tuple3.
type Tuple4 interface {
	Tuple
	T1() reflect.Type
	V1() interface{}

	T2() reflect.Type
	V2() interface{}

	T3() reflect.Type
	V3() interface{}

	T4() reflect.Type
	V4() interface{}
}

type _tuple4 struct {
	types  [4]reflect.Type
	values [4]interface{}
	vals   [4]reflect.Value
}

var _ Tuple4 = _tuple4{}

func (t _tuple4) Get() interface{} {
	return t.values
}

func (t _tuple4) rv() reflect.Value {
	return reflect.ValueOf(t)
}

func (t _tuple4) String() string {
	return fmt.Sprintf("(%v,%v,%v,%v)", t.values[0], t.values[1], t.values[2], t.values[3])
}

func (t _tuple4) Dimension() int {
	return 4
}

func (t _tuple4) toValues() []reflect.Value {
	return t.vals[0:]
}

func (t _tuple4) reduce() Tuple {
	return formTuple3(t.types[0], t.types[1], t.types[2], t.values[0], t.values[1], t.values[2])
}

func (t _tuple4) T(n int) reflect.Type {
	return t.types[n]
}

func (t _tuple4) V(n int) interface{} {
	return t.values[n]
}

func (t _tuple4) T1() reflect.Type {
	return t.types[0]
}

func (t _tuple4) T2() reflect.Type {
	return t.types[1]
}

func (t _tuple4) T3() reflect.Type {
	return t.types[2]
}

func (t _tuple4) T4() reflect.Type {
	return t.types[3]
}

func (t _tuple4) V1() interface{} {
	return t.values[0]
}

func (t _tuple4) V2() interface{} {
	return t.values[1]
}

func (t _tuple4) V3() interface{} {
	return t.values[2]
}

func (t _tuple4) V4() interface{} {
	return t.values[3]
}

// ----------------------------------------------------------------------------

// Tuple4Of returns a Tuple4.
func Tuple4Of(v1, v2, v3, v4 interface{}) Tuple4 {
	return formTuple4(reflect.TypeOf(v1), reflect.TypeOf(v2), reflect.TypeOf(v3), reflect.TypeOf(v4), v1, v2, v3, v4)
}

func formTuple4(t1, t2, t3, t4 reflect.Type, v1, v2, v3, v4 interface{}) Tuple4 {
	return _tuple4{
		types:  [4]reflect.Type{t1, t2, t3, t4},
		values: [4]interface{}{v1, v2, v3, v4},
		vals:   [4]reflect.Value{reflect.ValueOf(v1), reflect.ValueOf(v2), reflect.ValueOf(v3), reflect.ValueOf(v4)},
	}
}

// newTuple4 returns a Tuple4.
func newTuple4(t1, t2, t3, t4 reflect.Type, v1, v2, v3, v4 reflect.Value) _tuple4 {
	return _tuple4{
		types:  [4]reflect.Type{t1, t2, t3, t4},
		values: [4]interface{}{v1.Interface(), v2.Interface(), v3.Interface(), v4.Interface()},
		vals:   [4]reflect.Value{v1, v2, v3, v4},
	}
}
