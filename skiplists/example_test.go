package skiplists_test

import (
	"cmp"
	"fmt"
	"strings"

	"github.com/houz42/abstract/skiplists"
)

func Example() {
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

func Example_priorityQueue() {

	type process struct {
		pid      int
		niceness int
	}

	queue := skiplists.NewFunc[process](func(a, b process) int {
		return cmp.Compare(a.niceness, b.niceness)
	})

	queue.Insert(process{pid: 1, niceness: -20})
	queue.Insert(process{pid: 2, niceness: 0})
	queue.Insert(process{pid: 3, niceness: 10})
	queue.Insert(process{pid: 4, niceness: -1})

	for queue.Len() > 0 {
		p := queue.At(0)
		fmt.Printf("start process %d with niceness %d\n", p.pid, p.niceness)
		queue.DeleteAt(0)
	}

	// Output:
	// start process 1 with niceness -20
	// start process 4 with niceness -1
	// start process 2 with niceness 0
	// start process 3 with niceness 10
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
