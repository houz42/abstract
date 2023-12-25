package skiplists_test

import (
	"fmt"

	"github.com/houz42/abstract/skiplists"
)

func ExampleSkipList() {
	list := skiplists.New[int, string]()

	list.Insert(3, "three")
	list.Insert(2, "two")
	list.Insert(4, "four")
	list.Insert(1, "one")

	list.Delete(2)

	for i := 1; i <= 5; i++ {
		fmt.Println(list.Search(i))
	}

	// Output:
	// one true
	//  false
	// three true
	// four true
	//  false
}

func ExampleSkipList_At() {
	list := skiplists.New[int, string]()

	list.Insert(3, "three")
	list.Insert(5, "five")
	list.Insert(2, "two")
	list.Insert(4, "four")
	list.Insert(1, "one")

	fmt.Println(list.At(2))

	list.Delete(2)
	fmt.Println(list.At(2))

	list.DeleteAt(2)
	fmt.Println(list.At(2))

	// Output:
	// 3 three
	// 4 four
	// 5 five
}
