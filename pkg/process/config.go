package process

import (
	"github.com/tangxusc/higress-graphql-federation/pkg/config"

	"github.com/higress-group/wasm-go/pkg/wrapper"
)

func init() {
	//config.RegisterCtxOption(wrapper.ProcessRequestHeaders(onHttpRequestHeaders))
	config.RegisterCtxOption(wrapper.ProcessRequestBody(processRequestBody))
}
