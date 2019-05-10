package monadgo

import (
	"fmt"
	"reflect"
)

// Tuple3 represents scala-like Tuple3.
type Tuple3 struct {
	types  [3]reflect.Type
	values [3]interface{}
	vals   [3]reflect.Value
}

var _ Tuple3 = Tuple3{}

// Get ...
func (t Tuple3) Get() interface{} {
	return t.values
}

func (t Tuple3) rv() reflect.Value {
	return reflect.ValueOf(t)
}

func (t Tuple3) String() string {
	return fmt.Sprintf("(%v,%v,%v)", t.values[0], t.values[1], t.values[2])
}

// Dimension ...
func (t Tuple3) Dimension() int {
	return 3
}

func (t Tuple3) toValues() []reflect.Value {
	return t.vals[0:]
}

func (t Tuple3) reduce() Tuple {
	return formTuple2(t.types[0], t.types[1], t.values[0], t.values[1])
}

// T ...
func (t Tuple3) T(n int) reflect.Type {
	return t.types[n]
}

// V ...
func (t Tuple3) V(n int) interface{} {
	return t.values[n]
}

// T1 returns type of first element.
func (t Tuple3) T1() reflect.Type {
	return t.types[0]
}

// T2 returns type of second element.
func (t Tuple3) T2() reflect.Type {
	return t.types[1]
}

// T3 returns type of third element.
func (t Tuple3) T3() reflect.Type {
	return t.types[2]
}

// V1 returns value of first element.
func (t Tuple3) V1() interface{} {
	return t.values[0]
}

// V2 returns value of second element.
func (t Tuple3) V2() interface{} {
	return t.values[1]
}

// V3 returns value of third element.
func (t Tuple3) V3() interface{} {
	return t.values[2]
}

// ----------------------------------------------------------------------------

// Tuple3Of returns a Tuple3.
func Tuple3Of(v1, v2, v3 interface{}) Tuple3 {
	return formTuple3(reflect.TypeOf(v1), reflect.TypeOf(v2), reflect.TypeOf(v3), v1, v2, v3)
}

func formTuple3(t1, t2, t3 reflect.Type, v1, v2, v3 interface{}) Tuple3 {
	return Tuple3{
		types:  [3]reflect.Type{t1, t2, t3},
		values: [3]interface{}{v1, v2, v3},
		vals:   [3]reflect.Value{reflect.ValueOf(v1), reflect.ValueOf(v2), reflect.ValueOf(v3)},
	}
}

func newTuple3(t1, t2, t3 reflect.Type, v1, v2, v3 reflect.Value) Tuple3 {
	return Tuple3{
		types:  [3]reflect.Type{t1, t2, t3},
		values: [3]interface{}{v1.Interface(), v2.Interface(), v3.Interface()},
		vals:   [3]reflect.Value{v1, v2, v3},
	}
}
