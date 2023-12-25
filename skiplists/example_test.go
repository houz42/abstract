package skiplists_test

import (
	"cmp"
	"fmt"
	"strings"

	"github.com/houz42/abstract/skiplists"
)

func ExampleSkipList() {
	list := skiplists.New[int]()

	list.Insert(3)
	list.Insert(2)
	list.Insert(4)
	list.Insert(1)

	list.Delete(2)

	for i := 1; i <= 5; i++ {
		fmt.Println(list.Search(i))
	}

	// Output:
	// 1 true
	// 0 false
	// 3 true
	// 4 true
	// 0 false
}

func ExampleNewFunc() {
	list := skiplists.NewFunc[string](func(a, b string) int {
		return cmp.Compare(strings.ToLower(a), strings.ToLower(b))
	})

	list.Insert("Hello")
	list.Insert("gopher")
	list.Insert("Go")
	list.Insert("is")
	list.Insert("fun")

	for i := 0; i < 5; i++ {
		fmt.Println(list.At(i))
	}

	// Output:
	// fun
	// Go
	// gopher
	// Hello
	// is
}

func ExampleSkipList_UpdateAt() {
	list := skiplists.NewFunc[string](func(a, b string) int {
		return cmp.Compare(strings.ToLower(a), strings.ToLower(b))
	})

	list.Insert("Hello")
	list.Insert("gopher")
	list.Insert("Go")
	list.Insert("is")
	list.Insert("fun")

	list.UpdateAt(2, "gophers")

	for i := 0; i < 5; i++ {
		fmt.Println(list.At(i))
	}

	// Output:
	// fun
	// Go
	// gophers
	// Hello
	// is
}

func ExampleSkipList_At() {
	list := skiplists.New[int]()

	list.Insert(3)
	list.Insert(5)
	list.Insert(2)
	list.Insert(4)
	list.Insert(1)
	fmt.Println(list.At(2))

	list.Delete(2)
	fmt.Println(list.At(2))

	list.DeleteAt(2)
	fmt.Println(list.At(2))

	// Output:
	// 3
	// 4
	// 5
}
