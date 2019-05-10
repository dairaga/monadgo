package monadgo

import (
	"fmt"
	"reflect"
)

// Tuple4 represents scala-like Tuple4.
type Tuple4 struct {
	types  [4]reflect.Type
	values [4]interface{}
	vals   [4]reflect.Value
}

var _ Tuple4 = Tuple4{}

// Get ...
func (t Tuple4) Get() interface{} {
	return t.values
}

func (t Tuple4) rv() reflect.Value {
	return reflect.ValueOf(t)
}

func (t Tuple4) String() string {
	return fmt.Sprintf("(%v,%v,%v,%v)", t.values[0], t.values[1], t.values[2], t.values[3])
}

// Dimension ...
func (t Tuple4) Dimension() int {
	return 4
}

func (t Tuple4) toValues() []reflect.Value {
	return t.vals[0:]
}

func (t Tuple4) reduce() Tuple {
	return formTuple3(t.types[0], t.types[1], t.types[2], t.values[0], t.values[1], t.values[2])
}

// T ...
func (t Tuple4) T(n int) reflect.Type {
	return t.types[n]
}

// V ...
func (t Tuple4) V(n int) interface{} {
	return t.values[n]
}

// T1 returns type of first element.
func (t Tuple4) T1() reflect.Type {
	return t.types[0]
}

// T2 returns type of second element.
func (t Tuple4) T2() reflect.Type {
	return t.types[1]
}

// T3 returns type of third element.
func (t Tuple4) T3() reflect.Type {
	return t.types[2]
}

// T4 returns type of fourth element.
func (t Tuple4) T4() reflect.Type {
	return t.types[3]
}

// V1 returns value of first element.
func (t Tuple4) V1() interface{} {
	return t.values[0]
}

// V2 returns value of second element.
func (t Tuple4) V2() interface{} {
	return t.values[1]
}

// V3 returns value of third element.
func (t Tuple4) V3() interface{} {
	return t.values[2]
}

// V4 returns value of fourth element.
func (t Tuple4) V4() interface{} {
	return t.values[3]
}

// ----------------------------------------------------------------------------

// Tuple4Of returns a Tuple4.
func Tuple4Of(v1, v2, v3, v4 interface{}) Tuple4 {
	return formTuple4(reflect.TypeOf(v1), reflect.TypeOf(v2), reflect.TypeOf(v3), reflect.TypeOf(v4), v1, v2, v3, v4)
}

func formTuple4(t1, t2, t3, t4 reflect.Type, v1, v2, v3, v4 interface{}) Tuple4 {
	return Tuple4{
		types:  [4]reflect.Type{t1, t2, t3, t4},
		values: [4]interface{}{v1, v2, v3, v4},
		vals:   [4]reflect.Value{reflect.ValueOf(v1), reflect.ValueOf(v2), reflect.ValueOf(v3), reflect.ValueOf(v4)},
	}
}

// newTuple4 returns a Tuple4.
func newTuple4(t1, t2, t3, t4 reflect.Type, v1, v2, v3, v4 reflect.Value) Tuple4 {
	return Tuple4{
		types:  [4]reflect.Type{t1, t2, t3, t4},
		values: [4]interface{}{v1.Interface(), v2.Interface(), v3.Interface(), v4.Interface()},
		vals:   [4]reflect.Value{v1, v2, v3, v4},
	}
}
