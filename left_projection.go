package monadgo

import (
	"fmt"
	"reflect"
)

// LeftProjection represents scala-like LeftProjection.
type LeftProjection interface {
	Any
	E() Either

	Exists(f interface{}) bool
	Filter(f interface{}) Option
	FlatMap(f interface{}) Either
	Forall(f interface{}) bool
	Foreach(f interface{})
	GetOrElse(z interface{}) interface{}
	Map(f interface{}) Either
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
