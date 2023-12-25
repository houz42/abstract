// Package skiplists implements a [skip list],
// which can be used as an alternative to a balanced tree.
// Skip lists have the same expected time bounds as balanced trees but are simpler, faster, and use less space.
// They also provide additional fast random access methods.
//
// The implementation is based on [A Skip List Cookbook].
//
// [A Skip List Cookbook]: https://api.drum.lib.umd.edu/server/api/core/bitstreams/17176ef8-8330-4a6c-8b75-4cd18c570bec/content
// [skip list]: https://en.wikipedia.org/wiki/Skip_list
package skiplists

import (
	"cmp"
	"fmt"
	"math/bits"
	"math/rand"
)

type node[V any] struct {
	val   V
	width []int // for fast random access
	next  []*node[V]
}

// SkipList is a probabilistic data structure logically represents an ordered sequence of elements
// that allows $O(\log n)$ average complexity for search and insertion,
// as well as random accesses.
//
// A SkipList is not safe for concurrent use by multiple goroutines.
type SkipList[V any] struct {
	head  *node[V]
	level int
	size  int
	cmp   func(a, b V) int
}

// New returns a [SkipList] of any ordered elements
func New[V cmp.Ordered]() *SkipList[V] {
	return NewFunc[V](cmp.Compare)
}

// NewFunc returns a SkipList of any type when a custom cmp function is provided.
// The `cmp` function should return:
//   - -1 if a is less than b,
//   - 0 if a equals b,
//   - +1 if a is greater than b.
//
// NewFunc is intended to be used for types that cannot be ordered by the cmp.Ordered interface
// or for ordered types with a custom comparison operation (e.g., comparing floats within an approximate epsilon).
func NewFunc[V any](cmp func(a, b V) int) *SkipList[V] {
	return &SkipList[V]{
		head: &node[V]{
			width: []int{0},
			next:  []*node[V]{nil},
		},
		level: 1,
		cmp:   cmp,
	}
}

// Reverse returns a new SkipList which sort the elements in reversed order.
func (sl *SkipList[V]) Reverse() *SkipList[V] {
	return &SkipList[V]{
		head:  sl.head,
		level: sl.level,
		size:  sl.size,
		cmp:   func(a, b V) int { return sl.cmp(b, a) },
	}
}

// Len returns number of elements in the SkipList
func (sl *SkipList[V]) Len() int { return sl.size }

// Search returns an element and true if it is in the SkipList,
// otherwise zero value of type V and false.
// The returned value is useful when V is a custom type and the provided `cmp` method may
// return 0 (means equal) for different values, e.g. `cmp` only compares one field of a struct.
func (sl *SkipList[V]) Search(val V) (V, bool) {
	node := sl.head

	level := min(sl.level-1, sl.bestEntryLevel())
	for ; level >= 0; level-- {
		for node.next[level] != nil && sl.cmp(node.next[level].val, val) < 0 {
			node = node.next[level]
		}
	}
	node = node.next[0]

	if node == nil || sl.cmp(node.val, val) != 0 {
		var v V
		return v, false
	}
	return node.val, true
}

