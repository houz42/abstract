//go:build goexperiment.rangefunc

package lists_test

import (
	"fmt"

	"github.com/houz42/abstract/lists"
)

func ExampleList_All() {
	// Create a new list and put some numbers in it.
	l := lists.New[int]()
	e4 := l.PushBack(4)
	e1 := l.PushFront(1)
	l.InsertBefore(3, e4)
	l.InsertAfter(2, e1)

	// Iterate through list and print its contents.
	for e := range l.All() {
		fmt.Println(e)
	}
	// Output:
	// 1
	// 2
	// 3
	// 4
}

func ExampleList_Backward() {
	// Create a new list and put some numbers in it.
	l := lists.New[int]()
	e4 := l.PushBack(4)
	e1 := l.PushFront(1)
	l.InsertBefore(3, e4)
	l.InsertAfter(2, e1)

	// Iterate through list and print its contents.
	for e := range l.Backward() {
		fmt.Println(e)
	}
	// Output:
	// 4
	// 3
	// 2
	// 1
}
