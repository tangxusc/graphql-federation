package main

import (
	"github.com/higress-group/proxy-wasm-go-sdk/proxywasm"
	"github.com/higress-group/wasm-go/pkg/wrapper"
	"github.com/tangxusc/higress-graphql-federation/pkg/config"
	_ "github.com/tangxusc/higress-graphql-federation/pkg/process"
)

func main() {}

func init() {
	proxywasm.LogDebugf("higress-graphql-federation init,process count:%+v", len(config.Options))
	wrapper.SetCtx("higress-graphql-federation", config.Options...)
	proxywasm.LogDebug("higress-graphql-federation init end")
}
