package monadgo

import "reflect"

// container wraps a GO value.
type container struct {
	x     interface{}
	v     reflect.Value
	empty bool
}

var nullContainer = container{
	x: null,
	v: nullValue,
}

var nothingContainer = container{
	x:     nothing,
	v:     nothingValue,
	empty: true,
}

func containerOf(x interface{}) container {
	if x == nil {
		return nullContainer
	}

	switch v := x.(type) {
	case reflect.Value:
		if !v.IsValid() {
			return nothingContainer
		}

		return container{
			x: v.Interface(),
			v: v,
		}
	default:
		return container{
			x: x,
			v: reflect.ValueOf(v),
		}
	}
}

// ----------------------------------------------------------------------------

func (c container) Get() interface{} {
	return c.x
}

func (c container) rv() reflect.Value {
	return c.v
}

func (c container) invoke(f interface{}) reflect.Value {
	return funcOf(f).call(c.v)
}

func (c container) _map(f interface{}) container {
	if c.empty {
		return c
	}

	return containerOf(c.invoke(f))
}

func (c container) _flatMap(f interface{}) interface{} {
	if c.empty {
		return c
	}
	return c.invoke(f).Interface()
}

func (c container) Forall(f interface{}) bool {
	if c.empty {
		return true
	}

	return c.invoke(f).Bool()
}

func (c container) Foreach(f interface{}) {
	if !c.empty {
		c.invoke(f)
	}
}

func (c container) Empty() bool {
	return c.empty
}
