package monadgo

import (
	"reflect"
)

func checkAndInvoke(x interface{}) interface{} {
	xval := reflect.ValueOf(x)

	if xval.Kind() == reflect.Func {
		return xval.Call(nil)[0].Interface()
	}

	return x
}

// ----------------------------------------------------------------------------

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

func oneToSlice(v reflect.Value) reflect.Value {
	s := makeSlice(v.Type(), 1, 1)
	s.Index(0).Set(v)
	return s
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

// ----------------------------------------------------------------------------

// makeMap returns a reflect.Value of go map with k->v.
func makeMap(k, v reflect.Type, size int) reflect.Value {
	if size < 0 {
		return reflect.MakeMap(reflect.MapOf(k, v))
	}

	return reflect.MakeMapWithSize(reflect.MapOf(k, v), size)
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
	return nothingValue
}

// mergeKeyValue returns a reflect.Value of go map that add k->v to m.
func mergeKeyValue(m, k, v reflect.Value) reflect.Value {
	if !m.IsValid() {
		m = makeMap(k.Type(), v.Type(), -1)
	}

	m.SetMapIndex(k, v)
	return m
}
