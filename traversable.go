package monadgo

// Traversable represents a scala-like Traversable trait.
type Traversable interface {
	sequence

	// Size returns the size.
	Size() int

	// Map applies function f to all elements in seq s.
	// f: func(T) X
	// returns a Traversable with element type X.
	Map(f interface{}) Traversable

	// FlatMap applies f to all elements and builds a new Travesable from result.
	// f: func(T) Traversable[X], X can be Go slice, or map.
	// returns a new Traversable with element type X.
	FlatMap(f interface{}) Traversable

	// Forall tests whether a predicate holds for all elements.
	// f: func(T) bool
	Forall(f interface{}) bool

	// Foreach applies f to all element.
	// f: func(T)
	Foreach(f interface{})

	// Fold folds the elements using specified associative binary operator.
	// z: func() T or value of type T.
	// f: func(T, T) T
	// returns value with type T
	Fold(z, f interface{}) interface{}

	// Reduce reduces the elements of this using the specified associative binary operator.
	// f: func(T, T)
	// returns value with type T.
	Reduce(f interface{}) interface{}

	// GroupBy returns Map with K -> Go slice. Key is the result of f. Collect elements into a slice with same resulting key value.
	// f: func(T) K
	// returns Map(K -> Go slice)
	GroupBy(f interface{}) Map

	// Exists tests whether a predicate holds for at least one element of this sequence.
	// f: func(T) bool
	Exists(f interface{}) bool

	// Find returns the first element satisfying f,
	// otherwise return None.
	Find(f interface{}) Option

	// Filter retuns all elements satisfying f.
	// f: func(T) bool
	Filter(f interface{}) Traversable

	// MkString displays all elements in a string using start, end, and separator sep.
	MkString(start, sep, end string) string

	// Split splits this into a unsatisfying and satisfying pair according to f.
	// f: func(T) bool
	Split(f interface{}) Tuple2

	// Collect returns elements satisfying pf.
	// pf is a partial function consisting of Condition func(T) bool and Action func(T) X.
	// returns a new Traversable[X]
	Collect(pf PartialFunc) Traversable
}
