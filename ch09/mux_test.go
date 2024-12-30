package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

/*#1*/
func drainAndClose(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			/*#2*/ next.ServeHTTP(w, r)
			_, _ = io.Copy(io.Discard, r.Body)
			_ = r.Body.Close()
		},
	)
}

func TestSimpleMux(t *testing.T) {
	serveMux := http.NewServeMux() // 멀티 플렉서 생성
	/*#3*/ serveMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	serveMux.HandleFunc( /*#4*/ "/hello", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, "Hello friend.")
	})

	serveMux.HandleFunc( /*#5*/ "/hello/there/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, "why hello there")
	})

	mux := drainAndClose(serveMux)

	testCases := []struct {
		path     string
		response string
		code     int
	}{
		/*#1*/ {"http://test/", "", http.StatusNoContent},
		{"http://test/hello", "Hello friend.", http.StatusOK},
		{"http://test/hello/there/", "why, hello there.", http.StatusOK},
		/*#2*/ {"http://test/hello/there", "<a href=\"/hello/there/\">Moved Permanently</a>.\n\n", http.StatusMovedPermanently},
		/*#3*/ {"http://test/hello/there/you", "Why, hello there.", http.StatusOK},
		/*#4*/ {"http://test/hello/and/goodbye", "", http.StatusNoContent},
		{"http://test/something/else/entirely", "", http.StatusNoContent},
		{"http://test/hello/you", "", http.StatusNoContent},
	}
	for i, c := range testCases {
		r := httptest.NewRequest(http.MethodGet, c.path, nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, r)
		resp := w.Result()

		if actual := resp.StatusCode; c.code != actual {
			t.Errorf("%d: expected code %d; actual %d", i, c.code, actual)
		}

		b, err := /*#5*/ io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}

		_ = /*#6*/ resp.Body.Close()

		if actual := string(b); c.response != actual {
			t.Errorf("%d: expected response %q; actual %q", i, c.response, actual)
		}
	}
}
