package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/golang_network/ch09/handler"
)

func TestSimpleHTTPServer(t *testing.T) {
	srv := &http.Server{
		Addr: "127.0.0.1:8081",
		Handler:/*#2*/ http.TimeoutHandler(handler.DefaultHandler(), 2*time.Minute, ""),
		IdleTimeout:       5 * time.Minute,
		ReadHeaderTimeout: time.Minute,
	}

	l, err := /*#2*/ net.Listen("tcp", srv.Addr)
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		err := /*#3*/ srv.Serve(l)
		if err != http.ErrServerClosed {
			t.Error(err)
		}
	}()

	testCases := []struct {
		method   string
		body     io.Reader
		code     int
		response string
	}{
		/*#1*/ {http.MethodGet, nil, http.StatusOK, "Hello, friend!"},
		/*#2*/ {http.MethodPost, bytes.NewBufferString("<world>"), http.StatusOK, "Hello, &lt;world&gt;!"},
		/*#3*/ {http.MethodHead, nil, http.StatusMethodNotAllowed, ""},
	}
	client := new(http.Client)
	path := fmt.Sprintf("http://%s/", srv.Addr)

	//////////////////////////////////////////////////////+++++++++++++++++++++++++++++++++++++++++++

	for i, c := range testCases {
		r, err := /*#1*/ http.NewRequest(c.method, path, c.body)
		if err != nil {
			t.Errorf("%d: %v", i, err)
			continue
		}

		resp, err := /*#2*/ client.Do(r)
		if err != nil {
			t.Errorf("%d: %v", i, err)
			continue
		}

		if resp.StatusCode != c.code {
			t.Errorf("%d: unexpected status code: %q", i, resp.Status)
		}

		b, err := /*#3*/ io.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("%d: %v", i, err)
			continue
		}

		_ = /*#4*/ resp.Body.Close()

		if c.response != string(b) {
			t.Errorf("%d: expected %q; actual %q", i, c.response, b)
		}
	}
	if err := /*#5*/ srv.Close(); err != nil {
		t.Fatal(err)
	}

}
