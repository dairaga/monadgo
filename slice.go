package monadgo

// Slice wraps Go slice and implement scala monadic functins.
type Slice interface {
	Traversable

	Len() int
	Cap() int

	Head() interface{}

	HeadOption() Option

	Tail() Traversable

	Take(n int) Traversable

	TakeWhile(f interface{}) Traversable

	Drop(n int) Traversable

	IndexWhere(f interface{}, start int) int

	LastIndexWhere(f interface{}, end int) int

	Reverse() Traversable

	Scan(z, f interface{}) Traversable
}

type slice = seq

// SliceOf returns a Slice.
func SliceOf(x interface{}) Slice {
	return seqOf(x)
}
