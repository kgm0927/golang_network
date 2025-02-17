package tomladapter

import (
	"encoding/json"

	"github.com/caddyserver/caddy/v2/caddyconfig"
	"github.com/pelletier/go-toml"
)

func init() {
	caddyconfig.RegisterAdapter( /*#1*/ "toml", Adapter{})
}

// 어댑터는 TOML 형식의 Caddy 구성 파일을 JSON으로 변환
type Adapter struct{}

// TOML 형식의 보디를 JSON으로 변환
func (a Adapter) Adapt(body []byte, _ map[string]interface{}) ([]byte, []caddyconfig.Warning, error) {
	tree, err := /*#3*/ toml.LoadBytes(body)
	if err != nil {
		return nil, nil, err
	}
	b, err := json.Marshal( /*#4*/ tree.ToMap())

	return b, nil, err
}
