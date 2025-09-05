package process

import (
	"github.com/tangxusc/higress-graphql-federation/pkg/config"

	"github.com/higress-group/proxy-wasm-go-sdk/proxywasm/types"
	"github.com/higress-group/wasm-go/pkg/wrapper"
)

func processRequestBody(ctx wrapper.HttpContext, cfg config.FederationConfig, body []byte) types.Action {
	return types.ActionContinue
}
