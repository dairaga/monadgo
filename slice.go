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

type _slice struct {
	seq
}

var _ Slice = _slice{}

var sliceBuilder CanBuildFrom = func(v reflect.Value) reflect.Value {
	return reflect.ValueOf(SliceOf(v))
}

// ----------------------------------------------------------------------------

func (s _slice) toSeq() seq {
	return s.seq
}

func (s _slice) String() string {
	return fmt.Sprintf("%v", s.Get())
}

func (s _slice) Tail() Traversable {
	return SliceOf(s.seq.tail())
}

func (s _slice) Map(f interface{}) Traversable {
	return s.seq._map(f, sliceBuilder).(Traversable)
}

func (s _slice) FlatMap(f interface{}) Traversable {
	return s.seq._flatmap(f, sliceBuilder).(Traversable)
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

func (s _slice) Take(n int) Traversable {
	return SliceOf(s.seq.take(n))
}

func (s _slice) TakeWhile(f interface{}) Traversable {
	return SliceOf(s.seq.takeWhile(f))
}

// ----------------------------------------------------------------------------

// SliceOf returns a monadgo Slice value.
func SliceOf(x interface{}) Slice {
	return _slice{seqOf(x)}
}

func newSlice(v reflect.Value) Slice {
	return _slice(v)
}
