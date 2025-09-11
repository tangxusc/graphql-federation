package config

import (
	"context"

	"github.com/higress-group/proxy-wasm-go-sdk/proxywasm"
	"github.com/higress-group/wasm-go/pkg/wrapper"
	"github.com/tidwall/gjson"
)

func init() {
	RegisterCtxOption(wrapper.ParseConfig(parseFederationConfig))
}

type FederationConfig struct {
	DebugMode  bool `json:"debugMode"`
	Logger     WasmLogger
	CancelFunc context.CancelFunc
	Ctx        context.Context
}

func parseFederationConfig(json gjson.Result, config *FederationConfig) error {
	config.DebugMode = json.Get("debugMode").Bool()

	config.Logger = WasmLogger{}

	config.Ctx, config.CancelFunc = context.WithCancel(context.TODO())
	proxywasm.LogDebugf("higress-graphql-federation init,tickers count:%+v", len(Tickers))
	for _, k := range Tickers {
		k(config)
	}
	return nil
}
