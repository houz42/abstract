// Package sets defines a set type and various methods useful with set operations.
package sets

// Set contains zero or more unique elements.
// The zero value for Set is nil and cannot be used directly.
// You should always use [New] or “make(sets.Set[int])” to create an empty set.
//
// A set is not safe for concurrent use by multiple goroutines.
type Set[E comparable] map[E]struct{}

// New creates a set with optional initial elements.
func New[E comparable](values ...E) Set[E] {
	m := make(Set[E], len(values))
	for _, v := range values {
		m[v] = struct{}{}
	}

	return m
}

// Len returns number of elements in the set.
func (s Set[E]) Len() int { return len(s) }

// Contains reports if v is in the set.
func (s Set[E]) Contains(v E) bool {
	_, ok := s[v]
	return ok
}

// Set insert an element into a set.
// If same element is already in the set, Set does nothing.
func (s Set[E]) Set(v E) { s[v] = struct{}{} }

// Unset removes an element from a set.
// If the element is not in the set, Unset does nothing.
func (s Set[E]) Unset(v E) { delete(s, v) }

// Map returns a new set in which each element is a mapping of the original ones.
// The mapping is done by calling `fn` on each element in the original set.
//
// Sadly, we cannot yet define a map method for set type as:
//
//	func (s Set[E]) Map[F comparable](fn func(E) F) Set[F]
//
// see go 1.18 [release note] for more details
//
// [release note]: https://tip.golang.org/doc/go1.18
func Map[E, F comparable](s Set[E], fn func(E) F) Set[F] {
	t := make(Set[F], len(s))
	for v := range s {
		t.Set(fn(v))
	}
	return t
}

// Map returns a new set whose elements are one-to-one mapping of the original set.
// The new set can only contain elements of same type with the original one.
// If other element type is expected, you have to use package function [Map] instead.
func (s Set[E]) Map(fn func(E) E) Set[E] {
	return Map(s, fn)
}

// Filter returns a new set contains elements in s satisfies fn.
func (s Set[E]) Filter(fn func(E) bool) Set[E] {
	t := make(Set[E])
	for v := range s {
		if fn(v) {
			t.Set(v)
		}
	}
	return t
}

// Equal reports whether two slices are equal: the same length and all elements equal.
// If the lengths are different, Equal returns false.
// Otherwise, the elements are compared one by one, and the comparison stops at the first element in s but not in t.
func (s Set[E]) Equal(t Set[E]) bool {
	return s.Len() == t.Len() && s.Subset(t)
}

// Clone returns a new set contains exactly same elements in s.
func (s Set[E]) Clone() Set[E] {
	t := make(Set[E], len(s))
	for v := range s {
		t.Set(v)
	}
	return t
}

// Subset reports whether s is subset of t:
// all elements in s are also in t.
func (s Set[E]) Subset(t Set[E]) bool {
	for v := range s {
		if !t.Contains(v) {
			return false
		}
	}

	return true
}

// Superset reports whether s is superset of t:
// all elements in t are also in s.
func (s Set[E]) Superset(t Set[E]) bool {
	return t.Subset(s)
}

// Union returns a new set contains elements either in s or in t.
func (s Set[E]) Union(t Set[E]) Set[E] {
	u := make(Set[E], s.Len()+t.Len())

	for v := range s {
		u.Set(v)
	}
	for v := range t {
		u.Set(v)
	}

	return u
}

// Intersection returns a new set contains elements both in s and in t.
func (s Set[E]) Intersection(t Set[E]) Set[E] {
	i := make(Set[E])

	for v := range s {
		if t.Contains(v) {
			i.Set(v)
		}
	}

	return i
}
