//go:build goexperiment.rangefunc

package heaps

import "iter"

// Drain returns an iterator that pops elements from the heap in the order of the heap.
// It is intentionally named "Drain" to distinguish it from other types' "All" methods,
// as it pops out elements with each call to yield.
func (h *Heap[E]) Drain() iter.Seq2[int, E] {
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
