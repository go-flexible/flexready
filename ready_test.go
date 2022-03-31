package flexready_test

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/go-flexible/flexready"
)

func ExampleNew() {
	checks := flexready.Checks{
		"google": func() error {
			_, err := http.Get("https://google.com")
			return err
		},
	}
	flexready.New(checks)
}

func ExampleNew_withOptions() {
	checks := flexready.Checks{
		"google": func() error {
			_, err := http.Get("https://google.com")
			return err
		},
	}

	logger := log.New(io.Discard, "", 0)
	flexready.New(checks, flexready.WithLogger(logger))
}

type health struct {
	Messages string `json:"messages"`
	OK       bool   `json:"ok"`
}
type response map[string]health

func TestNew(t *testing.T) {
	t.Run("nil parameters must not return empty server", func(t *testing.T) {
		srv := flexready.New(nil)
		notEqual(t, srv, nil)
	})
	t.Run("checkers are being called", func(t *testing.T) {
		srv := flexready.New(flexready.Checks{
			"ok":     func() error { return nil },
			"not_ok": func() error { return errors.New("oops") },
		}, flexready.WithAddress(":0"))

		req := httptest.NewRequest(http.MethodGet, "/ready", nil)
		rec := httptest.NewRecorder()
		srv.Server.Handler.ServeHTTP(rec, req)

		var res response
		if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
			t.Fatal(err)
		}

		equal(t, res["ok"].OK, true)
		equal(t, res["ok"].Messages, "")
		equal(t, res["not_ok"].OK, false)
		equal(t, res["not_ok"].Messages, "oops")
	})
}

func equal(t *testing.T, got, want interface{}) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got: %#[1]v (%[1]T), but wanted: %#[2]v (%[2]T)", got, want)
	}
}

func notEqual(t *testing.T, got, want interface{}) {
	t.Helper()
	if reflect.DeepEqual(got, want) {
		t.Fatalf("got: %#[1]v (%[1]T), but wanted: %#[2]v (%[2]T)", got, want)
	}
}
