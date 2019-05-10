package monadgo

import (
	"fmt"
	"reflect"
	"strings"
)

// Tuple represents scala-like Tuple.
type Tuple interface {
	Any

	// Dimension returns number of dimension.
	Dimension() int

	// T returns the type of n-index element.
	T(n int) reflect.Type

	// V returns the value of n-index element.
	V(n int) interface{}

	// toValues returns slice of reflect.Value in Tuple.
	toValues() []reflect.Value

	// reduce reduces tuple.
	reduce() Tuple
}

// ----------------------------------------------------------------------------

// TupleN implements Tuple.
type TupleN struct {
	d      int
	types  []reflect.Type
	values []interface{}
	vals   []reflect.Value
}

var _ Tuple = TupleN{}

// Get returns all elements in Go slice.
func (t TupleN) Get() interface{} {
	return t.values
}

func (t TupleN) rv() reflect.Value {
	return reflect.ValueOf(t)
}

func (t TupleN) String() string {
	sb := new(strings.Builder)
	sb.WriteByte('(')
	sb.WriteString(fmt.Sprintf("%v", t.values[0]))
	for i := 1; i < t.d; i++ {
		sb.WriteString(fmt.Sprintf(",%v", t.values[i]))
	}
	sb.WriteByte(')')
	return sb.String()
}

// Dimension return number of dimension of t.
func (t TupleN) Dimension() int {
	return t.d
}

// T returns type of n-index element.
func (t TupleN) T(n int) reflect.Type {
	return t.types[n]
}

// V returns value of n-index element.
func (t TupleN) V(n int) interface{} {
	return t.values[n]
}

// reduce returns a new reduced tuple from t.
func (t TupleN) reduce() Tuple {
	d := t.d - 1
	switch d {
	case 4:
		return formTuple4(t.T(0), t.T(1), t.T(2), t.T(3), t.V(0), t.V(1), t.V(2), t.V(3))
	default:
		return formTuple(t.types[0:d], t.values[0:d])
	}
}

func (t TupleN) toValues() []reflect.Value {
	return t.vals
}

// ----------------------------------------------------------------------------

// TupleOf returns a general Tuple.
func TupleOf(v []interface{}) Tuple {
	size := len(v)
	switch size {
	case 2:
		return Tuple2Of(v[0], v[1])
	case 3:
		return Tuple3Of(v[0], v[1], v[2])
	case 4:
		return Tuple4Of(v[0], v[1], v[2], v[3])
	}

	t := make([]reflect.Type, size, size)
	for i := 0; i < size; i++ {
		t[i] = reflect.TypeOf(v[i])
	}

	return formTuple(t, v)
}

// ----------------------------------------------------------------------------

// tupleOf ...
func formTuple(t []reflect.Type, v []interface{}) Tuple {

	switch len(t) {
	case 2:
		return formTuple2(t[0], t[1], v[0], v[1])
	case 3:
		return formTuple3(t[0], t[1], t[2], v[0], v[1], v[2])
	case 4:
		return formTuple4(t[0], t[1], t[2], t[3], v[0], v[1], v[2], v[3])
	}

	ret := &TupleN{
		d:      len(t),
		types:  t,
		values: v,
		vals:   make([]reflect.Value, len(t), len(t)),
	}

	for i, x := range v {
		ret.vals[i] = reflect.ValueOf(x)
	}

	return ret
}

// ----------------------------------------------------------------------------

// newTuple ...
func newTuple(t []reflect.Type, v []reflect.Value) Tuple {

	switch len(t) {
	case 2:
		return newTuple2(t[0], t[1], v[0], v[1])
	case 3:
		return newTuple3(t[0], t[1], t[2], v[0], v[1], v[2])
	case 4:
		return newTuple4(t[0], t[1], t[2], t[3], v[0], v[1], v[2], v[3])
	}

	ret := TupleN{
		d:      len(t),
		types:  t,
		values: make([]interface{}, len(t), len(t)),
		vals:   v,
	}

	for i, x := range v {
		ret.values[i] = x.Interface()
	}

	return ret
}
