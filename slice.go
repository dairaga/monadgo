package monadgo

import (
	"fmt"
	"reflect"
)

// Slice represents a scala-like List.
type Slice interface {
	//Any
	Traversable

	Len() int
	Cap() int

	//Map(f interface{}) Slice
	//FlatMap(f interface{}) Slice
}

type _slice reflect.Value

var _ Slice = _slice{}

// ----------------------------------------------------------------------------

func (s _slice) Get() interface{} {
	return reflect.Value(s).Interface()
}

func (s _slice) String() string {
	return fmt.Sprintf("%v", s.Get())
}

func (s _slice) Len() int {
	return reflect.Value(s).Len()
}

func (s _slice) Size() int {
	return s.Len()
}

func (s _slice) Cap() int {
	return reflect.Value(s).Cap()
}

func (s _slice) Head() interface{} {
	if s.Len() <= 0 {
		return nil
	}

	return reflect.Value(s).Index(0).Interface()
}

func (s _slice) Tail() Traversable {
	sval := reflect.Value(s)
	return newSlice(sval.Slice(1, sval.Len()))
}

func (s _slice) Map(f interface{}) Traversable {
	sval := reflect.Value(s)
	len := sval.Len()
	fw := funcOf(f)

	var ret reflect.Value

	for i := 0; i < len; i++ {
		result := fw.call(sval.Index(i))
		ret = appendSlice(ret, result)
	}

	return newSlice(ret)
}

func (s _slice) FlatMap(f interface{}) Traversable {
	sval := reflect.Value(s)
	len := sval.Len()

	var ret reflect.Value
	fw := funcOf(f)
	for i := 0; i < len; i++ {
		seq := TraversableOf(fw.call(sval.Index(i)))
		ret = mergeSlice(ret, reflect.ValueOf(seq.ToSlice().Get()))
	}

	return newSlice(ret)
}

func (s _slice) Fold(z interface{}, f interface{}) interface{} {
	sval := reflect.Value(s)
	len := sval.Len()
	fw := foldOf(f)

	for i := 0; i < len; i++ {
		z = fw.fold(z, sval.Index(i).Interface())
	}
	return z
}

func (s _slice) Reduce(f interface{}) interface{} {
	sval := reflect.Value(s)
	len := sval.Len()

	if len <= 0 {
		panic("empty list can not reduce")
	}

	if len == 1 {
		return sval.Index(0).Interface()
	}

	return s.Tail().Fold(sval.Index(0).Interface(), f)
}

func (s _slice) Forall(f interface{}) bool {
	sval := reflect.Value(s)
	len := sval.Len()
	fw := funcOf(f)

	for i := 0; i < len; i++ {
		if !fw.invoke(sval.Index(i).Interface()).(bool) {
			return false
		}

	}

	return true
}

func (s _slice) Foreach(f interface{}) {
	sval := reflect.Value(s)
	len := sval.Len()
	fw := funcOf(f)
	for i := 0; i < len; i++ {
		fw.call(sval.Index(i))
	}
}

func (s _slice) ToSlice() Slice {
	return s
}

func (s _slice) ToMap() Map {
	sval := reflect.Value(s)
	if !sval.Type().Elem().Implements(typeTuple2) {
		return nil
	}

	var ret reflect.Value

	len := s.Len()
	for i := 0; i < len; i++ {
		ret = mergeMap(ret, sval.Index(i))
	}

	return newMap(ret)
}

func (s _slice) Scan(z, f interface{}) Traversable {
	zval := oneToSlice(reflect.ValueOf(z))
	len := s.Len()
	if s.Len() <= 0 {
		return newSlice(zval)
	}

	sval := reflect.Value(s)
	fw := foldOf(f)

	for i := 0; i < len; i++ {
		zval = appendSlice(zval, fw.call(zval.Index(i), sval.Index(i)))
	}

	return newSlice(zval)
}

func (s _slice) GroupBy(f interface{}) Map {
	len := s.Len()
	if len <= 0 {
		panic("can not group by on empty slice")
	}
	sval := reflect.Value(s)
	ftyp := reflect.TypeOf(f)
	elm := sval.Type()
	fw := funcOf(f)
	m := makeMap(ftyp.Out(0), elm, -1)

	for i := 0; i < len; i++ {
		k := fw.call(sval.Index(i))
		m.SetMapIndex(k, appendSlice(m.MapIndex(k), sval.Index(i)))
	}

	return newMap(m)
}

func (s _slice) Take(n int) Traversable {
	if n >= s.Len() {
		n = s.Len()
	}

	return newSlice(reflect.Value(s).Slice(0, n))
}

func (s _slice) TakeWhile(f interface{}) Traversable {
	n := 0
	fw := funcOf(f)
	sval := reflect.Value(s)
	len := sval.Len()

	for i := 0; i < len; i++ {
		if !fw.call(sval.Index(i)).Bool() {
			break
		}
		n = i
	}
	if n > 0 {
		n++
	}

	return newSlice(sval.Slice(0, n))
}

// ----------------------------------------------------------------------------

// SliceOf returns a monadgo Slice value.
func SliceOf(x interface{}) Slice {
	xval := reflect.ValueOf(x)
	if xval.Kind() != reflect.Slice && xval.Kind() != reflect.Array {
		panic("x must be a slice or array")
	}

	if xval.Kind() == reflect.Array {
		// clone to a slice if x is array, or some slice operation on array will panic.
		yval := makeSlice(xval.Type().Elem(), xval.Len(), xval.Len())
		reflect.Copy(yval, xval)
		xval = yval
	}

	return newSlice(xval)
}

func newSlice(v reflect.Value) Slice {
	return _slice(v)
}
