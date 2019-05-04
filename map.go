package monadgo

import (
	"fmt"
	"reflect"
)

// Map respresents a scala-like Map.
type Map interface {
	//Any
	//Size() int
	Range() *PairIter
	//Map(f interface{}) Traversable
	//FlatMap(f interface{}) Traversable

	Traversable
}

type _map reflect.Value

var _ Map = _map{}

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

func (m _map) Size() int {
	return reflect.Value(m).Len()
}

func (m _map) Range() *PairIter {
	return newPairIter(reflect.Value(m))
}

func (m _map) Map(f interface{}) Traversable {
	seq := m.ToSlice().Map(f)
	if m := seq.ToMap(); m != nil {
		return m
	}
	return seq
}

func (m _map) FlatMap(f interface{}) Traversable {
	seq := m.ToSlice().FlatMap(f)
	if m := seq.ToMap(); m != nil {
		return m
	}

	return seq
}

func (m _map) Fold(z interface{}, f interface{}) interface{} {
	return m.ToSlice().Fold(z, f)
}

func (m _map) Foreach(f interface{}) {
	m.ToSlice().Foreach(f)
}

func (m _map) Forall(f interface{}) bool {
	return m.ToSlice().Forall(f)
}

func (m _map) ToSlice() Slice {
	size := m.Size()
	ret := make([]Pair, 0, size)

	it := m.Range()

	for it.Next() {
		ret = append(ret, it.Pair())
	}

	return SliceOf(ret)
}

func (m _map) ToMap() Map {
	return m
}

func (m _map) Head() interface{} {
	it := m.Range()
	if it.Next() {
		return it.Pair()
	}
	return nil
}

func (m _map) Tail() Traversable {
	return m.ToSlice().Tail().ToMap()
}

func (m _map) Reduce(f interface{}) interface{} {
	return m.ToSlice().Reduce(f)
}

func (m _map) Scan(z, f interface{}) Traversable {
	return m.ToSlice().Scan(z, f).ToMap()
}

func (m _map) GroupBy(f interface{}) Map {
	len := m.Size()
	if len <= 0 {
		panic("can not group by on empty map")
	}

	mval := reflect.Value(m)
	ftyp := reflect.TypeOf(f)
	elm := mval.Type()
	fw := funcOf(f)
	newm := makeMap(ftyp.Out(0), elm, -1)

	it := m.Range()

	for it.Next() {
		p := it.Pair()
		k := reflect.ValueOf(fw.invoke(p))

		newm.SetMapIndex(k, mergeMap(newm.MapIndex(k), reflect.ValueOf(p)))
	}

	return newMap(newm)
}

// Take returns the first n elements in map.
func (m _map) Take(n int) Traversable {
	return m.ToSlice().Take(n).ToMap()
}

func (m _map) TakeWhile(f interface{}) Traversable {
	return m.ToSlice().TakeWhile(f).ToMap()
}

// ----------------------------------------------------------------------------

// MapOf returns a Map.
func MapOf(x interface{}) Map {
	xval := reflect.ValueOf(x)
	if xval.Kind() != reflect.Map {
		panic("x must be map")
	}

	return newMap(xval)
}

func newMap(v reflect.Value) Map {
	return _map(v)
}
