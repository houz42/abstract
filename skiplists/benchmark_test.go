package skiplists_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/houz42/abstract/skiplists"
)

func BenchmarkSkipList(b *testing.B) {
	for size := 1000; size < 1_000_000; size *= 10 {
		b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
			basicReadWrite(b, size)
			randomAccess(b, size)
		})
	}
}

func newList(size int) *skiplists.SkipList[int] {
	list := skiplists.New[int]()
	for i := 0; i < size; i++ {
		list.Set(i)
	}
	return list
}

func basicReadWrite(b *testing.B, size int) {
	perm := rand.Perm(size)

	random := make([]int, 0, size)
	for i := 0; i < size; i++ {
		random = append(random, rand.Intn(size))
	}

	b.ResetTimer()

	b.Run("set", func(b *testing.B) {
		b.Run("forward", func(b *testing.B) {
			for x := 0; x < b.N; x++ {
				list := skiplists.New[int]()
				for i := 0; i < size; i++ {
					list.Set(i)
				}
				if list.Len() != size {
					b.Fatal()
				}
			}
		})

		b.Run("backward", func(b *testing.B) {
			for x := 0; x < b.N; x++ {
				list := skiplists.New[int]()
				for i := 0; i < size; i++ {
					list.Set(size - i)
				}
				if list.Len() != size {
					b.Fatal()
				}
			}
		})

		b.Run("random", func(b *testing.B) {
			for x := 0; x < b.N; x++ {
				list := skiplists.New[int]()
				for _, i := range random {
					list.Set(i)
				}
				if list.Len() > size {
					b.Fatal()
				}
			}
		})
	})

	b.Run("get", func(b *testing.B) {
		list := newList(size)
		b.ResetTimer()

		b.Run("forward", func(b *testing.B) {
			var v int
			var ok bool
			for x := 0; x < b.N; x++ {
				for i := 0; i < size; i++ {
					v, ok = list.Get(i)
					if !ok || v != i {
						b.Fatal()
					}
				}
			}
		})

		b.Run("backward", func(b *testing.B) {
			var v int
			var ok bool
			for x := 0; x < b.N; x++ {
				for i := size - 1; i >= 0; i-- {
					v, ok = list.Get(i)
					if !ok || v != i {
						b.Fatal()
					}
				}
			}
		})

		b.Run("random", func(b *testing.B) {
			var v int
			var ok bool
			for x := 0; x < b.N; x++ {
				for _, i := range perm {
					v, ok = list.Get(i)
					if !ok || v != i {
						b.Fatal()
					}
				}
			}
		})
	})

	b.Run("remove", func(b *testing.B) {
		b.Run("forward", func(b *testing.B) {
			for x := 0; x < b.N; x++ {
				list := newList(size)
				// b.ResetTimer()

				for i := 0; i < size; i++ {
					list.Unset(i)
				}
				if list.Len() != 0 {
					b.Fatal()
				}
			}
		})

		b.Run("backward", func(b *testing.B) {
			for x := 0; x < b.N; x++ {
				list := newList(size)
				// b.ResetTimer()

				for i := size - 1; i >= 0; i-- {
					list.Unset(i)
				}
				if list.Len() != 0 {
					b.Fatal()
				}
			}
		})

		b.Run("random", func(b *testing.B) {
			for x := 0; x < b.N; x++ {
				list := newList(size)
				// b.ResetTimer()

				for _, i := range perm {
					list.Unset(i)
				}
				if list.Len() != 0 {
					b.Fatal()
				}
			}
		})
	})

}

func randomAccess(b *testing.B, size int) {
	perm := rand.Perm(size)
	b.ResetTimer()

	b.Run("at", func(b *testing.B) {
		list := newList(size)
		b.ResetTimer()

		b.Run("forward", func(b *testing.B) {
			for x := 0; x < b.N; x++ {
				for i := 0; i < size; i++ {
					if list.At(i) != i {
						b.Fatal()
					}
				}
			}
		})

		b.Run("backward", func(b *testing.B) {
			for x := 0; x < b.N; x++ {
				for i := size - 1; i >= 0; i-- {
					if list.At(i) != i {
						b.Fatal()
					}
				}
			}
		})

		b.Run("random", func(b *testing.B) {
			for x := 0; x < b.N; x++ {
				for _, i := range perm {
					if list.At(i) != i {
						b.Fatal()
					}
				}
			}
		})

	})

	b.Run("remove at", func(b *testing.B) {
		b.Run("forward", func(b *testing.B) {
			for x := 0; x < b.N; x++ {
				list := newList(size)
				// b.ResetTimer()

				for i := 0; i < size; i++ {
					list.RemoveAt(0)
				}

				if list.Len() != 0 {
					b.Fatal()
				}
			}
		})

		b.Run("backward", func(b *testing.B) {
			for x := 0; x < b.N; x++ {
				list := newList(size)
				// b.ResetTimer()

				for i := size - 1; i >= 0; i-- {
					list.RemoveAt(list.Len() - 1)
				}

				if list.Len() != 0 {
					b.Fatal()
				}
			}
		})

		b.Run("random", func(b *testing.B) {
			for x := 0; x < b.N; x++ {
				list := newList(size)
				// b.ResetTimer()

				for _, i := range perm {
					list.RemoveAt(i % list.Len())
				}

				if list.Len() != 0 {
					b.Fatal()
				}
			}
		})

	})
}
