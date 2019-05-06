package monadgo

import (
	"reflect"
)

type container struct {
	x interface{}
	v reflect.Value
}

func containerFromValue(v reflect.Value) container {
	if !v.IsValid() {
		panic("invalid value")
	}

	return container{
		x: v.Interface(),
		v: v,
	}
}

func containerOf(x interface{}) container {
	if x == nil {
		return container{
			x: null,
			v: nullValue,
		}
	}

	switch v := x.(type) {
	case reflect.Value:
		return containerFromValue(v)
	default:
		return containerFromValue(reflect.ValueOf(x))
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

func (c container) scan(z, f interface{}) seq {
	zval := reflect.ValueOf(z)
	s := makeSlice(c.v.Type(), 2, 2)
	s.Index(0).Set(zval)
	s.Index(1).Set(foldOf(f).call(zval, c.v))

	return seqFromValue(s)
}
