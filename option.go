package monadgo

import (
	"fmt"
	"reflect"
)

// OptionOrElse is function type for Option.OrElse.
type OptionOrElse func() Option

// Option represents scala-like Option[T].
type Option interface {
	Any

	// Defined return true if this is Some, otherwise return false.
	Defined() bool

	// Map applies function f if this is Some.
	// f: func(T) X
	// returns Option[X]
	Map(f interface{}) Option

	// FlatMap binds the function f across Some.
	// f: func(T) Option.
	// returns a new Option.
	FlatMap(f interface{}) Option

	// Foreach executes the given side-effecting function f if this is a Some.
	// f: func(T)
	Foreach(f interface{})

	// Forall returns true if this option is empty,
	// or the predicate p returns true when applied to this Some's value.
	// f: func(T) bool
	Forall(f interface{}) bool

	// Fold returns the result of applying f to this Option's value if this is Some,
	// otherwise return z.
	// z: func() X or value with type X
	// f: func(T) X
	// returns value with type X
	Fold(z, f interface{}) interface{}

	// OrElse returns this if it is nonempty,
	// otherwise return the result from f.
	OrElse(f OptionOrElse) Option

	// GetOrElse returns the option's value if the option is Some,
	// otherwise return the result z.
	// z: func() X or value with type X
	// returns value with type X
	GetOrElse(z interface{}) interface{}
}

// ----------------------------------------------------------------------------

type traitOption struct {
	empty bool
	v     reflect.Value
}

func (o *traitOption) Get() interface{} {
	return o.v.Interface()
}

func (o *traitOption) rv() reflect.Value {
	return o.v
}

func (o *traitOption) String() string {
	if o.empty {
		return "None"
	}
	return fmt.Sprintf("Some(%v)", o.Get())
}

func (o *traitOption) Defined() bool {
	return !o.empty
}

func (o *traitOption) Map(f interface{}) Option {
	if !o.empty {
		return optionCBF(funcOf(f).call(o.v))
	}

	return None
}

func (o *traitOption) FlatMap(f interface{}) Option {
	if !o.empty {
		return funcOf(f).call(o.v).Interface().(Option)
	}

	return None
}

func (o *traitOption) Fold(z, f interface{}) interface{} {
	if o.empty {
		return checkAndInvoke(z)
	}

	return funcOf(f).call(o.v).Interface()
}

func (o *traitOption) GetOrElse(z interface{}) interface{} {
	if o.empty {
		return checkAndInvoke(z)
	}

	return o.Get()
}

func (o *traitOption) OrElse(f OptionOrElse) Option {
	if o.empty {
		return f()
	}

	return o
}

func (o *traitOption) Forall(f interface{}) bool {
	if o.empty {
		return true
	}

	return funcOf(f).call(o.v).Bool()
}

func (o *traitOption) Foreach(f interface{}) {
	if !o.empty {
		funcOf(f).call(o.v)
	}
}

func (o *traitOption) Empty() bool {
	return o.empty
}

// ------------------------------------------------

func optionCBF(x ...interface{}) Option {
	len := len(x)

	switch len {
	case 0:
		return &traitOption{
			empty: false,
			v:     unitValue,
		}
	case 1:
		switch v := x[0].(type) {
		case Option:
			return v
		case *_nothing:
			return None
		case reflect.Value:
			return optionCBF(v.Interface())
		default:
			if v == nil {
				return &traitOption{
					empty: false,
					v:     nullValue,
				}
			}

			return &traitOption{
				empty: false,
				v:     reflect.ValueOf(x[0]),
			}
		}

	default:
		return &traitOption{
			empty: false,
			v:     reflect.ValueOf(TupleOf(x)),
		}
	}
}

// ------------------------------------------------

// None represents scala-like None.
var None Option = &traitOption{empty: true, v: nothingValue}

// SomeOf returns Option, maybe None.
func SomeOf(x ...interface{}) Option {
	return optionCBF(x...)
}

// OptionOf returns a Option.
func OptionOf(x ...interface{}) Option {
	return optionCBF(x...)
}
