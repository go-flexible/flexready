package main

import (
	"context"
	"errors"
	"io"
	"log"

	"github.com/go-flexible/flexready"
)

func main() {
	logger := log.New(io.Discard, "", 0)

	// Prepare your readyserver.
	readysrv := flexready.New(flexready.Checks{
		"redis":       func() error { return redisCheck(nil) },
		"cockroachdb": func() error { return cockroachCheck(nil) },
	}, flexready.WithAddress(":9999"), flexready.WithLogger(logger))

	// Run it, or better yet, let `flex` run it for you!
	_ = readysrv.Run(context.Background())
}

// redis is broken.
func redisCheck(redisClient any) error {
	return errors.New("connection to redis is broken")
}

// all is well.
func cockroachCheck(dbClient any) error {
	return nil
}
