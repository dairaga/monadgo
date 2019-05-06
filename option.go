package monadgo

import "fmt"

// OptionOrElse is function type for Option.OrElse.
type OptionOrElse func() Option

// Option represents scala-like Option.
type Option interface {
	Any

	Defined() bool

	Map(f interface{}) Option
	FlatMap(f interface{}) Option

	Foreach(f interface{})
	Forall(f interface{}) bool

	Fold(z, f interface{}) interface{}

	OrElse(OptionOrElse) Option
	GetOrElse(f interface{}) interface{}
}

// OptionOf returns a Option.
func OptionOf(x interface{}) Option {
	c, ok := x.(container)
	if !ok {
		c = containerOf(x)
	}

	return &traitOption{c}
}

// ----------------------------------------------------------------------------

type traitOption struct {
	container
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
		return &traitOption{o._map(f)}
	}

	return None
}

func (o *traitOption) FlatMap(f interface{}) Option {
	if !o.empty {
		return o._flatMap(f).(Option)
	}

	return None
}

func (o *traitOption) Fold(z, f interface{}) interface{} {
	if o.empty {
		return checkAndInvoke(z)
	}

	return o.invoke(f).Interface()
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

// ------------------------------------------------

// None represents scala-like None.
var None Option = &traitOption{nothingContainer}

// SomeOf returns Option, maybe None.
func SomeOf(x interface{}) Option {
	return OptionOf(x)
}
