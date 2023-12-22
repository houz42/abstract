// Package lists is a simple wrapper to container/list which
// in turn implements a doubly linked list, to make it generic and type-safe.
package lists

import "container/list"

// Element is an element of a linked list.
type Element[T any] struct {
	impl *list.Element
}

func wrapElement[T any](ele *list.Element) *Element[T] {
	// to be aligned with the container/list return values
	if ele == nil {
		return nil
	}
	return &Element[T]{impl: ele}
}

// Next returns the next list element or nil.
func (ele *Element[T]) Next() *Element[T] { return wrapElement[T](ele.impl.Next()) }

// Prev returns the previous list element or nil.
func (ele *Element[T]) Prev() *Element[T] { return wrapElement[T](ele.impl.Prev()) }

// Value returns the value stored with this element.
func (ele *Element[T]) Value() T {
	return ele.impl.Value.(T)
}

// List represents a doubly linked list.
// The zero value for List is an empty list ready to use.
type List[T any] struct {
	impl *list.List
}

// New returns an initialized list.
func New[T any]() *List[T] { return &List[T]{impl: list.New()} }

// Len returns the number of elements of list l.
// The complexity is O(1).
func (l *List[T]) Len() int { return l.impl.Len() }

// Front returns the first element of list l or nil if the list is empty.
func (l *List[T]) Front() *Element[T] { return wrapElement[T](l.impl.Front()) }

// Back returns the last element of list l or nil if the list is empty.
func (l *List[T]) Back() *Element[T] { return wrapElement[T](l.impl.Back()) }

// PushFront inserts a new element e with value v at the front of list l and returns e.
func (l *List[T]) PushFront(v T) *Element[T] {
	return wrapElement[T](l.impl.PushFront(v))
}

// PushBack inserts a new element e with value v at the back of list l and returns e.
func (l *List[T]) PushBack(v T) *Element[T] {
	return wrapElement[T](l.impl.PushBack(v))
}

// InsertBefore inserts a new element e with value v immediately before mark and returns e.
// If mark is not an element of l, the list is not modified.
// The mark must not be nil.
func (l *List[T]) InsertBefore(v T, mark *Element[T]) *Element[T] {
	return wrapElement[T](l.impl.InsertBefore(v, mark.impl))
}

// InsertAfter inserts a new element e with value v immediately after mark and returns e.
// If mark is not an element of l, the list is not modified.
// The mark must not be nil.
func (l *List[T]) InsertAfter(v T, mark *Element[T]) *Element[T] {
	return wrapElement[T](l.impl.InsertAfter(v, mark.impl))
}

// MoveToFront moves element e to the front of list l.
// If e is not an element of l, the list is not modified.
// The element must not be nil.
func (l *List[T]) MoveToFront(e *Element[T]) { l.impl.MoveToFront(e.impl) }

// MoveToBack moves element e to the back of list l.
// If e is not an element of l, the list is not modified.
// The element must not be nil.
func (l *List[T]) MoveToBack(e *Element[T]) { l.impl.MoveToBack(e.impl) }

// MoveBefore moves element e to its new position before mark.
// If e or mark is not an element of l, or e == mark, the list is not modified.
// The element and mark must not be nil.
func (l *List[T]) MoveBefore(e, mark *Element[T]) { l.impl.MoveBefore(e.impl, mark.impl) }

// MoveAfter moves element e to its new position after mark.
// If e or mark is not an element of l, or e == mark, the list is not modified.
// The element and mark must not be nil.
func (l *List[T]) MoveAfter(e, mark *Element[T]) { l.impl.MoveAfter(e.impl, mark.impl) }

// PushBackList inserts a copy of another list at the back of list l.
// The lists l and other may be the same. They must not be nil.
func (l *List[T]) PushBackList(other *List[T]) { l.impl.PushBackList(other.impl) }

// PushFrontList inserts a copy of another list at the front of list l.
// The lists l and other may be the same. They must not be nil.
func (l *List[T]) PushFrontList(other *List[T]) { l.impl.PushFrontList(other.impl) }

// Remove removes e from l if e is an element of list l.
// It returns the element value e.Value.
// The element must not be nil.
func (l *List[T]) Remove(e *Element[T]) T {
	return l.impl.Remove(e.impl).(T)
}
