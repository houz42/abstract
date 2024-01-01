package heaps_test

import (
	"fmt"

	"github.com/houz42/abstract/heaps"
)

func Example() {
	h := heaps.New(2, 1, 5, 6)
	h.Push(3)
	h.RemoveAt(3)
	fmt.Printf("minimum: %d\n", h.Top())
	for h.Len() > 0 {
		fmt.Printf("%d ", h.Pop())
	}

	// Output:
	// minimum: 1
	// 1 2 3 5
}

func Example_priorityQueue() {
	type process struct {
		pid      int
		niceness int
	}

	queue := heaps.NewFunc[process](func(x, y process) bool { return x.niceness < y.niceness })

	queue.Push(process{pid: 1, niceness: -20})
	queue.Push(process{pid: 2, niceness: 0})
	queue.Push(process{pid: 3, niceness: 10})
	queue.Push(process{pid: 4, niceness: -1})

	for queue.Len() > 0 {
		p := queue.Pop()
		fmt.Printf("start process %d with niceness %d\n", p.pid, p.niceness)
	}

	// Output:
	// start process 1 with niceness -20
	// start process 4 with niceness -1
	// start process 2 with niceness 0
	// start process 3 with niceness 10
}

func ExampleHeap_RemoveAt() {
	h := heaps.New(1, 5, 3, 2)
	fmt.Println("removed:", h.RemoveAt(2))

	for h.Len() > 0 {
		fmt.Println(h.Pop())
	}

	// Output:
	// removed: 3
	// 1
	// 2
	// 5
}

func ExampleHeap_Reverse() {
	type plan struct {
		name     string
		severity int
	}

	queue := heaps.NewFunc[plan](func(x, y plan) bool { return x.severity < y.severity }).Reverse()

	queue.Push(plan{name: "call friend", severity: 1})
	queue.Push(plan{name: "finish report", severity: 3})
	queue.Push(plan{name: "buy food", severity: 2})

	for queue.Len() > 0 {
		plan := queue.Pop()
		fmt.Printf("%d: %s\n", plan.severity, plan.name)
	}

	// Output:
	// 3: finish report
	// 2: buy food
	// 1: call friend
}
