package main

import (
	_ "caddy/caddy-restrict-prefix"
	_ "caddy/caddy-toml-adapter"

	cmd "github.com/caddyserver/caddy/v2/cmd"
	_ "github.com/caddyserver/caddy/v2/modules/standard"
)

func main() {
	cmd.Main()

}
