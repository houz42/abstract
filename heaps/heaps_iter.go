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
