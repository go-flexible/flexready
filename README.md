# flexready

<a href="https://pkg.go.dev/github.com/go-flexible/flexready"><img src="https://pkg.go.dev/badge/github.com/go-flexible/flexready.svg" alt="Go Reference"></a>
[![Go](https://github.com/go-flexible/flexready/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/go-flexible/flexready/actions/workflows/go.yml)

A [flex](https://github.com/go-flexible/flex) compatible readiness server.

## Install

```shell
go get github.com/go-flexible/flexready
```

## Example

```go
// Prepare your readyserver.
readysrv := flexready.New(flexready.Checks{
        "redis":       func() error { return redisCheck(nil) },
        "cockroachdb": func() error { return cockroachCheck(nil) },
}, flexready.WithAddress(":9999"))

// Run it, or better yet, let `flex` run it for you!
_ = readysrv.Run(context.Background())

// Ready server is now available on http://localhost:9999/ready
```
