//go:build goexperiment.rangefunc

package skiplists_test

import (
	"fmt"

	"github.com/houz42/abstract/skiplists"
)

func ExampleSkipList_All() {
	list := skiplists.New[string]()

	list.Insert("Hello")
	list.Insert("gopher")
	list.Insert("Go")
	list.Insert("is")
	list.Insert("fun")

	for i, v := range list.All() {
		fmt.Println(i, v)
	}

	// Output:
	// 0 Go
	// 1 Hello
	// 2 fun
	// 3 gopher
	// 4 is
}
