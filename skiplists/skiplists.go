package skiplists

import (
	"cmp"
	"fmt"
	"math/bits"
	"math/rand"
)

type node[K cmp.Ordered, V any] struct {
	key   K
	val   V
	width []int // for fast random access
	next  []*node[K, V]
}

type SkipList[K cmp.Ordered, V any] struct {
	head  *node[K, V]
	level int
	size  int

	// TODO: any key type with custom comparator
	// cmp   func(a, b K) int
}

func New[K cmp.Ordered, V any]() *SkipList[K, V] {
	return &SkipList[K, V]{
		head: &node[K, V]{
			width: []int{0},
			next:  []*node[K, V]{nil},
		},
		level: 1,
	}
}

func (sl *SkipList[K, V]) Len() int { return sl.size }

func (sl *SkipList[K, V]) Search(key K) (V, bool) {
	node := sl.head

	// TODO: jump to best entry level
	for level := sl.level - 1; level >= 0; level-- {
		for node.next[level] != nil && node.next[level].key < key {
			node = node.next[level]
		}
	}
	node = node.next[0]

	if node == nil || node.key != key {
		var v V
		return v, false
	}
	return node.val, true
}

func (sl *SkipList[K, V]) Insert(key K, val V) {
	// nodes in each level just before the target
	updates := make([]*node[K, V], sl.level)

	// distances between each updates, used to fix width values
	jumps := make([]int, sl.level)

	nd := sl.head
	for level := sl.level - 1; level >= 0; level-- {
		jump := 0
		for nd.next[level] != nil && nd.next[level].key < key {
			jump += nd.width[level]
			nd = nd.next[level]
		}
		updates[level] = nd
		jumps[level] = jump
	}
	nd = nd.next[0]

	if nd != nil && nd.key == key {
		nd.val = val
		return
	}

	newLevel := sl.randomLevel()
	if newLevel > sl.level {
		updates = append(updates, make([]*node[K, V], newLevel-sl.level)...)
		jumps = append(jumps, make([]int, newLevel-sl.level)...)
		sl.head.next = append(sl.head.next, make([]*node[K, V], newLevel-sl.level)...)
		sl.head.width = append(sl.head.width, make([]int, newLevel-sl.level)...)

		for level := sl.level; level < newLevel; level++ {
			// updates = append(updates, sl.head)
			// sl.head.width = append(sl.head.width, sl.size)
			updates[level] = sl.head
			sl.head.width[level] = sl.size
		}

		sl.level = newLevel
	}

	node := &node[K, V]{
		key:   key,
		val:   val,
		width: make([]int, newLevel),
		next:  make([]*node[K, V], newLevel),
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

func (sl *SkipList[K, V]) randomLevel() int {
	maxLevel := bits.Len64(uint64(sl.size + 1))
	level := bits.TrailingZeros64(rand.Uint64())
	return int(min(level, maxLevel)) + 1
}

func (sl *SkipList[K, V]) Delete(key K) {
	update := make([]*node[K, V], sl.level)

	nd := sl.head
	for level := sl.level - 1; level >= 0; level-- {
		for nd.next[level] != nil && nd.next[level].key < key {
			nd = nd.next[level]
		}
		update[level] = nd
	}
	nd = nd.next[0]

	if nd == nil || nd.key != key {
		return
	}

	for level := 0; level < sl.level; level++ {
		if update[level].next[level] == nd {
			update[level].width[level] += nd.width[level] - 1
			update[level].next[level] = nd.next[level]
		} else {
			update[level].width[level]--
		}
	}

	for sl.level > 1 && sl.head.next[sl.level-1] == nil {
		sl.level--
	}
	sl.head.next = sl.head.next[:sl.level]
	sl.head.width = sl.head.width[:sl.level]

	sl.size--
}

func (sl *SkipList[K, V]) At(at int) (K, V) {
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

		if pos == at {
			return node.key, node.val
		}
	}

	return node.key, node.val
}

func (sl *SkipList[K, V]) DeleteAt(at int) {
	if at < 0 || at >= sl.size {
		panic(fmt.Errorf("runtime error: index out of range [%d] with skip list length %d", at, sl.size))
	}

	updates := make([]*node[K, V], sl.level)

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
