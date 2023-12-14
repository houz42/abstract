//go:build goexperiment.rangefunc

package sets

import "iter"

// All returns all elements in the set as an iterator.
//
// All is useful if an iterator is wanted for some interface,
// but it provides no extra convenience by range over the returned iterator as you can always do:
//
//	for v := range set {
//		...
//	}
func (s Set[E]) All() iter.Seq[E] {
	return func(yield func(E) bool) {
		for v := range s {
			if !yield(v) {
				return
			}
		}
	}
}
