package monadgo

import (
	"fmt"
	"reflect"
)

// Traversable represents a scala-like Traversable trait.
type Traversable interface {
	Any

	// Forall tests whether a predicate holds for all elements.
	Forall(f interface{}) bool

	// Foreach applies f to all element.
	Foreach(f interface{})

	// Fold folds the elements using specified associative binary operator.
	Fold(z, f interface{}) interface{}

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
