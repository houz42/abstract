// Package heaps provides a heap implementation for any type.
package heaps

import (
	"cmp"
	"container/heap"
)

// Heap is a tree with the property that each node is the minimum-valued (or maximum if reversed)
// node in its subtree.
// The minimum (maximum) element in the tree is the root, at index 0.
//
// A newly created Heap is a min-heap.
// If a max-heap is wanted, call [Reverse] on it.
//
// A heap is a common way to implement a priority queue. To build a priority queue,
// pass in a less method which orders the elements by their priorities,
// so [Push] adds items while [Pop] removes the highest-priority item from the queue.
// See the example for more details.
//
// A Heap is not safe for concurrent use by multiple goroutines.
type Heap[E any] struct {
	impl *heapImpl[E]
}

// New creates a new min-heap for ordered element types.
// The initial values are optional.
func New[E cmp.Ordered](values ...E) *Heap[E] {
	return NewFunc(func(x, y E) bool { return x < y }, values...)
}

// NewFunc creates a new min-heap for any type.
// The initial values are optional.
func NewFunc[E any](less func(x, y E) bool, values ...E) *Heap[E] {
	impl := &heapImpl[E]{
		values: values,
		less:   less,
	}
	heap.Init(impl)
	return &Heap[E]{impl: impl}
}

// Reverse returns a new Heap in which the elements will be pop out in reserved sequence to the original one.
// That is, if h is a min-heap, a max-heap will be returned, or vice versa.
func (h *Heap[E]) Reverse() *Heap[E] {
	r := &Heap[E]{
		impl: &heapImpl[E]{
			values: make([]E, 0, len(h.impl.values)),
			less:   func(x, y E) bool { return h.impl.less(y, x) },
		},
	}
	copy(r.impl.values, h.impl.values)
	heap.Init(r.impl)

	return r
}

// Len returns number of elements in the heap.
func (h *Heap[E]) Len() int { return len(h.impl.values) }

// Push pushes the element x onto the heap.
// The complexity is O(log n) where n = h.Len().
func (h *Heap[E]) Push(x E) {
	heap.Push(h.impl, x)
}

// Pop removes and returns the first element from the heap.
// The complexity is O(log n) where n = h.Len().
// Pop is equivalent to [Remove](h).
func (h *Heap[E]) Pop() E {
	v := heap.Pop(h.impl)
	return v.(E)
}

// Top returns the first element from the heap.
// The complexity is O(1).
func (h *Heap[E]) Top() E {
	return h.impl.values[0]
}

// Remove removes and returns the element at index i from the heap.
// The complexity is O(log n) where n = h.Len().
func (h *Heap[E]) Remove(i int) E {
	v := heap.Remove(h.impl, i)
	return v.(E)
}

type heapImpl[E any] struct {
	values []E
	less   func(x, y E) bool
}

func (h *heapImpl[E]) Len() int           { return len(h.values) }
func (h *heapImpl[E]) Less(i, j int) bool { return h.less(h.values[i], h.values[j]) }
func (h *heapImpl[E]) Swap(i, j int)      { h.values[i], h.values[j] = h.values[j], h.values[i] }

func (h *heapImpl[E]) Push(x any) {
	h.values = append(h.values, x.(E))
}

func (h *heapImpl[E]) Pop() any {
	old := h.values
	n := len(old)
	x := old[n-1]
	h.values = old[0 : n-1]
	return x
}
