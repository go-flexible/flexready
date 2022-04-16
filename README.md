# flexready

<a href="https://pkg.go.dev/github.com/go-flexible/flexready"><img src="https://pkg.go.dev/badge/github.com/go-flexible/flexready.svg" alt="Go Reference"></a>
[![Go](https://github.com/go-flexible/flexready/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/go-flexible/flexready/actions/workflows/go.yml)

A [flex](https://github.com/go-flexible/flex) compatible readiness server.

## Install

```shell
go get github.com/go-flexible/flexready
```
## Configuration

The readiness server can be configured through the environment to match setup in
the infrastructure.

- `FLEX_READYSRV_ADDR` default: `0.0.0.0:3674`
- `FLEX_READYSRV_LIVENESS_PATH` default: `/live`
- `FLEX_READYSRV_READINESS_PATH` default: `/ready`

## Example

```go
// Prepare your readyserver.
readysrv := flexready.New(flexready.Checks{
        "redis":       func() error { return redisCheck(nil) },
        "cockroachdb": func() error { return cockroachCheck(nil) },
}, flexready.WithAddress(":9999"))

// Run it, or better yet, let `flex` run it for you!
_ = readysrv.Run(context.Background())

// Liveness endpoint:  http://localhost:9999/live
// Readiness endpoint: http://localhost:9999/ready
```
