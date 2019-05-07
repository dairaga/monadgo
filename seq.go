package monadgo

import (
	"fmt"
	"reflect"
)

type sequence interface {
	Any
	toSeq() seq
}

type seq struct {
	x     interface{}
	v     reflect.Value
	len   int
	cap   int
	empty bool
}

var emptySeq = seq{
	x:     nothings,
	v:     nothingsValue,
	len:   0,
	cap:   0,
	empty: true,
}

func seqFromValue(v reflect.Value) seq {
	if !v.IsValid() {
		return emptySeq
	}

	if v.Type().Implements(typeSeq) {
		return v.Interface().(sequence).toSeq()
	}

	if v.Kind() == reflect.Slice {
		return seq{
			x:   v.Interface(),
			v:   v,
			len: v.Len(),
			cap: v.Cap(),
		}
	}

	if v.Kind() == reflect.Array {
		// clone to a slice if x is array, or some slice operation on array will panic.
		z := makeSlice(v.Type().Elem(), v.Len(), v.Len())
		reflect.Copy(z, v)
		return seq{
			x:   z.Interface(),
			v:   z,
			len: z.Len(),
			cap: z.Cap(),
		}
	}

	if v.Kind() == reflect.Map {
		return newMap(v).toSeq()
	}

	z := makeSlice(v.Type(), 1, 1)
	z.Index(0).Set(v)
	return seqFromValue(z)

}

func seqOf(x interface{}) seq {
	if x == nil {
		return emptySeq
	}

	switch v := x.(type) {
	case reflect.Value:
		return seqFromValue(v)
	default:
		return seqFromValue(reflect.ValueOf(x))
	}
}

// ----------------------------------------------------------------------------

func (s seq) Get() interface{} {
	return s.x
}

func (s seq) String() string {
	return fmt.Sprintf("%v", s.x)
}

func (s seq) rv() reflect.Value {
	return s.v
}

func (s seq) toSeq() seq {
	return s
}

func (s seq) Size() int {
	return s.len
}

func (s seq) Len() int {
	return s.len
}

func (s seq) Cap() int {
	return s.cap
}

func (s seq) Map(f interface{}) Traversable {
	if s.empty {
		return emptySeq
	}

	fw := funcOf(f)
	ret := makeSlice(fw.out[0])

	for i := 0; i < s.len; i++ {
		ret = appendSlice(ret, fw.call(s.v.Index(i)))
	}
	return seqFromValue(ret)
}

func (s seq) FlatMap(f interface{}) Traversable {
	if s.empty {
		return emptySeq
	}
	fw := funcOf(f)
	var elm reflect.Type

	if fw.out[0].Kind() == reflect.Slice {
		elm = fw.out[0].Elem()
	} else if fw.out[0].Kind() == reflect.Map {
		elm = typePair
	} else {
		elm = fw.out[0]
	}

	ret := makeSlice(elm, 0, 0)

	for i := 0; i < s.len; i++ {
		result := fw.call(s.v.Index(i))
		ret = mergeSlice(ret, seqFromValue(result).v)
	}

	if elm == typePair {
		return mapFromValue(ret)
	}

	return SliceOf(ret)
}

// Forall tests whether a predicate holds for all elements.
func (s seq) Forall(f interface{}) bool {
	if s.empty {
		return false
	}

	fw := funcOf(f)

	for i := 0; i < s.len; i++ {
		if !fw.call(s.v.Index(i)).Bool() {
			return false
		}
	}

	return true
}

// Foreach applies f to all element.
func (s seq) Foreach(f interface{}) {
	fw := funcOf(f)

	for i := 0; i < s.len; i++ {
		fw.call(s.v.Index(i))
	}
}

// Fold folds the elements using specified associative binary operator.
func (s seq) Fold(z, f interface{}) interface{} {
	z = checkAndInvoke(z)

	if s.len <= 0 {
		return z
	}

	fw := foldOf(f)
	zval := reflect.ValueOf(z)

	for i := 0; i < s.len; i++ {
		zval = fw.call(zval, s.v.Index(i))
	}
	return zval.Interface()
}
