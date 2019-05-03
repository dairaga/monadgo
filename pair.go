package monadgo

import (
	"reflect"
)

// Pair represents a scala-like Pair.
type Pair interface {
	Tuple2
	Key() interface{}
	Value() interface{}
}

type _pair struct {
	_tuple2
}

var _ Pair = _pair{}

func (p _pair) Key() interface{} {
	return p._tuple2.V1()
}

func (p _pair) Value() interface{} {
	return p._tuple2.V2()
}

// ----------------------------------------------------------------------------

// PairOf returns a pair.
func PairOf(k, v interface{}) Pair {
	return _pair{newTuple2(reflect.TypeOf(k), reflect.TypeOf(v), reflect.ValueOf(k), reflect.ValueOf(v))}
}

func pairFromTuple2(t _tuple2) Pair {
	return _pair{t}
}

// ----------------------------------------------------------------------------

// PairIter represents a iterator of Pair for Map.
type PairIter struct {
	it *reflect.MapIter
}

// Next returns if iterator does not reach end.
func (pit *PairIter) Next() bool {
	return pit.it.Next()
}

// Pair returns current value.
func (pit *PairIter) Pair() Pair {
	return PairOf(pit.it.Key().Interface(), pit.it.Value().Interface())
}

func newPairIter(m reflect.Value) *PairIter {
	return &PairIter{m.MapRange()}
}
