//go:build goexperiment.rangefunc

package skiplists

import "iter"

// All returns an iterator that yields all the ordered elements in the SkipList.
func (sl *SkipList[V]) All() iter.Seq2[int, V] {
	return func(yield func(int, V) bool) {
		node := sl.head.next[0]
		i := 0
		for node != nil {
			if !yield(i, node.val) {
				return
			}
			i++
			node = node.next[0]
		}
	}
}
