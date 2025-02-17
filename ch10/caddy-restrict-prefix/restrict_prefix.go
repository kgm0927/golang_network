package restrictprefix

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"go.uber.org/zap"
)

func init() {
	/*#1*/ caddy.RegisterModule(RestrictPrefix{})
}

// RestrictPrefix은 URI의 일부가 주어진 접두사와 일치하는 요청을 제한하는 미들웨어
type RestrictPrefix struct {
	Prefix string      `json:"prefix,omitempty"` /*#2*/
	logger *zap.Logger /*#3*/
}

// CaddyModule은 Caddy의 모듈 정보를 반환함.
func (RestrictPrefix) /*#4*/ CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		/*#5*/ ID: "http.handlers.restrict_prefix",
		/*#6*/ New: func() caddy.Module { return new(RestrictPrefix) },
	}
}

// Zap 로거를 RestrictPrefix로 프로비저닝
func (p *RestrictPrefix) /*#1*/ Provision(ctx caddy.Context) error {
	p.logger = /*#2*/ ctx.Logger(p)
	return nil
}

// 모듈 구성에서 접두사를 검증하고 필요시 기본 접두사 "."으로 설정
func (p *RestrictPrefix) /*#3*/ Validate() error {
	if p.Prefix == "" {
		p.Prefix = "."
	}
	return nil
}

// serverHTTP는 caddyhttp.MiddlewareHandler 인터페이스 구현
func (p RestrictPrefix) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	/*#1*/ for _, part := range strings.Split(r.URL.Path, "/") {
		if strings.HasPrefix(part, p.Prefix) {
			/*#2*/ http.Error(w, "Not Found", http.StatusNotFound)
			if p.logger != nil {
				/*#3*/ p.logger.Debug(fmt.Sprintf("restricted prefix: %q in %s", part, r.URL.Path))
			}
			return nil
		}
	}
	return /*#4*/ next.ServeHTTP(w, r)
}

var (
	/*#5*/ _ caddy.Provisioner           = (*RestrictPrefix)(nil)
	_        caddy.Validator             = (*RestrictPrefix)(nil)
	_        caddyhttp.MiddlewareHandler = (*RestrictPrefix)(nil)
)
