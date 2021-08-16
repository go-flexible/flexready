package main

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-flexible/flexready"
)

func main() {
	// Configure the server, or pass nil for package defaults.
	config := &flexready.Config{
		Server: &http.Server{Addr: ":9999"},
	}

	// Prepare your readyserver.
	readysrv := flexready.New(config, flexready.Checks{
		"redis":       func() error { return redisCheck(nil) },
		"cockroachdb": func() error { return cockroachCheck(nil) },
	})

	// Run it, or better yet, let `flex` run it for you!
	_ = readysrv.Run(context.Background())
}

// redis is broken.
func redisCheck(redisClient interface{}) error {
	return errors.New("connection to redis is broken")
}

// all is well.
func cockroachCheck(dbClient interface{}) error {
	return nil
}
