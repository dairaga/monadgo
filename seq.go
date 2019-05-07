package monadgo

import (
	"fmt"
	"reflect"
	"strings"
)

type sequence interface {
	Any
	toSeq() seq
}

type seq struct {
	x     interface{}
	t     reflect.Type
	v     reflect.Value
	len   int
	cap   int
	empty bool
}

var emptySeq = seq{
	x:     nothings,
	t:     typeNothings,
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
			t:   v.Type(),
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
			t:   z.Type(),
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
	if s.len <= 0 {
		return
	}

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

func (s seq) Head() interface{} {
	if s.len <= 0 {
		return nil
	}

	return s.v.Index(0).Interface()
}

func (s seq) HeadOption() Option {
	if s.len <= 0 {
		return None
	}
	return OptionOf(s.v.Index(0))
}

func (s seq) Tail() Traversable {
	if s.len < 1 {
		return seqFromValue(s.v.Slice(0, 0))
	}

	return seqFromValue(s.v.Slice(1, s.len))
}

func (s seq) Reduce(f interface{}) interface{} {
	if s.len <= 0 {
		panic("empty list can not reduce")
	}

	if s.len == 1 {
		return s.v.Index(0).Interface()
	}

	fw := foldOf(f)
	zval := s.v.Index(0)

	for i := 1; i < s.len; i++ {
		zval = fw.call(zval, s.v.Index(i))
	}

	return zval.Interface()
}

func (s seq) Scan(z, f interface{}) Traversable {
	z = checkAndInvoke(z)
	zval := oneToSlice(reflect.ValueOf(z))
	if s.len <= 0 {
		return seqFromValue(zval)
	}
	fw := foldOf(f)

	for i := 0; i < s.len; i++ {
		zval = appendSlice(zval, fw.call(zval.Index(i), s.v.Index(i)))
	}

	return seqFromValue(zval)
}

func (s seq) GroupBy(f interface{}) Map {
	if s.len <= 0 {
		panic("can not group by on empty slice")
	}
	fw := funcOf(f)
	m := makeMap(fw.out[0], s.t, -1)

	for i := 0; i < s.len; i++ {
		k := fw.call(s.v.Index(i))
		m.SetMapIndex(k, appendSlice(m.MapIndex(k), s.v.Index(i)))
	}

	return newMap(m)
}

func (s seq) Take(n int) Traversable {
	if n >= s.len {
		n = s.len
	}
	return seqFromValue(s.v.Slice(0, n))
}

func (s seq) TakeWhile(f interface{}) Traversable {
	n := 0
	fw := funcOf(f)

	for i := 0; i < s.len; i++ {
		if !fw.call(s.v.Index(i)).Bool() {
			break
		}
		n = i
	}
	if n > 0 {
		n++
	}

	return seqFromValue(s.v.Slice(0, n))
}

func (s seq) Drop(n int) Traversable {
	if n > s.len {
		n = s.len
	}

	return seqFromValue(s.v.Slice(n, s.len))
}

func (s seq) Exists(f interface{}) bool {
	if s.len <= 0 {
		return false
	}

	fw := funcOf(f)

	for i := 0; i < s.len; i++ {
		if fw.call(s.v.Index(i)).Bool() {
			return true
		}
	}
	return false
}

func (s seq) Find(f interface{}) Option {
	if s.len <= 0 {
		return None
	}

	fw := funcOf(f)

	for i := 0; i < s.len; i++ {
		x := s.v.Index(i)
		if fw.call(x).Bool() {
			return SomeOf(x)
		}
	}

	return None
}

func (s seq) Filter(f interface{}) Traversable {
	ret := reflect.MakeSlice(s.t, 0, 0)

	fw := funcOf(f)

	for i := 0; i < s.len; i++ {
		x := s.v.Index(i)
		if fw.call(x).Bool() {
			ret = appendSlice(ret, x)
		}
	}

	return seqFromValue(ret)
}

func (s seq) IndexWhere(f interface{}, start int) int {
	if s.len <= 0 {
		return -1
	}
	if start < 0 {
		start = 0
	}

	fw := funcOf(f)
	for i := start; i < s.len; i++ {
		if fw.call(s.v.Index(i)).Bool() {
			return i
		}
	}

	return -1
}

func (s seq) LastIndexWhere(f interface{}, end int) int {
	if s.len <= 0 {
		return -1
	}
	if end >= s.len {
		end = s.len - 1
	}

	fw := funcOf(f)
	for i := end; i >= 0; i-- {
		if fw.call(s.v.Index(i)).Bool() {
			return i
		}
	}
	return -1
}

func (s seq) MkString(start, sep, end string) string {
	sb := new(strings.Builder)
	sb.WriteString(start)
	for i := 0; i < s.len; i++ {
		sb.WriteString(fmt.Sprintf("%v", s.v.Index(i).Interface()))
		sb.WriteString(sep)
	}
	sb.WriteString(end)

	return sb.String()
}

func (s seq) Reverse() Traversable {
	ret := reflect.MakeSlice(s.t, s.len, s.len)

	for i := 0; i < s.len; i++ {
		ret.Index(i).Set(s.v.Index(s.len - 1 - i))
	}

	return seqFromValue(ret)
}

func (s seq) Span(f interface{}) Tuple2 {
	if s.len <= 0 {
		return Tuple2Of(nothings, nothings)
	}

	left := reflect.MakeSlice(s.t, 0, 0)
	right := reflect.MakeSlice(s.t, 0, 0)

	fw := funcOf(f)

	for i := 0; i < s.len; i++ {
		x := s.v.Index(i)

		if fw.call(x).Bool() {
			right = appendSlice(right, x)
		} else {
			left = appendSlice(left, x)
		}
	}

	return Tuple2Of(left.Interface(), right.Interface())
}

func (s seq) Collect(pf PartialFunc) Traversable {
	ret := makeSlice(pf.action.out[0], 0, 0)
	if s.len <= 0 {
		return seqFromValue(ret)
	}

	for i := 0; i < s.len; i++ {
		result := pf.Call(s.v.Index(i))
		if result != nothingValue {
			ret = appendSlice(ret, result)
		}
	}

	return seqFromValue(ret)
}
