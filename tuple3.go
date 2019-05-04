package monadgo

import (
	"fmt"
	"reflect"
)

// Tuple3 represents scala-like Tuple3.
type Tuple3 interface {
	Tuple
	T1() reflect.Type
	V1() interface{}

	T2() reflect.Type
	V2() interface{}

	T3() reflect.Type
	V3() interface{}
}

type _tuple3 struct {
	types  [3]reflect.Type
	values [3]interface{}
	vals   [3]reflect.Value
}

var _ Tuple3 = _tuple3{}

func (t _tuple3) Get() interface{} {
	return t.values
}

func (t _tuple3) rv() reflect.Value {
	return reflect.ValueOf(t)
}

func (t _tuple3) String() string {
	return fmt.Sprintf("(%v,%v,%v)", t.values[0], t.values[1], t.values[2])
}

func (t _tuple3) Dimension() int {
	return 3
}

func (t _tuple3) toValues() []reflect.Value {
	return t.vals[0:]
}

func (t _tuple3) reduce() Tuple {
	return formTuple2(t.types[0], t.types[1], t.values[0], t.values[1])
}

func (t _tuple3) T(n int) reflect.Type {
	return t.types[n]
}

func (t _tuple3) V(n int) interface{} {
	return t.values[n]
}

func (t _tuple3) T1() reflect.Type {
	return t.types[0]
}

func (t _tuple3) T2() reflect.Type {
	return t.types[1]
}

func (t _tuple3) T3() reflect.Type {
	return t.types[2]
}

func (t _tuple3) V1() interface{} {
	return t.values[0]
}

func (t _tuple3) V2() interface{} {
	return t.values[1]
}

func (t _tuple3) V3() interface{} {
	return t.values[2]
}

// ----------------------------------------------------------------------------

// Tuple3Of returns a Tuple3.
func Tuple3Of(v1, v2, v3 interface{}) Tuple3 {
	return formTuple3(reflect.TypeOf(v1), reflect.TypeOf(v2), reflect.TypeOf(v3), v1, v2, v3)
}

func formTuple3(t1, t2, t3 reflect.Type, v1, v2, v3 interface{}) Tuple3 {
	return _tuple3{
		types:  [3]reflect.Type{t1, t2, t3},
		values: [3]interface{}{v1, v2, v3},
		vals:   [3]reflect.Value{reflect.ValueOf(v1), reflect.ValueOf(v2), reflect.ValueOf(v3)},
	}
}

// Tuple3Of returns a Tuple3.
func newTuple3(t1, t2, t3 reflect.Type, v1, v2, v3 reflect.Value) _tuple3 {
	return _tuple3{
		types:  [3]reflect.Type{t1, t2, t3},
		values: [3]interface{}{v1.Interface(), v2.Interface(), v3.Interface()},
		vals:   [3]reflect.Value{v1, v2, v3},
	}
}
