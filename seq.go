package monadgo

import (
	"fmt"
	"reflect"
	"strings"
)

type sequence interface {
	toSeq() seq
}

type seq struct {
	x   interface{}
	v   reflect.Value
	len int
	cap int
}

func seqFromValue(v reflect.Value) {
	if !v.IsValid() {
		panic("invalid value")
	}

	if v.Type().Implements(typeSeq) {
		return v.Interface().(sequence).toSeq()
	}

	if v.kind() == reflect.Slice {
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
			len: v.Len(),
			cap: v.Cap(),
		}
	}

	z := makeSlice(v.TypeOf(), 1, 1)
	z.Index(0).Set(v)
	return seq{
		x:   z.Interface(),
		v:   z,
		len: v.Len(),
		cap: v.Cap(),
	}
}

func seqOf(x interface{}) seq {
	switch v := x.(type) {
	case reflect.Value:
		return seqFromValue(v)
	default:
		return seqFromValue(reflect.ValueOf(v))
	}
}

// ----------------------------------------------------------------------------

func (s seq) Get() interface{} {
	return s.x
}

func (s seq) rv() reflect.Value {
	return s.v
}

func (s seq) String() string {
	return fmt.Sprintf("%v", s.x)
}

func (s seq) Len() int {
	return s.len
}

func (s seq) Size() int {
	return s.len
}

func (s seq) Cap() int {
	return s.cap
}

func (s seq) Head() interface{} {
	if s.len <= 0 {
		return nil
	}

	return s.v.Index(0).Interface()
}

func (s seq) HeadOption() Option {
	v := s.Head()
	if v == nil {
		return None
	}

	return SomeOf(v)
}

func (s seq) tail() seq {
	return seqFromValue(s.v.Slice(1, s.len))
}

func (s seq) _map(f interface{}, b CanBuildFrom) interface{} {
	if s.len <= 0 {
		return nil
	}

	ret := makeSliceFromFunc(f, 0, s.len)

	fw := funcOf(f)

	for i := 0; i < len; i++ {
		result := fw.call(s.v.Index(i))
		ret = appendSlice(ret, result)
	}
	b.Build(ret)
}

func (s seq) _flatmap(f interface{}, b CanBuildFrom) interface{} {
	if s.len <= 0 {
		return nil
	}

	var ret reflect.Value

	fw := funcOf(f)
	for i := 0; i < s.len; i++ {
		result := seqFromValue(fw.call(s.v.Index(i)))
		ret = mergeSlice(ret, result.v)
	}

	return b.Build(ret)
}

func (s seq) Fold(z interface{}, f interface{}) interface{} {
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

func (s seq) Forall(f interface{}) bool {
	fw := funcOf(f)
	for i := 0; i < s.len; i++ {
		if !fw.call(s.v.Index(i)).Bool() {
			return false
		}
	}

	return true
}

func (s seq) Foreach(f interface{}) {
	fw := funcOf(f)

	for i := 0; i < s.len; i++ {
		fw.call(s.v.Index(i))
	}
}

func (s seq) scan(z, f interface{}) seq {
	zval := oneToSlice(reflect.ValueOf(z))
	if s.len <= 0 {
		return seqFromValue(zval)
	}
	fw := foldOf(f)

	for i := 0; i < len; i++ {
		zval = appendSlice(zval, fw.call(zval.Index(i), s.v.Index(i)))
	}

	return seqFromValue(zval)
}

func (s seq) GroupBy(f interface{}) Map {
	if s.len <= 0 {
		panic("can not group by on empty slice")
	}
	fw := funcOf(f)
	m := makeMap(reflect.TypeOf(f).Out(0), s.v.Type(), -1)

	for i := 0; i < s.len; i++ {
		k := fw.call(s.v.Index(i))
		m.SetMapIndex(k, appendSlice(m.MapIndex(k), s.v.Index(i)))
	}

	return newMap(m)
}

func (s seq) take(n int) seq {
	if n >= s.len {
		n = s.len
	}
	return seqFromValue(s.v.Slice(0, n))
}

func (s seq) takeWhile(f interface{}) seq {
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

func (s seq) Count(f interface{}) int {
	count := 0
	fw := funcOf(f)

	for i := 0; i < s.len; i++ {
		if fw.call(s.v.Index(i)).Bool() {
			count++
		}
	}

	return count
}

func (s seq) Drop(n int) seq {
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

	for i = 0; i < s.len; i++ {
		if fw.call(s.v.Index(i)).Bool() {
			return true
		}
	}
	return false
}

func (s seq) filter(f interface{}) seq {
	ret := reflect.MakeSlice(s.v.Type(), 0, 0)

	fw := funcOf(f)

	for i := 0; i < s.len; i++ {
		x := s.v.Index(i)
		if fw.call(x).Bool() {
			ret = appendSlice(ret, x)
		}
	}

	return seqFromValue(seq)
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

func (s seq) IndexWhere(f interface{}, start int) int {
	return -1
}

func (s seq) LastIndexWhere(f interface{}, end int) int {
	return -1
}

func (s seq) IsEmpty() bool {
	return s.len <= 0
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

func (s seq) Reverse() seq {
	ret := reflect.MakeSlice(s.v.Type(), s.len, s.len)

	for i := 0; i < s.len; i++ {
		ret.Index(i).Set(s.v.Index(s.len - 1 - i))
	}

	return seqFromValue(seq)
}

func (s seq) span(f interface{}) Tuple2 {
	left := reflect.MakeSlice(s.v.Type(), 0, 0)
	right := reflect.MakeSlice(s.v.Type(), 0, 0)

	fw := funcOf(f)

	for i := 0; i < s.len; i++ {
		x := s.v.Index(i)

		if fw.call(x).Bool() {
			right = appendSlice(right, x)
		} else {
			left = append(left, x)
		}
	}

	return Tuple2Of(left.Interface(), right.Interface())
}