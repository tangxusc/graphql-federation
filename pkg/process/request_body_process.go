package process

import (
	"net/http"

	"github.com/higress-group/proxy-wasm-go-sdk/proxywasm"
	"github.com/tangxusc/higress-graphql-federation/pkg/config"

	"github.com/higress-group/proxy-wasm-go-sdk/proxywasm/types"
	"github.com/higress-group/wasm-go/pkg/wrapper"
)

func processRequestBody(ctx wrapper.HttpContext, cfg config.FederationConfig, body []byte) types.Action {
	proxywasm.LogDebug("+++++++++++++++++++++++++++")
	client := wrapper.NewClusterClient(wrapper.FQDNCluster{
		FQDN: "httpbin",
		Port: 80,
	})

	client.Call("POST", "/", nil, []byte{},
		func(statusCode int, responseHeaders http.Header, responseBody []byte) {
			proxywasm.LogDebugf("=================statusCode:%v,header:%+v,responseBody:%s", statusCode, responseHeaders, responseBody)
			proxywasm.ResumeHttpRequest()
		})
	proxywasm.LogDebug("执行请求--------------")
	return types.ActionPause
}
