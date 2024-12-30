package handler

import (
	"html/template"
	"io"
	"net/http"
)

var t = /*#1*/ template.Must(template.New("hello").Parse("Hello, {{.}}!"))

func DefaultHandler() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			/*#2*/ defer func(r io.ReadCloser) {
				_, _ = io.Copy(io.Discard, r)
				_ = r.Close()
			}(r.Body)

			var b []byte

			/*#3*/
			switch r.Method {
			case http.MethodGet:
				b = []byte("friend")
			case http.MethodPost:
				var err error
				b, err = io.ReadAll(r.Body)
				/*#4*/ if err != nil {
					http.Error(w, "Internal server error",
						http.StatusInternalServerError)
					return
				}
			default:
				// "Allow" 헤더가 없기 때문에 RFC 규격을 따르지 않음
				/*#5*/
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				return
			}
			_ = /*#6*/ t.Execute(w, string(b))
		},
	)
}
