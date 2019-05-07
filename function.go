package monadgo

import (
	"reflect"
)

type funcTR struct {
	in  [2]reflect.Type
	out [1]reflect.Type
	x   reflect.Value
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

// funcOf wraps original function f to one-input and one-output function.
// Input may be Unit if no input from f, original input, or tuple binding inputs from f.
// Output may be Unit if no output from f, Null if nil returns, or tuple binding outpus from f.
func funcOf(f interface{}) funcTR {
	ftyp := reflect.TypeOf(f)
	fval := reflect.ValueOf(f)

	switch ftyp.NumIn() {
	case 0: // input unit
		resultF := reflect.FuncOf([]reflect.Type{typeAny}, typeValues, false)
		return funcTR{
			in:  [2]reflect.Type{typeAny},
			out: [1]reflect.Type{bindType(ftyp)},
			x: reflect.MakeFunc(resultF, func(args []reflect.Value) []reflect.Value {
				out := reflect.ValueOf(bindValues(ftyp, fval.Call(nil)))
				return []reflect.Value{reflect.ValueOf(out)}
			})}
	case 1: // input original value from f.
		resultF := reflect.FuncOf([]reflect.Type{ftyp.In(0)}, typeValues, false)
		return funcTR{
			in:  [2]reflect.Type{ftyp.In(0)},
			out: [1]reflect.Type{bindType(ftyp)},
			x: reflect.MakeFunc(resultF, func(args []reflect.Value) []reflect.Value {
				out := reflect.ValueOf(bindValues(ftyp, fval.Call(args)))
				return []reflect.Value{reflect.ValueOf(out)}
			})}
	default: // bind all input from f to tuple.
		resultF := reflect.FuncOf([]reflect.Type{typeTuple}, typeValues, false)
		return funcTR{
			in:  [2]reflect.Type{typeTuple},
			out: [1]reflect.Type{bindType(ftyp)},
			x: reflect.MakeFunc(resultF, func(args []reflect.Value) []reflect.Value {
				t := args[0].Interface().(Tuple)

				out := reflect.ValueOf(bindValues(ftyp, fval.Call(t.toValues())))
				return []reflect.Value{reflect.ValueOf(out)}
			})}
	}
}

// foldOf wraps original function f to two-input and one-output function.
// Output may be Unit if no output from f, Null if nil returns, or tuple binding outpus from f.
func foldOf(f interface{}) funcTR {
	ftyp := reflect.TypeOf(f)

	if ftyp.NumIn() < 1 {
		panic("fold function must have one argument at last.")
	}

	fval := reflect.ValueOf(f)

	switch ftyp.NumIn() {
	case 1: // input (z) from f and unit.
		resultF := reflect.FuncOf([]reflect.Type{ftyp.In(0), typeAny}, typeValues, false)
		return funcTR{
			in:  [2]reflect.Type{ftyp.In(0), typeAny},
			out: [1]reflect.Type{bindType(ftyp)},
			x: reflect.MakeFunc(resultF, func(args []reflect.Value) []reflect.Value {
				out := reflect.ValueOf(bindValues(ftyp, fval.Call(args[0:1])))
				return []reflect.Value{reflect.ValueOf(out)}
			})}
	case 2: // inputs (z, x) from f.
		resultF := reflect.FuncOf([]reflect.Type{ftyp.In(0), ftyp.In(1)}, typeValues, false)
		return funcTR{
			in:  [2]reflect.Type{ftyp.In(0), ftyp.In(1)},
			out: [1]reflect.Type{bindType(ftyp)},
			x: reflect.MakeFunc(resultF, func(args []reflect.Value) []reflect.Value {
				out := reflect.ValueOf(bindValues(ftyp, fval.Call(args)))
				return []reflect.Value{reflect.ValueOf(out)}
			})}
	default:
		if ftyp.NumIn()&1 == 1 {
			// num of input is odd.
			resultF := reflect.FuncOf([]reflect.Type{ftyp.In(0), typeTuple}, typeValues, false)
			return funcTR{
				in:  [2]reflect.Type{ftyp.In(0), typeTuple},
				out: [1]reflect.Type{bindType(ftyp)},
				x: reflect.MakeFunc(resultF, func(args []reflect.Value) []reflect.Value {
					t := args[1].Interface().(Tuple)
					vals := append(args[0:1], t.toValues()...)
					out := reflect.ValueOf(bindValues(ftyp, fval.Call(vals)))
					return []reflect.Value{reflect.ValueOf(out)}
				})}
		}
		// num of input is even.
		resultF := reflect.FuncOf([]reflect.Type{typeTuple, typeTuple}, typeValues, false)
		return funcTR{
			in:  [2]reflect.Type{typeTuple, typeTuple},
			out: [1]reflect.Type{bindType(ftyp)},
			x: reflect.MakeFunc(resultF, func(args []reflect.Value) []reflect.Value {
				t0 := args[0].Interface().(Tuple)
				t1 := args[1].Interface().(Tuple)
				vals := append(t0.toValues(), t1.toValues()...)
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
		return unit
	case 1:
		if outs[0].Interface() == nil {
			return null
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

func bindType(ftyp reflect.Type) reflect.Type {

	switch ftyp.NumOut() {
	case 0:
		return typeUnit
	case 1:
		return ftyp.Out(0)
	case 2:
		return typePair
	case 3:
		return typeTuple3
	case 4:
		return typeTuple4
	default:
		return typeTuple
	}
}
