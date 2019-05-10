package monadgo

import (
	"fmt"
	"reflect"
)

// Map respresents a scala-like Map.
type Map interface {

	// Range returns a pair iterator.
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

	if t.rv().Type().Elem().ConvertibleTo(typePair) {
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

	if v.Kind() == reflect.Slice && v.Type().Elem().ConvertibleTo(typePair) {

		var ret reflect.Value

		len := v.Len()
		for i := 0; i < len; i++ {
			ret = mergeMap(ret, v.Index(i))
		}

		return newMap(ret)
	}

	if v.Type().ConvertibleTo(typePair) {
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

// toSeq returns a seq with element type Pair.
func (m _map) toSeq() seq {
	ret := makeSlice(typePair, 0, m.Size())

	it := m.Range()
	for it.Next() {
		ret = appendSlice(ret, reflect.ValueOf(it.Pair()))
	}
	return seqFromValue(ret)
}

// Size returns the size.
func (m _map) Size() int {
	return m.v.Len()
}

// Range returns a pair iterator.
func (m _map) Range() *PairIter {
	return newPairIter(m.v)
}

// Map applies function f to all elements in map m.
// f: func(Pair) X or func(K,V) X. X can be Pair or others.
// returns a Map if X is Pair.
func (m _map) Map(f interface{}) Traversable {
	ret := m.toSeq().Map(f)

	return m.mapCBF(ret)
}

// FlatMap applies f to all elements and builds a new Travesable from result.
// f: func(Pair) X or func(K,V) X, X can be Go slice, or map.
// returns a Map if X is a Go slice with element type Pair.
func (m _map) FlatMap(f interface{}) Traversable {
	ret := m.toSeq().FlatMap(f)
	return m.mapCBF(ret)
}

// Fold folds the elements using specified associative binary operator.
// z: func() T or value of type T. T can be monadgo Pair or Go tuple (K,V).
// f: func(T, T) T
// returns value with type T
func (m _map) Fold(z, f interface{}) interface{} {
	return m.toSeq().Fold(z, f)
}

// Foreach applies f to all element.
// f: func(T). T can be monadgo Pair or Go tuple (K,V).
func (m _map) Foreach(f interface{}) {
	m.toSeq().Foreach(f)
}

// Forall tests whether a predicate holds for all elements.
// f: func(T) bool. T can be monadgo Pair or Go tuple (K,V).
func (m _map) Forall(f interface{}) bool {
	return m.toSeq().Forall(f)
}

// Reduce reduces the elements of this using the specified associative binary operator.
// f: func(T, T) T. T can be monadgo Pair or Go tuple (K,V).
// returns Pair.
func (m _map) Reduce(f interface{}) interface{} {
	return m.toSeq().Reduce(f)
}

// GroupBy returns Map with K -> Go map. Key is the result of f. Collect elements into a map with same resulting key value.
// f: func(T) X. T can be monadgo Pair or Go tuple (K,V).
// returns Map(X -> Go map[K]V)
func (m _map) GroupBy(f interface{}) Map {
	x := m.toSeq().GroupBy(f)
	x2 := x.Map(func(p Pair) Pair {
		return PairOf(p.Key(), MapOf(p.Value()).Get())
	})

	return x2.(Map)
}

// Exists tests whether a predicate holds for at least one element of this sequence.
// f: func(T) bool. T can be monadgo Pair or Go tuple (K,V).
func (m _map) Exists(f interface{}) bool {
	return m.toSeq().Exists(f)
}

// Find returns the first pair satisfying f,
// otherwise return None.
// f: func(T) bool. T can be monadgo Pair or Go tuple (K,V).
func (m _map) Find(f interface{}) Option {
	return m.toSeq().Find(f)
}

// Filter retuns all elements satisfying f.
// f: func(T) bool. T can be monadgo Pair or Go tuple (K,V).
func (m _map) Filter(f interface{}) Traversable {
	return m.mapCBF(m.toSeq().Filter(f))
}

// MkString displays all elements in a string using start, end, and separator sep.
func (m _map) MkString(start, sep, end string) string {
	return m.toSeq().MkString(start, sep, end)
}

// Split splits this into a unsatisfying and satisfying pair according to f.
// f: func(T) bool. T can be monadgo Pair or Go tuple (K,V).
func (m _map) Split(f interface{}) Tuple2 {
	t2 := m.toSeq().Split(f)
	return Tuple2Of(
		m.mapCBF(t2.V1()).Get(),
		m.mapCBF(t2.V2()).Get(),
	)
}

// Collect returns elements satisfying pf.
// pf is a partial function consisting of Condition func(T) bool and Action func(T) X. T can be monadgo Pair or Go tuple (K,V).
// returns a Map if type of X is Pair.
func (m _map) Collect(pf PartialFunc) Traversable {
	return m.mapCBF(m.toSeq().Collect(pf))
}
