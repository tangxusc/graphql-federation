package main

import (
	jj "encoding/json"

	"github.com/higress-group/proxy-wasm-go-sdk/proxywasm"
	"github.com/higress-group/proxy-wasm-go-sdk/proxywasm/types"
	logs "github.com/higress-group/wasm-go/pkg/log"
	"github.com/higress-group/wasm-go/pkg/wrapper"
	"github.com/tidwall/gjson"
)

func main() {}

func init() {
	wrapper.SetCtx(
		"higress-graphql-federation",
		wrapper.ParseConfig(parseConfigNew),
		wrapper.ProcessRequestBody(func(context wrapper.HttpContext, config MyConfig, body []byte) types.Action {
			return types.ActionContinue
		}),

		// 为处理请求头，设置自定义函数
		wrapper.ProcessRequestHeadersBy(onHttpRequestHeaders),
	)
}
func body(ctx wrapper.HttpContext, config *MyConfig, body []byte) types.Action {
	return types.ActionContinue
}

func parseConfigNew(json gjson.Result, config *MyConfig) error {
	return nil
}

// 自定义插件配置
type MyConfig struct {
	mockEnable bool
}

// 在控制台插件配置中填写的yaml配置会自动转换为json，此处直接从json这个参数里解析配置即可
func parseConfig(json gjson.Result, config *MyConfig, log logs.Log) error {
	// 解析出配置，更新到config中
	config.mockEnable = json.Get("mockEnable").Bool()
	return nil
}

func onHttpRequestHeaders(ctx wrapper.HttpContext, config MyConfig, log logs.Log) types.Action {
	d, err := jj.Marshal(config)
	if err != nil {
		log.Error(err.Error())

	} else {
		log.Infof("config:", string(d))

	}
	log.Error("onHttpRequestHeaders")
	proxywasm.AddHttpRequestHeader("hello", "world")
	if config.mockEnable {
		proxywasm.SendHttpResponse(200, nil, []byte("hello world"), -1)
	}
	return types.HeaderContinue
}
