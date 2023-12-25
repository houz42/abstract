//go:build goexperiment.rangefunc

package heaps_test

import (
	"fmt"

	"github.com/houz42/abstract/heaps"
)

func ExampleHeap_PopAll() {
	h := heaps.New(9, 5, 2, 7)
	for i, v := range h.Drain() {
		fmt.Println(i, v)
	}

	// Output:
	// 0 2
	// 1 5
	// 2 7
	// 3 9
}
