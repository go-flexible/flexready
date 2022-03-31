package main

import (
	"context"
	"errors"

	"github.com/go-flexible/flexready"
)

func main() {
	// Prepare your readyserver.
	readysrv := flexready.New(flexready.Checks{
		"redis":       func() error { return redisCheck(nil) },
		"cockroachdb": func() error { return cockroachCheck(nil) },
	}, flexready.WithAddress(":9999"))

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
