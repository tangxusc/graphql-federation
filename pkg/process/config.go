package process

import (
	"github.com/tangxusc/higress-graphql-federation/pkg/config"

	"github.com/higress-group/wasm-go/pkg/wrapper"
)

func init() {
	ops := wrapper.ProcessRequestBody(processRequestBody)
	config.RegisterCtxOption(ops)
	config.RegisterCtxOption(wrapper.ProcessRequestHeaders(onHttpRequestHeaders))
}
