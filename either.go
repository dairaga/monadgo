package monadgo

import (
	"fmt"
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

	return e.container.Forall(f)
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
		e.container.Foreach(f)
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

func newEither(right bool, x ...interface{}) Either {
	switch len(x) {
	case 0:
		return eitherFromContainer(right, nothingContainer)
	case 1:
		return eitherFromContainer(right, containerOf(x[0]))
	default:
		return eitherFromContainer(right, containerOf(TupleOf(x)))
	}
}

// LeftOf returns Left of x.
func LeftOf(x ...interface{}) Either {
	return newEither(false, x...)
}

// RightOf returns Right of x.
func RightOf(x ...interface{}) Either {
	return newEither(true, x...)
}

// EitherOf returns a Either with one parameter.
// errOrFalse must be bool or error type.
// Return Left if errOrFalse is false or error existing,
// or Right of Null.
func EitherOf(errOrFalse interface{}) Either {
	if isErrorOrFalse(errOrFalse) {
		return LeftOf(errOrFalse)
	}
	return RightOf(errOrFalse)
}

// Either1Of returns a Either.
// errOrFalse must be bool or error type.
// Return Left if errOrFalse is false or error existing,
// or Right of x.
func Either1Of(x, errOrFalse interface{}) Either {
	if isErrorOrFalse(errOrFalse) {
		return LeftOf(errOrFalse)
	}

	return RightOf(x)
}

// Either2Of returns a Try.
// errOrFalse must be bool or error type.
// Return Left if errOrFalse is false or error existing,
// or Right of Tuple2(x1,x2).
func Either2Of(x1, x2, errOrFalse interface{}) Either {
	if isErrorOrFalse(errOrFalse) {
		return LeftOf(errOrFalse)
	}

	return RightOf(Tuple2Of(x1, x2))
}
