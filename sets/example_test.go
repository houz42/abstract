package sets_test

import (
	"fmt"

	"github.com/houz42/abstract/sets"
)

func Example() {

}

func ExampleSet_Map() {
	str := sets.Map(sets.New("hello", "gopher"), func(s string) int { return len(s) })
	for s := range str {
		fmt.Println(s)
	}

	// Unordered Output:
	// 5
	// 6
}

func ExampleSet_Filter() {
	even := sets.New(1, 2, 3, 4, 5).Filter(func(i int) bool { return i%2 == 0 })
	for i := range even {
		fmt.Println(i)
	}

	// Unordered Output:
	// 2
	// 4
}

func ExampleSet_Equal() {
	set := sets.New(1, 2, 3, 4, 5)
	fmt.Println(set.Equal(sets.New(5, 4, 3, 2, 1)))
	fmt.Println(set.Equal(sets.New(1, 2, 3)))

	// Output:
	// true
	// false
}

func ExampleSet_Clone() {
	set := sets.New(1, 2, 3)
	clone := set.Clone()
	for i := range clone {
		fmt.Println(i)
	}

	// Unordered Output:
	// 1
	// 2
	// 3
}

func ExampleSet_Subset() {
	set := sets.New(1, 2, 3)
	fmt.Println(set.Subset(sets.New(1, 2, 3, 4, 5)))
	fmt.Println(set.Subset(sets.New(1, 2, 3)))
	fmt.Println(set.Subset(sets.New(2, 4, 6)))

	// Output:
	// true
	// true
	// false
}

func ExampleSet_Superset() {
	set := sets.New(1, 2, 3, 4, 5)
	fmt.Println(set.Superset(sets.New[int]()))
	fmt.Println(set.Superset(sets.New(1, 2, 3)))
	fmt.Println(set.Superset(sets.New(1, 2, 3, 4, 5)))
	fmt.Println(set.Superset(sets.New(2, 4, 6)))

	// Output:
	// true
	// true
	// true
	// false
}

func ExampleSet_Union() {
	s := sets.New(1, 2, 3)
	t := sets.New(2, 4, 6)

	u := s.Union(t)
	for v := range u {
		fmt.Println(v)
	}

	// Unordered Output:
	// 1
	// 2
	// 3
	// 4
	// 6
}

func ExampleSet_Intersection() {
	s := sets.New(1, 2, 3, 4, 5)
	t := sets.New(2, 4, 6, 8)

	i := s.Intersection(t)
	for v := range i {
		fmt.Println(v)
	}

	// Unordered Output:
	// 2
	// 4
}
