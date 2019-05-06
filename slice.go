package monadgo

// Slice wraps Go slice and implement scala monadic functins.
type Slice interface {
	Traversable

	Len() int
	Cap() int
}

type slice = seq

// SliceOf returns a Slice.
func SliceOf(x interface{}) Slice {
	return seqOf(x)
}
