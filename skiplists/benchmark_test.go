package skiplists_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/houz42/abstract/skiplists"
)

func Benchmark(b *testing.B) {
	size := 1000

	b.Run("forward insert", func(b *testing.B) {
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

	b.Run("forward delete", func(b *testing.B) {
		for x := 0; x < b.N; x++ {
			list := skiplists.New[int]()
			for i := 0; i < size; i++ {
				list.Set(i)
			}

			for i := 0; i < size; i++ {
				if list.Delete(i) != i {
					b.Fatal()
				}
			}
		}
	})
}

func BenchmarkSkipList(b *testing.B) {
	size := 100
	for i := 0; i < 2; i++ {
		size *= 10

		b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
			// basicReadWrite(b, size)
			randomAccess(b, size)
		})
	}
}

func basicReadWrite(b *testing.B, size int) {
	perm := rand.Perm(size)
	b.ResetTimer()

	b.Run("insert", func(b *testing.B) {

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
				for _, i := range perm {
					list.Set(i)
				}
				if list.Len() != size {
					b.Fatal()
				}
			}
		})
	})

	b.Run("search", func(b *testing.B) {

		b.Run("forward", func(b *testing.B) {
			list := skiplists.New[int]()
			for i := 0; i < size; i++ {
				list.Set(i)
			}
			b.ResetTimer()

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
			list := skiplists.New[int]()
			for i := 0; i < size; i++ {
				list.Set(i)
			}
			b.ResetTimer()

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
			list := skiplists.New[int]()
			for i := 0; i < size; i++ {
				list.Set(i)
			}
			b.ResetTimer()

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

	b.Run("delete", func(b *testing.B) {

		b.Run("forward", func(b *testing.B) {
			for x := 0; x < b.N; x++ {
				list := skiplists.New[int]()
				for i := 0; i < size; i++ {
					list.Set(i)
				}

				for i := 0; i < size; i++ {
					if list.Delete(i) != i {
						b.Fatal()
					}
				}
			}
		})

		b.Run("backward", func(b *testing.B) {
			for x := 0; x < b.N; x++ {
				list := skiplists.New[int]()
				for i := 0; i < size; i++ {
					list.Set(i)
				}

				for i := size - 1; i >= 0; i-- {
					if list.Delete(i) != i {
						b.Fatal()
					}
				}
			}
		})

		b.Run("random", func(b *testing.B) {
			for x := 0; x < b.N; x++ {
				list := skiplists.New[int]()
				for i := 0; i < size; i++ {
					list.Set(i)
				}

				for _, i := range perm {
					if list.Delete(i) != i {
						b.Fatal()
					}
				}
			}
		})
	})

}

func randomAccess(b *testing.B, size int) {
	perm := rand.Perm(size)
	b.ResetTimer()

	b.Run("at", func(b *testing.B) {
		b.Run("forward", func(b *testing.B) {
			for x := 0; x < b.N; x++ {
				list := skiplists.New[int]()
				for i := 0; i < size; i++ {
					list.Set(i)
				}

				// b.ResetTimer()

				for i := 0; i < size; i++ {
					if list.At(i) != i {
						b.Fatal()
					}
				}
			}
		})

		b.Run("barkward", func(b *testing.B) {
			for x := 0; x < b.N; x++ {
				list := skiplists.New[int]()
				for i := 0; i < size; i++ {
					list.Set(i)
				}

				for i := size - 1; i >= 0; i-- {
					if list.At(i) != i {
						b.Fatal()
					}
				}
			}
		})

		b.Run("random", func(b *testing.B) {
			for x := 0; x < b.N; x++ {
				list := skiplists.New[int]()
				for i := 0; i < size; i++ {
					list.Set(i)
				}
				for _, i := range perm {
					if list.At(i) != i {
						b.Fatal()
					}
				}
			}
		})

	})

	b.Run("delete at", func(b *testing.B) {
		b.Run("forward", func(b *testing.B) {
			for x := 0; x < b.N; x++ {
				list := skiplists.New[int]()
				for i := 0; i < size; i++ {
					list.Set(i)
				}

				for i := 0; i < size; i++ {
					if list.DeleteAt(0) != i {
						b.Fatal()
					}
				}

				if list.Len() != 0 {
					b.Fatal()
				}
			}
		})

		b.Run("backward", func(b *testing.B) {
			for x := 0; x < b.N; x++ {
				list := skiplists.New[int]()
				for i := 0; i < size; i++ {
					list.Set(i)
				}

				for i := size - 1; i >= 0; i-- {
					if list.DeleteAt(list.Len()-1) != i {
						b.Fatal()
					}
				}

				if list.Len() != 0 {
					b.Fatal()
				}
			}
		})

		b.Run("random", func(b *testing.B) {
			for x := 0; x < b.N; x++ {
				list := skiplists.New[int]()
				for i := 0; i < size; i++ {
					list.Set(i)
				}

				for _, i := range perm {
					list.DeleteAt(i % list.Len())
				}

				if list.Len() != 0 {
					b.Fatal()
				}
			}
		})
	})
}
