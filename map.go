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

type _map struct {
	ktype reflect.Type
	vtype reflect.Type
	v     reflect.Value
}

var _ Map = _map{}

func mapCBF(k, v reflect.Type, x interface{}) Traversable {
	xtyp := reflect.TypeOf(x)
	if xtyp.Kind() == reflect.Slice {
		return mapCBF(k, v, seqOf(x))
	}

	t := x.(Traversable)
	if t.Size() <= 0 {
		return _map{
			ktype: k,
			vtype: v,
			v:     makeMap(k, v, 0),
		}
	}

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

		return newMap(ret)
	}

	if v.Type().Implements(typePair) {
		return newMap(oneToMap(v))
	}

	if v.Type().Implements(typeSeq) {
		return mapFromValue(v.Interface().(sequence).rv())
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
	t := v.Type()
	return _map{
		ktype: t.Key(),
		vtype: t.Elem(),
		v:     v,
	}
}

// ----------------------------------------------------------------------------

func (m _map) mapCBF(x interface{}) Traversable {
	return mapCBF(m.ktype, m.vtype, x)
}

func (m _map) Get() interface{} {
	return m.v.Interface()
}

func (m _map) rv() reflect.Value {
	return m.v
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
	return m.v.Len()
}

func (m _map) Range() *PairIter {
	return newPairIter(m.v)
}

func (m _map) Map(f interface{}) Traversable {
	ret := m.toSeq().Map(f)

	return m.mapCBF(ret)
}

func (m _map) FlatMap(f interface{}) Traversable {
	ret := m.toSeq().FlatMap(f)
	return m.mapCBF(ret)
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

func (m _map) Reduce(f interface{}) interface{} {
	return m.toSeq().Reduce(f)
}

func (m _map) GroupBy(f interface{}) Map {
	x := m.toSeq().GroupBy(f)
	x2 := x.Map(func(p Pair) Pair {
		return PairOf(p.Key(), MapOf(p.Value()).Get())
	})

	return x2.(Map)
}

func (m _map) Exists(f interface{}) bool {
	return m.toSeq().Exists(f)
}

func (m _map) Find(f interface{}) Option {
	return m.toSeq().Find(f)
}

func (m _map) Filter(f interface{}) Traversable {
	return m.mapCBF(m.toSeq().Filter(f))
}

func (m _map) MkString(start, sep, end string) string {
	return m.toSeq().MkString(start, sep, end)
}

func (m _map) Span(f interface{}) Tuple2 {
	t2 := m.toSeq().Span(f)
	return Tuple2Of(
		m.mapCBF(t2.V1()).Get(),
		m.mapCBF(t2.V2()).Get(),
	)
}
