package monadgo

import (
	"fmt"
	"reflect"
)

// Either represents scala-like Either[A,B].
type Either interface {
	Any

	// IsLeft returns true if this is a Left, false otherwise.
	IsLeft() bool

	// IsRight returns true if this is a Right, false otherwise.
	IsRight() bool

	// Left projects this Either as a Left.
	Left() LeftProjection

	// FilterOrElse returns Right with the existing value of Right if this is a Right and the given predicate p holds for the right value,
	// or Left(zero) if this is a Right and the given predicate p does not hold for the right value,
	// or Left with the existing value of Left if this is a Left.
	// p: func(B) bool
	// zero: func() T
	// returns Either[T, B]
	FilterOrElse(p interface{}, zero interface{}) Either

	// Exists returns false if Left
	// or returns the result of the application of the given predicate to the Right value.
	// f: func(B) bool
	Exists(f interface{}) bool

	// Returns true if Left
	// or returns the result of the application of the given predicate to the Right value.
	// f: func(B) bool
	Forall(f interface{}) bool

	// Foreach executes the given side-effecting function f if this is a Right.
	// f: func(B)
	Foreach(f interface{})

	// Fold applies z if this is a Left or f if this is a Right.
	// z: func(A) X
	// f: func(B) X
	// returns value with type X.
	Fold(z, f interface{}) interface{}

	// Map is applying function f if this is a Right.
	// f: func(B) C
	// return Either[A, C]
	Map(f interface{}) Either

	// FlatMap binds the function f across Right.
	// f: func(B) Either
	// returns new Either.
	FlatMap(f interface{}) Either

	// GetOrElse returns the value from this Right,
	// or z if this is a Left.
	GetOrElse(z interface{}) interface{}

	// ToOption returns a Some containing the Right value if it exists,
	// or a None if this is a Left.
	ToOption() Option
}

type traitEither struct {
	right bool
	v     reflect.Value
}

func (e *traitEither) Get() interface{} {
	return e.v.Interface()
}

func (e *traitEither) rv() reflect.Value {
	return e.v
}

func (e *traitEither) String() string {
	if e.right {
		return fmt.Sprintf("Right(%v)", e.Get())
	}
	return fmt.Sprintf("Left(%v)", e.Get())
}

func (e *traitEither) IsLeft() bool {
	return !e.right
}

func (e *traitEither) IsRight() bool {
	return e.right
}

func (e *traitEither) Left() LeftProjection {
	return &leftProjection{
		e: e,
	}
}

func (e *traitEither) Forall(f interface{}) bool {
	if e.IsLeft() {
		return true
	}

	return funcOf(f).call(e.v).Bool()
}

func (e *traitEither) FilterOrElse(p, z interface{}) Either {
	if e.IsLeft() {
		return e
	}

	if funcOf(p).call(e.v).Bool() {
		return e
	}

	return LeftOf(checkAndInvoke(z))
}

func (e *traitEither) Exists(f interface{}) bool {
	if e.IsLeft() {
		return false
	}
	return funcOf(f).call(e.v).Bool()
}

func (e *traitEither) Foreach(f interface{}) {
	if e.IsRight() {
		funcOf(f).call(e.v)
	}
}

func (e *traitEither) Fold(z, f interface{}) interface{} {
	if e.IsLeft() {
		return funcOf(z).call(e.v).Interface()
	}

	return funcOf(f).call(e.v).Interface()
}

func (e *traitEither) Map(f interface{}) Either {
	if e.right {
		return &traitEither{
			right: true,
			v:     funcOf(f).call(e.v),
		}
	}

	return e
}

func (e *traitEither) FlatMap(f interface{}) Either {
	if e.right {
		return funcOf(f).call(e.v).Interface().(Either)
	}
	return e
}

func (e *traitEither) GetOrElse(z interface{}) interface{} {
	if e.IsLeft() {
		return checkAndInvoke(z)
	}
	return e.Get()
}

func (e *traitEither) ToOption() Option {
	if e.IsLeft() {
		return None
	}

	return OptionOf(e.Get())
}

// ----------------------------------------------------------------------------

// LeftOf returns Left of x.
func LeftOf(x ...interface{}) Either {
	return eitherCBF(false, x...)
}

// RightOf returns Right of x.
func RightOf(x ...interface{}) Either {
	return eitherCBF(true, x...)
}

func eitherCBF(right bool, x ...interface{}) Either {
	len := len(x)

	if len == 0 {
		return &traitEither{
			right: right,
			v:     unitValue,
		}
	}

	switch len {
	case 1:
		switch v := x[0].(type) {
		case Either:
			return v
		case reflect.Value:
			return eitherCBF(right, v.Interface())
		default:
			if v == nil {
				return &traitEither{
					right: right,
					v:     nullValue,
				}
			}
			return &traitEither{
				right: right,
				v:     reflect.ValueOf(x[0]),
			}
		}

	default:
		return &traitEither{
			right: right,
			v:     reflect.ValueOf(TupleOf(x)),
		}
	}
}
