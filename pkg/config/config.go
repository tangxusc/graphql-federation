package config

import (
	"github.com/higress-group/wasm-go/pkg/wrapper"
)

var Options []wrapper.CtxOption[FederationConfig]=make([]wrapper.CtxOption[FederationConfig], 0)

func RegisterCtxOption(opt wrapper.CtxOption[FederationConfig]) {
	Options = append(Options, opt)
}
