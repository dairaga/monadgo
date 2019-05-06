package monadgo

import (
	"reflect"
)

type funcTR struct {
	x reflect.Value
}

func (f funcTR) call(v ...reflect.Value) reflect.Value {
	return f.x.Call(v)[0].Interface().(reflect.Value)
}

func (f funcTR) invoke(v interface{}) interface{} {
	return f.call(reflect.ValueOf(v)).Interface()
}

func (f funcTR) fold(z, v interface{}) interface{} {
	return f.call(reflect.ValueOf(z), reflect.ValueOf(v)).Interface()
}

// funcOf ...
func funcOf(f interface{}) funcTR {
	ftyp := reflect.TypeOf(f)
	fval := reflect.ValueOf(f)

	switch ftyp.NumIn() {
	case 0:
		resultF := reflect.FuncOf([]reflect.Type{typeUnit}, []reflect.Type{typeValue}, false)
		return funcTR{reflect.MakeFunc(resultF, func(args []reflect.Value) []reflect.Value {
			out := reflect.ValueOf(bindValues(ftyp, fval.Call(nil)))
			return []reflect.Value{reflect.ValueOf(out)}
		})}
	case 1:
		resultF := reflect.FuncOf([]reflect.Type{ftyp.In(0)}, []reflect.Type{typeValue}, false)
		return funcTR{reflect.MakeFunc(resultF, func(args []reflect.Value) []reflect.Value {
			out := reflect.ValueOf(bindValues(ftyp, fval.Call(args)))
			return []reflect.Value{reflect.ValueOf(out)}
		})}
	default:
		resultF := reflect.FuncOf([]reflect.Type{typeTuple}, []reflect.Type{typeValue}, false)
		return funcTR{reflect.MakeFunc(resultF, func(args []reflect.Value) []reflect.Value {
			t := args[0].Interface().(Tuple)

			out := reflect.ValueOf(bindValues(ftyp, fval.Call(t.toValues())))
			return []reflect.Value{reflect.ValueOf(out)}
		})}
	}
}

func foldOf(f interface{}) funcTR {
	ftyp := reflect.TypeOf(f)

	if ftyp.NumIn() < 1 {
		panic("fold function must have one argument at last.")
	}

	fval := reflect.ValueOf(f)

	switch ftyp.NumIn() {
	case 1:
		resultF := reflect.FuncOf([]reflect.Type{ftyp.In(0), typeUnit}, []reflect.Type{typeValue}, false)
		return funcTR{reflect.MakeFunc(resultF, func(args []reflect.Value) []reflect.Value {
			out := reflect.ValueOf(bindValues(ftyp, fval.Call(args[0:1])))
			return []reflect.Value{reflect.ValueOf(out)}
		})}
	case 2:
		resultF := reflect.FuncOf([]reflect.Type{ftyp.In(0), ftyp.In(1)}, []reflect.Type{typeValue}, false)
		return funcTR{reflect.MakeFunc(resultF, func(args []reflect.Value) []reflect.Value {
			out := reflect.ValueOf(bindValues(ftyp, fval.Call(args)))
			return []reflect.Value{reflect.ValueOf(out)}
		})}
	default:
		resultF := reflect.FuncOf([]reflect.Type{ftyp.In(0), typeTuple}, []reflect.Type{typeValue}, false)
		return funcTR{reflect.MakeFunc(resultF, func(args []reflect.Value) []reflect.Value {
			t := args[1].Interface().(Tuple)
			vals := append(args[0:1], t.toValues()...)
			out := reflect.ValueOf(bindValues(ftyp, fval.Call(vals)))
			return []reflect.Value{reflect.ValueOf(out)}
		})}
	}

}

// ----------------------------------------------------------------------------

// bindValues converts []reflect.Value to callable type.
// return Unit if output size is 0, native if size is 1, and Tuple if size more than 1.
func bindValues(ftyp reflect.Type, outs []reflect.Value) interface{} {

	switch len(outs) {
	case 0:
		return Unit
	case 1:
		if outs[0].Interface() == nil {
			return Null
		}
		return outs[0].Interface()
	case 2:
		return pairFromTuple2(newTuple2(ftyp.Out(0), ftyp.Out(1), outs[0], outs[1]))
	case 3:
		return newTuple3(ftyp.Out(0), ftyp.Out(1), ftyp.Out(2), outs[0], outs[1], outs[2])
	case 4:
		return newTuple4(ftyp.Out(0), ftyp.Out(1), ftyp.Out(2), ftyp.Out(3), outs[0], outs[1], outs[2], outs[3])
	default:
		types := make([]reflect.Type, ftyp.NumOut(), ftyp.NumOut())
		for i := 0; i < ftyp.NumOut(); i++ {
			types[i] = ftyp.Out(i)
		}
		return newTuple(types, outs)
	}
}

// checkFuncAndInvoke check f is a function or not.
// Return result from invoke f and convert output of f with bindValues and true if f is a function.
func checkFuncAndInvoke(f interface{}) (interface{}, bool) {
	ftyp := reflect.TypeOf(f)
	fval := reflect.ValueOf(f)

	if ftyp.Kind() == reflect.Func {
		return bindValues(ftyp, fval.Call(nil)), true
	}

	return nil, false
}
