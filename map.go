package monadgo

import (
	"fmt"
	"reflect"
)

// Map respresents a scala-like Map.
type Map interface {
	Range() *PairIter
	Traversable
}

type _map reflect.Value

var _ Map = _map{}

func mapCBF(t Traversable) Traversable {
	if t.rv().Type().Elem().Implements(typePair) {
		var ret reflect.Value
		v := t.rv()
		len := t.Size()
		for i := 0; i < len; i++ {
			ret = mergeMap(ret, v.Index(i))
		}
		return newMap(ret)
	}

	return t
}

func mapFromValue(v reflect.Value) Map {
	if !v.IsValid() {
		panic("v is invalid")
	}

	if v.Kind() == reflect.Map {
		return newMap(v)
	}

	if v.Kind() == reflect.Slice && v.Type().Elem().Implements(typePair) {

		var ret reflect.Value

		len := v.Len()
		for i := 0; i < len; i++ {
			ret = mergeMap(ret, v.Index(i))
		}

		return newMap(v)
	}

	if v.Type().Implements(typePair) {
		p := v.Interface().(Pair)
		return newMap(makeMap(p.T1(), p.T2(), -1))
	}

	panic(fmt.Sprintf("%v can not convert to map", v.Interface()))

}

// MapOf returns a Map.
func MapOf(x interface{}) Map {
	if x == nil {
		panic("x is nil")
	}

	switch v := x.(type) {
	case reflect.Value:
		return mapFromValue(v)
	default:
		return mapFromValue(reflect.ValueOf(x))
	}
}

func newMap(v reflect.Value) Map {
	return _map(v)
}

// ----------------------------------------------------------------------------

func (m _map) Get() interface{} {
	return reflect.Value(m).Interface()
}

func (m _map) rv() reflect.Value {
	return reflect.Value(m)
}

func (m _map) String() string {
	return fmt.Sprintf("%v", m.Get())
}

func (m _map) toSeq() seq {
	ret := makeSlice(typePair, 0, m.Size())

	it := m.Range()
	for it.Next() {
		ret = appendSlice(ret, reflect.ValueOf(it.Pair()))
	}
	return seqFromValue(ret)
}

func (m _map) Size() int {
	return reflect.Value(m).Len()
}

func (m _map) Range() *PairIter {
	return newPairIter(reflect.Value(m))
}

func (m _map) Map(f interface{}) Traversable {
	ret := m.toSeq().Map(f)

	return mapCBF(ret)
}

func (m _map) FlatMap(f interface{}) Traversable {
	ret := m.toSeq().FlatMap(f)
	return mapCBF(ret)
}

func (m _map) Fold(z, f interface{}) interface{} {
	return m.toSeq().Fold(z, f)
}

func (m _map) Foreach(f interface{}) {
	m.toSeq().Foreach(f)
}

func (m _map) Forall(f interface{}) bool {
	return m.toSeq().Forall(f)
}

// ----------------------------------------------------------------------------
