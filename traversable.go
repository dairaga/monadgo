package monadgo

import (
	"fmt"
	"reflect"
)

// Traversable represents a scala-like Traversable trait.
type Traversable interface {
	Any

	Map(f interface{}) Traversable

	FlatMap(f interface{}) Traversable

	// Forall tests whether a predicate holds for all elements.
	Forall(f interface{}) bool

	// Foreach applies f to all element.
	Foreach(f interface{})

	// Fold folds the elements using specified associative binary operator.
	Fold(z, f interface{}) interface{}

	// Reduce reduces the elements of this iterable collection using the specified associative binary operator.
	//Reduce(f interface{}) interface{}

	//Scan(z, f interface{})

	//GroupBy(f interface{})

	//Take(n int)

	//TakeWhile(f interface{})

	// ToSeq converts to slice.
	ToSeq() interface{}
}

// TraversableOf returns a Traversable.
func TraversableOf(xval reflect.Value) Traversable {
	switch xval.Kind() {
	case reflect.Slice:
		return newSlice(xval)
	case reflect.Map:
		return newMap(xval)
	default:
		panic(fmt.Sprintf("%v can not cast to Traversable", xval.Type()))
	}
}
