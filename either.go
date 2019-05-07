package monadgo

import (
	"fmt"
)

// Either represents scala-like Either.
type Either interface {
	Any
	IsLeft() bool
	IsRight() bool

	Left() LeftProjection
	//Right() Either

	FilterOrElse(p interface{}, z interface{}) Either

	Exists(f interface{}) bool
	Forall(f interface{}) bool
	Foreach(f interface{})
	Fold(z, f interface{}) interface{}

	Map(f interface{}) Either
	FlatMap(f interface{}) Either

	GetOrElse(z interface{}) interface{}

	ToOption() Option
}

type traitEither struct {
	right bool
	container
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

	return e.Forall(f)
}

func (e *traitEither) FilterOrElse(p, z interface{}) Either {
	if e.IsLeft() {
		return e
	}

	if e.invoke(p).Bool() {
		return e
	}

	return LeftOf(checkAndInvoke(z))
}

func (e *traitEither) Exists(f interface{}) bool {
	if e.IsLeft() {
		return false
	}

	return e.invoke(f).Bool()
}

func (e *traitEither) Foreach(f interface{}) {
	if e.IsRight() {
		e.Foreach(f)
	}
}

func (e *traitEither) Fold(z, f interface{}) interface{} {
	if e.IsLeft() {
		return e.invoke(z).Interface()
	}

	return e.invoke(f).Interface()
}

func (e *traitEither) Map(f interface{}) Either {
	if e.right {
		return eitherFromContainer(true, e._map(f))
	}

	return e
}

func (e *traitEither) FlatMap(f interface{}) Either {
	if e.right {
		return e._flatMap(f).(Either)
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

func eitherFromContainer(right bool, c container) Either {
	return &traitEither{
		right:     right,
		container: c,
	}
}

func eitherOf(right bool, x ...interface{}) Either {
	switch len(x) {
	case 0:
		return &traitEither{
			right:     right,
			container: nothingContainer,
		}
	case 1:
		return &traitEither{
			right:     right,
			container: containerOf(x),
		}
	default:
		return &traitEither{
			right:     right,
			container: containerOf(TupleOf(x)),
		}
	}
}

// LeftOf returns Left of x.
func LeftOf(x ...interface{}) Either {
	return eitherOf(false, x...)
}

// RightOf returns Right of x.
func RightOf(x ...interface{}) Either {
	return eitherOf(true, x...)
}
