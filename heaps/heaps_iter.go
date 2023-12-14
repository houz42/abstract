//go:build goexperiment.rangefunc

package heaps

import "iter"

// PopAll returns an iterator that pops elements in the heap by the heap order.
func (h *Heap[E]) PopAll() iter.Seq2[int, E] {
	return func(yield func(int, E) bool) {
		i := 0
		for h.Len() > 0 {
			if !yield(i, h.Pop()) {
				return
			}
			i++
		}
	}
}

// All returns an iterator to access the elements in the heap.
// The elements will not be pop out, and the order is not following the heap order.
func (h *Heap[E]) All() iter.Seq2[int, E] {
	return func(yield func(int, E) bool) {
		for i, v := range h.impl.values {
			if !yield(i, v) {
				return
			}
			i++
		}
	}
}
