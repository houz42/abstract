//go:build goexperiment.rangefunc

package lists

import "iter"

// All returns an iterator that yields elements in the list in order.
func (l *List[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		for e := l.Front(); e != nil; e = e.Next() {
			if !yield(e.Value()) {
				return
			}
		}
	}
}

// Backward returns an iterator that yields elements in the list in backward order.
func (l *List[T]) Backward() iter.Seq[T] {
	return func(yield func(T) bool) {
		for e := l.Back(); e != nil; e = e.Prev() {
			if !yield(e.Value()) {
				return
			}
		}
	}
}
