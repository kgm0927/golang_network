package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandlerWriteHeader(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		_, _ = /*#1*/ w.Write([]byte("Bad request"))
		/*#2*/ w.WriteHeader(http.StatusBadRequest)
	}
	r := httptest.NewRequest(http.MethodGet, "http://test", nil)
	w := httptest.NewRecorder()
	handler(w, r)
	t.Logf("Response status: %q" /*#3*/, w.Result().Status)

	handler = func(w http.ResponseWriter, r *http.Request) {
		/*#4*/ w.WriteHeader(http.StatusBadRequest)
		_, _ = /*#5*/ w.Write([]byte("Bad request"))
	}

	r = httptest.NewRequest(http.MethodGet, "http://test", nil)
	w = httptest.NewRecorder()
	handler(w, r)

	t.Logf("Response status: %q" /*#6*/, w.Result().Status)
}
