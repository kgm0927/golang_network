package handler

import (
	"database/sql"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"sort"
	"strings"
)

type Methods map[string]http.Handler

func (h Methods) /*#2*/ ServeHTTP(w http.ResponseWriter, r *http.Request) {
	/*#3*/ defer func(r io.ReadCloser) {
		_, _ = io.Copy(io.Discard, r)
		_ = r.Close()
	}(r.Body)

	if handler, ok := h[r.Method]; ok {
		if handler == nil {
			/*#4*/ http.Error(w, "Internal server error", http.StatusInternalServerError)
		} else {
			/*#5*/ handler.ServeHTTP(w, r)
		}
		return
	}

	/*#6*/
	w.Header().Add("Allow", h.allowedMethods())
	if r.Method != /*#7*/ http.MethodOptions {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h Methods) allowedMethods() string {
	a := make([]string, 0, len(h))

	for k := range h {
		a = append(a, k)
	}
	sort.Strings(a)
	return strings.Join(a, ", ")
}

func DefaultMethodHandler() http.Handler {
	return Methods{
		/*#1*/ http.MethodGet: http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte("Hello, friend!"))
			},
		),

		/*#2*/ http.MethodPost: http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				b, err := io.ReadAll(r.Body)
				if err != nil {
					http.Error(w, "Internal server error", http.StatusInternalServerError)
					return
				}
				_, _ = fmt.Fprintf(w, "Hello, %s!", html.EscapeString(string(b)))
			},
		),
	}
}

type Handlers struct {
	db *sql.DB
	/*#1*/ log *log.Logger
}

func (h *Handlers) Handler1() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			err := h.db.Ping()
			if err != nil {
				/*#2*/ h.log.Printf("db ping: %v", err) // 객체에 접근
			}
			// 데이터베이스와 관련된 작업 수행
		},
	)
}

func (h *Handlers) Handler2() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// ...
		},
	)
}
