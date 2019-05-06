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
		v = nothingValue
	}

	return container{
		x: v.Interface(),
		v: v,
	}

}

func containerOf(x interface{}) container {
	if x == nil {
		return container{
			x: Null,
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

func (c container) Forall(f interface{}) bool {
	return c.invoke(f).Bool()
}

func (c container) Foreach(f interface{}) {
	c.invoke(f)
}

func (c container) Exists(f interface{}) bool {
	return c.Forall(f)
}

func (c container) _map(f interface{}, b CanBuildFrom) interface{} {
	return b.Build(c.invoke(f)).Interface()
}

func (c container) flatMap(f interface{}) interface{} {
	return c.invoke(f).Interface()
}
