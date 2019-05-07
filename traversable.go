package monadgo

// Traversable represents a scala-like Traversable trait.
type Traversable interface {
	sequence

	Size() int

	Map(f interface{}) Traversable

	FlatMap(f interface{}) Traversable

	// Forall tests whether a predicate holds for all elements.
	Forall(f interface{}) bool

	// Foreach applies f to all element.
	Foreach(f interface{})

	// Fold folds the elements using specified associative binary operator.
	Fold(z, f interface{}) interface{}

	// Reduce reduces the elements of this iterable collection using the specified associative binary operator.
	Reduce(f interface{}) interface{}

	GroupBy(f interface{}) Map

	Exists(f interface{}) bool

	Find(f interface{}) Option

	Filter(f interface{}) Traversable

	MkString(start, sep, end string) string

	Span(f interface{}) Tuple2
}
