package monadgo

import (
	"fmt"
	"reflect"
)

// Map respresents a scala-like Map.
type Map interface {
	//Any
	Size() int
	Range() *PairIter
	Map(f interface{}) Traversable
	FlatMap(f interface{}) Traversable

	Traversable
}

type _map reflect.Value

var _ Map = _map{}

// ----------------------------------------------------------------------------

func (m _map) Get() interface{} {
	return reflect.Value(m).Interface()
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
	ftyp := reflect.TypeOf(f)
	isMap := ftyp.NumOut() == 2 || (ftyp.NumOut() == 1 && ftyp.Out(0).Implements(typePair))
	var ret reflect.Value

	fw := funcOf(f)

	it := m.Range()
	for it.Next() {
		result := fw.call(reflect.ValueOf(it.Pair()))
		if isMap {
			ret = mergeMap(ret, result)
		} else {
			ret = appendSlice(ret, result)
		}
	}

	return TraversableOf(ret)
}

func (m _map) FlatMap(f interface{}) Traversable {
	fout := reflect.TypeOf(f).Out(0)
	var ret reflect.Value
	fw := funcOf(f)

	it := m.Range()
	for it.Next() {
		result := fw.call(reflect.ValueOf(it.Pair()))
		if fout.Kind() == reflect.Map {
			ret = mergeMap(ret, result)
		} else if fout.Kind() == reflect.Slice && fout.Elem().Implements(typePair) {
			len := result.Len()
			for i := 0; i < len; i++ {
				ret = mergeMap(ret, result.Index(i))
			}
		} else {
			ret = mergeSlice(ret, result)
		}
	}

	return TraversableOf(ret)
}

func (m _map) Fold(z interface{}, f interface{}) interface{} {
	it := m.Range()
	fw := foldOf(f)

	for it.Next() {
		z = fw.fold(z, it.Pair())
	}

	return z
}

func (m _map) Foreach(f interface{}) {
	it := m.Range()
	fw := funcOf(f)
	for it.Next() {
		fw.invoke(it.Pair())
	}
}

func (m _map) Forall(f interface{}) bool {
	it := m.Range()
	fw := funcOf(f)
	for it.Next() {
		if !fw.invoke(it.Pair()).(bool) {
			return false
		}
	}
	return true
}

func (m _map) ToSeq() interface{} {
	size := m.Size()
	ret := make([]Pair, 0, size)

	it := m.Range()

	for it.Next() {
		ret = append(ret, it.Pair())
	}

	return ret
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
