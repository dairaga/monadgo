package monadgo

import (
	"reflect"
)

type container struct {
	x interface{}
	v reflect.Value
}

func genContainer(x interface{}) container {
	return container{
		x: x,
		v: reflect.ValueOf(x),
	}
}

func (c container) Get() interface{} {
	return c.x
}

func (c container) rv() reflect.Value {
	return c.v
}

func (c container) invoke(f interface{}) reflect.Value {
	return funcOf(f).call(c.v)
}

func (c container) _map(f interface{}, b CanBuildFrom) interface{} {
	return b.Build(c.invoke(f))
}

func (c container) FlatMap(f interface{}) interface{} {
	return c.invoke(f).Interface()
}

func (c container) Forall(f interface{}) bool {
	return c.invoke(f).Bool()
}

func (c container) Foreach(f interface{}) {
	c.invoke(f)
}

func (c container) Reduce(interface{}) interface{} {
	return c.x
}

func (c container) Scan(z, f interface{}) Traversable {
	zval := reflect.ValueOf(z)
	s := makeSlice(c.v.Type(), 2, 2)
	s.Index(0).Set(zval)
	s.Index(1).Set(foldOf(f).call(zval, c.v))

	return SliceOf()

}

/*
	GroupBy(f interface{}) Map

	Take(n int) Traversable

	TakeWhile(f interface{}) Traversable


		Collect(f interface{}) Traversable

		Count(f interface{}) int

		Drop(n int) Traversable

		Exists(f interface{}) Traversable

		Filter(f interface{}) Traversable

		Find(f interface{}) Option

		IndexWhere(f interface{}, start int) int

		LastIndexWhere(f interface{}, end int) int

		IsEmpty() bool

		MaxBy(f interface{}) interface{}

		MinBy(f interface{}) interface{}

		MkString(start, sep, end string) string

		Reverse() Traversable

		Span(f interface{}) Pair // PairOf Traversable
*/
