package main

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func blockIndefinitely(w http.ResponseWriter, r *http.Request) {
	select {}
}
func TestBlockIndefinitely(t *testing.T) {
	ts := /*#1*/ httptest.NewServer( /*#2*/ http.HandlerFunc( /*#3*/ blockIndefinitely))
	_, _ = http.Get( /*#4*/ ts.URL)
	t.Fatal("client did not indefinitely block")
}

func TestBlockIndefinitelyWithTimeout(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(blockIndefinitely))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := /*#1*/ http.NewRequestWithContext(ctx, http.MethodGet, ts.URL, nil)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Fatal(err)
		}
		return
	}

	_ = resp.Body.Close()
}
