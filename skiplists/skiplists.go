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

type SkipList[V any] struct {
	head  *node[V]
	level int
	size  int

	// cmp returns
	//
	//	-1 if x is less than y,
	//	 0 if x equals y,
	//	+1 if x is greater than y.
	//
	// For floating-point types, a NaN is considered less than any non-NaN,
	// a NaN is considered equal to a NaN, and -0.0 is equal to 0.0.
	cmp func(a, b V) int
}

func New[V cmp.Ordered]() *SkipList[V] {
	return NewFunc[V](cmp.Compare)
}

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

func (sl *SkipList[V]) Len() int { return sl.size }

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

func (sl *SkipList[V]) At(at int) V {
	if at < 0 || at >= sl.size {
		panic(fmt.Errorf("runtime error: index out of range [%d] with skip list length %d", at, sl.size))
	}

	node := sl.head
	pos := 0
	for level := sl.level - 1; level >= 0; level-- {
		for node != nil && pos+node.width[level] <= at {
			pos += node.width[level]
			node = node.next[level]
		}
	}

	return node.val
}

func (sl *SkipList[V]) DeleteAt(at int) {
	if at < 0 || at >= sl.size {
		panic(fmt.Errorf("runtime error: index out of range [%d] with skip list length %d", at, sl.size))
	}

	updates := make([]*node[V], sl.level)

	node := sl.head
	pos := 0
	for level := sl.level - 1; level >= 0; level-- {
		for node != nil && pos+node.width[level] < at {
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
