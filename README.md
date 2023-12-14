# Missing abstract data types for Go

We've already got generic type-safe [slices](https://pkg.go.dev/slices) and [maps](https://pkg.go.dev/maps) since Go 1.21, but we want more.

Sub packages provided here share some common patterns:

- no third-party dependencies (except `golang.org/x` packages)
- (most) methods could be chained:

    ```go
    sets.New[float64](1,2,3,4,5).Union(sets.New(1,3,5,7,9)).Unset(1).Map(math.Sqrt)
    heaps.New[float64]().Reverse().Push(9).Push(1)
    ```

- provides `Clone` methods to get a deep copy of the original one,
- provides experimental [Iterator]s, see below for usage and development guides.

## Roadmap

- [x] set
- [x] (type safe) heap
- [ ] (type safe) list
- [ ] (type safe) ring
- [ ] (type safe) stack?
- [ ] (type safe) queue?
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
