package middleware

import (
	"net/http"
	"path"
	"strings"
)

// 온점으로 시작하는 파일, 디렉터리를 서빙하지 못하게 하는 예시 코드이다.
func RestrictPrefix(prefix string, next http.Handler) http.Handler {
	return /*#1*/ http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			/*#2*/ for _, p := range strings.Split(path.Clean(r.URL.Path), "/") {
				if strings.HasPrefix(p, prefix) {
					/*#3*/ http.Error(w, "Not Found", http.StatusNotFound)
				}
			}
			next.ServeHTTP(w, r)
		},
	)
}
