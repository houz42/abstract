[![Go Reference](https://pkg.go.dev/badge/github.com/houz42/abstract.svg)](https://pkg.go.dev/github.com/houz42/abstract)

# Missing abstract data types for Go

We've already got generic type-safe [slices](https://pkg.go.dev/slices) and [maps](https://pkg.go.dev/maps) since Go 1.21, but we want more.

Sub packages provided here share some common patterns:

- no third-party dependencies (except `golang.org/x` packages),
- provides alternative constructors for custom types:
  - [`heaps.NewFunc[MyType](myLessFunc)`](https://pkg.go.dev/github.com/houz42/abstract@v0.0.0-20231225123224-3c21759614ba/heaps#NewFunc)
  - [`skiplists.NewFunc[MyType](myCmpFunc)`](https://pkg.go.dev/github.com/houz42/abstract@v0.0.0-20231225123224-3c21759614ba/skiplists#NewFunc)
- provides `Reverse` methods for ordered sequences, to be used for reversion ordinary
  - [`heaps.New[string].Reverse()`](https://pkg.go.dev/github.com/houz42/abstract@v0.0.0-20231225123224-3c21759614ba/heaps#Heap.Reverse)
  - [`skiplists.New[int].Reverse()`](https://pkg.go.dev/github.com/houz42/abstract@v0.0.0-20231225123224-3c21759614ba/skiplists#SkipList.Reverse)
- (most) methods could be chained:
  - `sets.New[float64](1,2,3,4,5).Union(sets.New(1,3,5,7,9)).Unset(1).Map(math.Sqrt)`
  - `heaps.New[float64]().Reverse().Push(9).Push(1)`
- provides `Clone` methods to get a deep copy of the original one,
- provides experimental [Iterator]s, see below for usage and development guides.

## Roadmap

- [x] set
- [x] heap
- [x] list
- [x] skip list
- [ ] ring
- [ ] stack?
- [ ] queue?
- [ ] chainable maps?
- [ ] chainable slices?

## Usage

`go get github.com/houz42/abstract@latest`

If you want to try the experimental [iterator] and [range func] features,
set environment variable `GOEXPERIMENT=rangefunc` before running any `go` command:

```sh
GOEXPERIMENT=rangefunc go install my/program
GOEXPERIMENT=rangefunc go build my/program
GOEXPERIMENT=rangefunc go test my/program
GOEXPERIMENT=rangefunc go test my/program -bench=.
```

See the [range func] wiki for more details.

## Development

To develop with the experimental features, follow the [gist] to
install [gotip] and least gopls and configure vs code.

[range func]: https://github.com/golang/go/wiki/RangefuncExperiment
[iterator]: https://github.com/golang/go/issues/61897
[gotip]: https://pkg.go.dev/golang.org/dl/gotip
[gist]: https://gist.github.com/nikgalushko/e1b5c85c64653dd554a7a904bbef4eee
