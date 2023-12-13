package sets_test

import (
	"math"
	"testing"

	"github.com/houz42/abstract/sets"
)

func TestSetValues(t *testing.T) {
	t.Run("from empty values", func(t *testing.T) {
		set := sets.New[int]()
		assertSize(t, set, 0)

		set.Set(1)
		set.Set(2)
		set.Set(3)
		set.Unset(2)

		assertSize(t, set, 2)
		assertElements(t, set, 1, 3)
	})

	t.Run("from existing values", func(t *testing.T) {
		set := sets.New(1, 2, 3)
		assertSize(t, set, 3)
		assertElements(t, set, 1, 2, 3)

		set.Set(9)
		set.Set(8)
		set.Set(1)
		set.Unset(2)
		set.Unset(7)

		assertSize(t, set, 4)
		assertElements(t, set, 1, 3, 8, 9)
	})
}

func TestEqual(t *testing.T) {
	t.Run("float Nan", func(t *testing.T) {
		s := sets.New[float64](0., 1., math.NaN())
		p := sets.New[float64](0., 1., math.NaN())
		if s.Equal(p) {
			t.Fatal()
		}
	})
}

func assertSize(t *testing.T, set sets.Set[int], want int) {
	if size := set.Size(); size != want {
		t.Fatalf("expecting size of set be %d, got %d", want, size)
	}
}

func assertElements(t *testing.T, set sets.Set[int], values ...int) {
	for _, v := range values {
		if !set.Contains(v) {
			t.Fatalf("%d should be in set but not found", v)
		}
	}
}
