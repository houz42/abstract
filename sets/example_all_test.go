//go:build goexperiment.rangefunc

package sets_test

import (
	"fmt"

	"github.com/houz42/abstract/sets"
)

func ExampleSet_All() {
	set := sets.New(1, 2, 3, 4, 5)
	for v := range set.All() {
		fmt.Println(v)
	}

	// Unordered Output:
	// 1
	// 2
	// 3
	// 4
	// 5
}
