package process

import (
	"github.com/tangxusc/higress-graphql-federation/pkg/config"

	"github.com/higress-group/proxy-wasm-go-sdk/proxywasm"
	"github.com/higress-group/proxy-wasm-go-sdk/proxywasm/types"
	"github.com/higress-group/wasm-go/pkg/wrapper"
)

func onHttpRequestHeaders(ctx wrapper.HttpContext, cfg config.FederationConfig) types.Action {
	proxywasm.LogDebug("onHttpRequestHeaders...")
	//proxywasm.AddHttpRequestHeader("hello", "world")
	return types.HeaderStopIteration
}
