# flexready

A [flex](https://github.com/go-flexible/flex) compatible readiness server.

```go
// Configure the server, or pass nil for package defaults.
config := &flexready.Config{
        Server: &http.Server{Addr: ":9999"},
}
// Prepare your readyserver.
readysrv := flexready.New(config, flexready.Checks{
        "redis":       func() error { return redisCheck(redisClient) },
        "cockroachdb": func() error { return cockroachCheck(dbClient) },
})
// Run it, or better yet, let `flex` run it for you!
_ = readysrv.Run(context.Background())

// Ready server is now available on http://localhost:9999/ready
```
