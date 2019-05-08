package monadgo

import (
	"fmt"
	"reflect"
)

// LeftProjection represents scala-like LeftProjection[A,B].
type LeftProjection interface {
	Any

	// E return internal Either value.
	E() Either

	// Exists returns false if Right,
	// or returns the result of the application of the given function to the Left value.
	// func(A) bool
	Exists(f interface{}) bool

	// Filter returns None if this is a Right,
	// or if the given predicate p does not hold for the left value,
	// otherwise return a Left.
	// f: func(A) bool
	Filter(f interface{}) Option

	// FlatMap binds the given function f across Left.
	// f: func(A) Either
	// returns a new Either.
	FlatMap(f interface{}) Either

	// Forall returns true if Right,
	// or returns the result of the application of the given function to the Left value.
	// f: func(A) bool
	Forall(f interface{}) bool

	// Foreach executes the given side-effecting function f if this is a Left.
	// f: func(A)
	Foreach(f interface{})

	// GetOrElse returns the value from this Left,
	// or z if this is a Right.
	GetOrElse(z interface{}) interface{}

	// Map applies f through Left.
	// f: func(A) X
	// returns Either[X, B]
	Map(f interface{}) Either

	// ToOption returns a Some containing the Left value if it exists,
	// or a None if this is a Right.
	ToOption() Option
}

type leftProjection struct {
	e *traitEither
}

var _ LeftProjection = &leftProjection{}

func (p *leftProjection) Get() interface{} {
	if p.e.IsRight() {
		return nothing
	}
	return p.e.Get()
}

func (p *leftProjection) String() string {
	return fmt.Sprintf("LeftProjection(%v)", p.Get())
}

func (p *leftProjection) rv() reflect.Value {
	return p.e.rv()
}

func (p *leftProjection) E() Either {
	return p.e
}

func (p *leftProjection) Exists(f interface{}) bool {
	if p.e.IsRight() {
		return false
	}

	return p.e.invoke(f).Bool()
}

func (p *leftProjection) Filter(f interface{}) Option {
	if p.e.IsRight() {
		return None
	}

	if p.e.invoke(f).Bool() {
		return OptionOf(p.E())
	}

	return None
}

func (p *leftProjection) FlatMap(f interface{}) Either {
	if p.e.IsRight() {
		return p.E()
	}

	return p.e._flatMap(f).(Either)
}

func (p *leftProjection) Forall(f interface{}) bool {
	if p.e.IsRight() {
		return true
	}

	return p.e.invoke(f).Bool()
}

func (p *leftProjection) Foreach(f interface{}) {
	if p.e.IsLeft() {
		p.e.invoke(f)
	}
}

func (p *leftProjection) GetOrElse(z interface{}) interface{} {
	if p.e.IsRight() {
		return checkAndInvoke(z)
	}

	return p.Get()
}

func (p *leftProjection) Map(f interface{}) Either {
	if p.e.IsRight() {
		return p.E()
	}

	return eitherFromContainer(false, p.e._map(f))
}

func (p *leftProjection) ToOption() Option {
	if p.e.IsRight() {
		return None
	}

	return SomeOf(p.e.Get())
}
