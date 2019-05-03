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

func (s _slice) Cap() int {
	return reflect.Value(s).Cap()
}

func (s _slice) Map(f interface{}) Traversable {
	sval := reflect.Value(s)
	len := sval.Len()

	ret := makeSlice(reflect.TypeOf(f).Out(0), len)
	fw := funcOf(f)

	for i := 0; i < len; i++ {
		result := fw.call(sval.Index(i))
		ret.Index(i).Set(result)
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
		ret = mergeSlice(ret, reflect.ValueOf(seq.ToSeq()))
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

func (s _slice) ToSeq() interface{} {
	return s.Get()
}

// ----------------------------------------------------------------------------

// SliceOf returns a monadgo Slice value.
func SliceOf(x interface{}) Slice {
	xval := reflect.ValueOf(x)
	if xval.Kind() != reflect.Slice && xval.Kind() != reflect.Array {
		panic("x must be a slice")
	}

	return newSlice(xval)
}

func newSlice(v reflect.Value) Slice {
	return _slice(v)
}