// Insert inserts an element into the SkipList.
// If the element is already in, the element will be overwritten with the input value.
func (sl *SkipList[V]) Insert(val V) {
	// nodes in each level just before the target
	updates := make([]*node[V], sl.level)

	// distances between each updates, used to fix width values
	jumps := make([]int, sl.level)

	nd := sl.head
	for level := sl.level - 1; level >= 0; level-- {
		jump := 0
		for nd.next[level] != nil && sl.cmp(nd.next[level].val, val) < 0 {
			jump += nd.width[level]
			nd = nd.next[level]
		}
		updates[level] = nd
		jumps[level] = jump
	}
	nd = nd.next[0]

	if nd != nil && sl.cmp(nd.val, val) == 0 {
		nd.val = val
		return
	}

	newLevel := sl.randomLevel()
	if newLevel > sl.level {
		updates = append(updates, make([]*node[V], newLevel-sl.level)...)
		jumps = append(jumps, make([]int, newLevel-sl.level)...)
		sl.head.next = append(sl.head.next, make([]*node[V], newLevel-sl.level)...)
		sl.head.width = append(sl.head.width, make([]int, newLevel-sl.level)...)

		for level := sl.level; level < newLevel; level++ {
			updates[level] = sl.head
			sl.head.width[level] = sl.size
		}

		sl.level = newLevel
	}

	node := &node[V]{
		val:   val,
		width: make([]int, newLevel),
		next:  make([]*node[V], newLevel),
	}

	for level := 0; level < newLevel; level++ {
		node.next[level] = updates[level].next[level]
		updates[level].next[level] = node

		if level == 0 {
			node.width[0] = 1
		} else {
			left := updates[level-1].width[level-1] + jumps[level-1]
			node.width[level] = updates[level].width[level] - left + 1
			updates[level].width[level] = left
		}
	}

	for level := newLevel; level < sl.level; level++ {
		updates[level].width[level]++
	}

	sl.size++
}

// Delete deletes an element from the SkipList.
// If the value is not found, nothing happens.
func (sl *SkipList[V]) Delete(val V) {
	updates := make([]*node[V], sl.level)

	node := sl.head
	for level := sl.level - 1; level >= 0; level-- {
		for node.next[level] != nil && sl.cmp(node.next[level].val, val) < 0 {
			node = node.next[level]
		}
		updates[level] = node
	}
	node = node.next[0]

	if node == nil || sl.cmp(node.val, val) != 0 {
		return
	}

	for level := 0; level < sl.level; level++ {
		if updates[level].next[level] == node {
			updates[level].width[level] += node.width[level] - 1
			updates[level].next[level] = node.next[level]
		} else {
			updates[level].width[level]--
		}
	}

	for sl.level > 1 && sl.head.next[sl.level-1] == nil {
		sl.level--
	}
	sl.head.next = sl.head.next[:sl.level]
	sl.head.width = sl.head.width[:sl.level]

	sl.size--
}

// At returns the i-th element in the SkipList.
// It panics if i is not valid, just like accessing slice element with an out-of-range index.
func (sl *SkipList[V]) At(i int) V {
	if i < 0 || i >= sl.size {
		panic(fmt.Errorf("runtime error: index out of range [%d] with skip list length %d", i, sl.size))
	}

	node := sl.head
	pos := 0
	for level := sl.level - 1; level >= 0; level-- {
		for node != nil && pos+node.width[level] <= i {
			pos += node.width[level]
			node = node.next[level]
		}
	}

	return node.val
}

// DeleteAt deletes i-th element in the SkipList.
// It panics if i is not valid, just like accessing slice element with an out-of-range index.
func (sl *SkipList[V]) DeleteAt(i int) {
	if i < 0 || i >= sl.size {
		panic(fmt.Errorf("runtime error: index out of range [%d] with skip list length %d", i, sl.size))
	}

	updates := make([]*node[V], sl.level)

	node := sl.head
	pos := 0
	for level := sl.level - 1; level >= 0; level-- {
		for node != nil && pos+node.width[level] < i {
			pos += node.width[level]
			node = node.next[level]
		}
		updates[level] = node
	}

	node = node.next[0]

	for level := 0; level < sl.level; level++ {
		if updates[level].next[level] == node {
			updates[level].width[level] += node.width[level] - 1
			updates[level].next[level] = node.next[level]
		} else {
			updates[level].width[level]--
		}
	}

	for sl.level > 1 && sl.head.next[sl.level-1] == nil {
		sl.level--
	}
	sl.head.next = sl.head.next[:sl.level]
	sl.head.width = sl.head.width[:sl.level]

	sl.size--
}

func (sl *SkipList[V]) bestEntryLevel() int {
	return bits.Len64(uint64(sl.size + 1))
}

func (sl *SkipList[V]) randomLevel() int {
	maxLevel := bits.Len64(uint64(sl.size + 1))
	level := bits.TrailingZeros64(rand.Uint64())
	return int(min(level, maxLevel)) + 1
}
