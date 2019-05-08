package monadgo

import "fmt"

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
