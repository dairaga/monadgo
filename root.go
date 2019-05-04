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

// interface types
var (
	typeValue = reflect.TypeOf(reflect.Value{})

	typeUnit = reflect.TypeOf((*Unit)(nil)).Elem()

	typeTuple  = reflect.TypeOf((*Tuple)(nil)).Elem()
	typeTuple2 = reflect.TypeOf((*Tuple2)(nil)).Elem()
	typeTuple3 = reflect.TypeOf((*Tuple3)(nil)).Elem()
	typeTuple4 = reflect.TypeOf((*Tuple4)(nil)).Elem()

	typePair  = reflect.TypeOf((*Pair)(nil)).Elem()
	typeError = reflect.TypeOf((*error)(nil)).Elem()

	typeSeq = reflect.TypeOf((*sequence)(nil)).Elem()
)

// ----------------------------------------------------------------------------

// makeMap returns a reflect.Value of go map with k->v.
func makeMap(k, v reflect.Type, size int) reflect.Value {
	if size < 0 {
		return reflect.MakeMap(reflect.MapOf(k, v))
	}

	return reflect.MakeMapWithSize(reflect.MapOf(k, v), size)
}

// makeSlice returns a reflect.Value of go slice.
func makeSlice(t reflect.Type, lenAndCap ...int) reflect.Value {
	switch len(lenAndCap) {
	case 1:
		return reflect.MakeSlice(reflect.SliceOf(t), lenAndCap[0], lenAndCap[0])
	case 2:
		return reflect.MakeSlice(reflect.SliceOf(t), lenAndCap[0], lenAndCap[1])
	default:
		return reflect.MakeSlice(reflect.SliceOf(t), 0, 0)
	}
}

// appendSlice returns a reflect.Value of go slice that y appends to x.
func appendSlice(x, y reflect.Value) reflect.Value {
	if !x.IsValid() {
		// if x is not valid, returns a slice of y.
		x = makeSlice(y.Type(), 1, 1)
		x.Index(0).Set(y)
		return x
	}

	return reflect.Append(x, y)

}

// mergeSlice returns a reflect.Value of go slice that x merges y.
func mergeSlice(x, y reflect.Value) reflect.Value {
	if !x.IsValid() {
		return y
	}

	if x.Kind() != reflect.Slice {
		x = oneToSlice(x)
	}

	if y.Kind() == reflect.Slice {
		return reflect.AppendSlice(x, y)
	}

	return reflect.Append(x, y)
}

func makeSliceFromFunc(f interface{}, lenAndCap ...int) reflect.Value {
	ftyp := reflect.TypeOf(f)

	switch ftyp.NumOut() {
	case 0:
		return makeSlice(typeUnit, lenAndCap...)
	case 1:
		return makeSlice(ftyp.Out(0), lenAndCap...)
	case 2:
		return makeSlice(typeTuple2, lenAndCap...)
	case 3:
		return makeSlice(typeTuple3, lenAndCap...)
	case 4:
		return makeSlice(typeTuple4, lenAndCap...)
	default:
		return makeSlice(typeTuple, lenAndCap...)
	}
}

func oneToSlice(v reflect.Value) reflect.Value {
	s := makeSlice(v.Type(), 1, 1)
	s.Index(0).Set(v)
	return s
}

// mergeMap returns a reflect.Value of go map that x merges y.
func mergeMap(x, y reflect.Value) reflect.Value {
	if !x.IsValid() {
		return oneToMap(y)
	}

	if x.Type().Implements(typePair) {
		x = oneToMap(x)
	}

	if y.Type().Implements(typePair) {
		py := y.Interface().(Pair)
		x.SetMapIndex(reflect.ValueOf(py.Key()), reflect.ValueOf(py.Value()))
	} else {
		ity := y.MapRange()
		for ity.Next() {
			x.SetMapIndex(ity.Key(), ity.Value())
		}
	}

	return x
}

func oneToMap(x reflect.Value) reflect.Value {
	if x.Type().Implements(typePair) {
		px := x.Interface().(Pair)
		m := makeMap(reflect.TypeOf(px.Key()), reflect.TypeOf(px.Value()), 1)
		m.SetMapIndex(reflect.ValueOf(px.Key()), reflect.ValueOf(px.Value()))
		return m
	} else if x.Kind() == reflect.Map {
		return x
	}
	return reflect.Value{}
}

// mergeKeyValue returns a reflect.Value of go map that add k->v to m.
func mergeKeyValue(m, k, v reflect.Value) reflect.Value {
	if !m.IsValid() {
		m = makeMap(k.Type(), v.Type(), -1)
	}

	m.SetMapIndex(k, v)
	return m
}
