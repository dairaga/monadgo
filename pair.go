package monadgo

import (
	"reflect"
)

/*
// Pair represents a scala-like Pair.
type Pair interface {
	Tuple2
	Key() interface{}
	Value() interface{}
}
*/

// Pair represents a scala-like Pair.
type Pair struct {
	Tuple2
}

func (p Pair) Key() interface{} {
	return p.Tuple2.V1()
}

func (p Pair) Value() interface{} {
	return p.Tuple2.V2()
}

// ----------------------------------------------------------------------------

// PairOf returns a pair.
func PairOf(k, v interface{}) Pair {
	return Pair{newTuple2(reflect.TypeOf(k), reflect.TypeOf(v), reflect.ValueOf(k), reflect.ValueOf(v))}
}

func pairFromTuple2(t Tuple2) Pair {
	return Pair{t}
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
