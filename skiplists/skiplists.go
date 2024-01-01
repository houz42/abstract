// Package skiplists implements a [skip list].
// A skip list can be used as an alternative to a balanced tree, and as a priority queue.
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
	opt   Options
}

type node[V any] struct {
	val   V
	width []int // for fast random access
	next  []*node[V]
}

// New returns a [SkipList] of any ordered elements
func New[V cmp.Ordered](options ...Option) *SkipList[V] {
	return NewFunc[V](cmp.Compare, options...)
}

// NewFunc returns a SkipList of any type when a custom cmp function is provided.
// The `cmp` function should return:
//   - -1 if a is less than b,
//   - 0 if a equals b,
//   - +1 if a is greater than b.
//
// NewFunc is intended to be used for types that cannot be ordered by the cmp.Ordered interface
// or for ordered types with a custom comparison operation (e.g., comparing floats within an approximate epsilon).
func NewFunc[V any](cmp func(a, b V) int, options ...Option) *SkipList[V] {
	sl := &SkipList[V]{
		head: &node[V]{
			width: []int{0},
			next:  []*node[V]{nil},
		},
		level: 1,
		cmp:   cmp,
		opt:   defaultOptions,
	}

	for _, op := range options {
		op(&sl.opt)
	}

	return sl
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

// Get returns an element and true if it is in the SkipList,
// otherwise zero value of type V and false.
// The returned value is useful when V is a custom type and the provided `cmp` method may
// return 0 (means equal) for different values, e.g., `cmp` only compares one field of a struct.
func (sl *SkipList[V]) Get(val V) (V, bool) {
	node := sl.head

	level := min(sl.level, maxLevel(sl.opt.LogP, sl.size)) - 1
	for ; level >= 0; level-- {
		for node.next[level] != nil && sl.cmp(node.next[level].val, val) < 0 {
			node = node.next[level]
		}
		if node.next[level] != nil && sl.cmp(node.next[level].val, val) == 0 {
			node = node.next[level]
			goto CHECK
		}
	}
	node = node.next[0]

CHECK:
	if node == nil || sl.cmp(node.val, val) != 0 {
		var v V
		return v, false
	}
	return node.val, true
}

// Set inserts an element into the SkipList.
// If the element is already in, the element will be overwritten with the input value.
func (sl *SkipList[V]) Set(val V) *SkipList[V] {
	// nodes in each level just before the target
	updates := make([]*node[V], sl.level)

	// indexes of each node when jumps to lower level
	jumps := make([]int, sl.level)
	pos := -1

	nd := sl.head
	for level := sl.level - 1; level >= 0; level-- {
		for nd.next[level] != nil && sl.cmp(nd.next[level].val, val) <= 0 {
			pos += nd.width[level]
			nd = nd.next[level]
		}
		updates[level] = nd
		jumps[level] = pos
	}

	if nd.next[0] != nil && sl.cmp(nd.next[0].val, val) == 0 {
		nd.val = val
		return sl
	}

	newLevel := sl.randomLevel()

	newNode := &node[V]{
		val:   val,
		width: make([]int, newLevel),
		next:  make([]*node[V], newLevel),
	}
	pos++ // index of new node

	// add new levels if needed
	if newLevel > sl.level {
		sl.head.next = append(sl.head.next[:sl.level], make([]*node[V], newLevel-sl.level)...)
		sl.head.width = append(sl.head.width[:sl.level], make([]int, newLevel-sl.level)...)

		for level := sl.level; level < newLevel; level++ {
			sl.head.next[level] = newNode
			sl.head.width[level] = pos + 1
			newNode.width[level] = sl.size - pos
		}

		sl.level = newLevel
	}

	for level := 0; level < min(newLevel, len(updates)); level++ {
		newNode.next[level] = updates[level].next[level]
		updates[level].next[level] = newNode

		left := pos - jumps[level]
		newNode.width[level] = updates[level].width[level] - left + 1
		updates[level].width[level] = left
	}

	for level := newLevel; level < len(updates); level++ {
		updates[level].width[level]++
	}

	sl.size++

	return sl
}

// Unset removes an element from the SkipList.
// If the element is not found, nothing happens.
func (sl *SkipList[V]) Unset(val V) *SkipList[V] {
	updates := make([]*node[V], sl.level)

	node := sl.head
	for level := sl.level - 1; level >= 0; level-- {
		for node.next[level] != nil && sl.cmp(node.next[level].val, val) < 0 {
			node = node.next[level]
		}
		updates[level] = node
	}

	// the node we are removing
	node = node.next[0]
	if node == nil || sl.cmp(node.val, val) != 0 {
		return sl
	}

	// remove node from each level
	for level := 0; level < sl.level; level++ {
		if updates[level].next[level] == node {
			updates[level].width[level] += node.width[level] - 1
			updates[level].next[level] = node.next[level]
		} else {
			updates[level].width[level]--
		}

		// remove higher levels contains nothing
		if level > 1 && updates[level] == sl.head && sl.head.next[level] == nil {
			sl.level = level - 1
			break
		}
	}

	sl.size--

	return sl
}

// At returns the i-th element in the SkipList.
// It panics if i is not valid, just like accessing slice element with an out-of-range index.
func (sl *SkipList[V]) At(i int) V {
	if i < 0 || i >= sl.size {
		panic(fmt.Errorf("runtime error: index out of range [%d] with skip list length %d", i, sl.size))
	}

	node := sl.head
	pos := -1
	for level := sl.level - 1; level >= 0; level-- {
		for node.next[level] != nil && pos+node.width[level] <= i {
			pos += node.width[level]
			node = node.next[level]
		}
		if pos == i {
			return node.val
		}
	}

	return node.val
}

// UpdateAt updates the value at the specified index in the SkipList.
// It panics with a runtime error if the index is out of range,
// or the new value will violate the SkipList's ordering property.
func (sl *SkipList[V]) UpdateAt(i int, val V) {
	if i < 0 || i >= sl.size {
		panic(fmt.Errorf("runtime error: index out of range [%d] with skip list length %d", i, sl.size))
	}

	node := sl.head
	pos := -1

	for level := sl.level - 1; level >= 0; level-- {
		for node != nil && pos+node.width[level] < i {
			pos += node.width[level]
			node = node.next[level]
		}
	}

	target := node.next[0]

	if sl.cmp(node.val, val) >= 0 {
		panic("updated value is less than or equal to the previous one")
	}
	if next := target.next[0]; next != nil && sl.cmp(val, next.val) >= 0 {
		panic("updated value is greater than or equal to the next one")
	}

	target.val = val
}

// RemoveAt removes the i-th element in the SkipList.
// It panics if i is not valid, just like accessing slice element with an out-of-range index.
func (sl *SkipList[V]) RemoveAt(i int) *SkipList[V] {
	if i < 0 || i >= sl.size {
		panic(fmt.Errorf("runtime error: index out of range [%d] with skip list length %d", i, sl.size))
	}

	// unlike in [Unset], we are sure a node will be removed,
	// so we can do the removing at the same time as we search for it.

	node := sl.head
	pos := -1
	for level := sl.level - 1; level >= 0; level-- {
		for node.next[level] != nil && pos+node.width[level] < i {
			pos += node.width[level]
			node = node.next[level]
		}

		// the node we are removing must be in the range
		node.width[level]--

		// found the node we are removing, merge its width with the following part
		if node.next[level] != nil && pos+node.width[level]+1 == i {
			node.width[level] += node.next[level].width[level]
			node.next[level] = node.next[level].next[level]

			if node == sl.head && node.next[level] == nil {
				sl.level--
			}
		}
	}

	sl.size--

	return sl
}

func maxLevel(logP, size int) int {
	return bits.Len64(uint64(size)+1) / logP
}

func (sl *SkipList[V]) randomLevel() int {
	level := sl.opt.maxLevel
	if sl.opt.SizeHint < sl.size {
		level = maxLevel(sl.opt.LogP, sl.size)
	}
	level = bits.TrailingZeros64(rand.Uint64()|(1<<(level*sl.opt.LogP))) / sl.opt.LogP
	return level + 1
}

// Options represents the internal configurations of a SkipList.
type Options struct {

	// LogP is the log base of the probability of each level (the P  value):
	//
	// 	$LogP = - \log_{2}{P}$
	//
	// Commonly used P values are:
	// 	- 0.5  (LogP = 1, the default)
	// 	- 0.25 (LogP = 2)
	//
	// see the [wiki page] for more details about the P value
	//
	// [wiki page]: https://en.wikipedia.org/wiki/Skip_list
	LogP int

	// SizeHint is the expected total size of the elements.
	// It is used to hint the max level of the SkipList.
	// If it is not set, max level will be calculated based on current size dynamically.
	SizeHint int
	maxLevel int
}

// Option changes a SkipList's Options
type Option func(*Options)

// SetLogP sets the log base of the P value.
func SetLogP(logP int) Option {
	return func(o *Options) {
		o.LogP = logP
		o.maxLevel = maxLevel(logP, o.SizeHint)
	}
}

// SetSizeHint sets the expected number of elements in the SkipList.
func SetSizeHint(hint int) Option {
	return func(o *Options) {
		o.SizeHint = hint
		o.maxLevel = maxLevel(o.LogP, hint)
	}
}

var defaultOptions = Options{
	LogP: 1, // P = 0.5
}
