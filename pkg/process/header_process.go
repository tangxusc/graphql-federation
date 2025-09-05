package process

import (
	"encoding/json"

	"github.com/tangxusc/higress-graphql-federation/pkg/config"

	"github.com/higress-group/proxy-wasm-go-sdk/proxywasm"
	"github.com/higress-group/proxy-wasm-go-sdk/proxywasm/types"
	"github.com/higress-group/wasm-go/pkg/wrapper"
)

func onHttpRequestHeaders(ctx wrapper.HttpContext, cfg config.FederationConfig) types.Action {
	d, err := json.Marshal(cfg)
	if err != nil {
		proxywasm.LogError(err.Error())

	} else {
		proxywasm.LogInfof("config:%s \n", string(d))

	}
	proxywasm.LogError("onHttpRequestHeaders")
	proxywasm.AddHttpRequestHeader("hello", "world")
	// if cfg.MockEnable {
	// 	proxywasm.SendHttpResponse(200, nil, []byte("hello world"), -1)
	// }
	return types.HeaderContinue
}
