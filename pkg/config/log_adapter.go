package config

import (
	"github.com/higress-group/proxy-wasm-go-sdk/proxywasm"
	"github.com/jensneuse/abstractlogger"
)

type WasmLogger struct {
}

func (w WasmLogger) Println(v ...interface{}) {
	proxywasm.LogDebugf("%v", v)
}

func (w WasmLogger) Printf(format string, v ...interface{}) {
	proxywasm.LogDebugf(format, v)
}

func (w WasmLogger) Debug(msg string, fields ...abstractlogger.Field) {
	proxywasm.LogDebugf(msg, fields)
}

func (w WasmLogger) Info(msg string, fields ...abstractlogger.Field) {
	proxywasm.LogInfof(msg, fields)
}

func (w WasmLogger) Warn(msg string, fields ...abstractlogger.Field) {
	proxywasm.LogWarnf(msg, fields)
}

func (w WasmLogger) Error(msg string, fields ...abstractlogger.Field) {
	proxywasm.LogErrorf(msg, fields)
}

func (w WasmLogger) Fatal(msg string, fields ...abstractlogger.Field) {
	proxywasm.LogErrorf(msg, fields)
}

func (w WasmLogger) Panic(msg string, fields ...abstractlogger.Field) {
	proxywasm.LogErrorf(msg, fields)
}

func (w WasmLogger) LevelLogger(level abstractlogger.Level) abstractlogger.LevelLogger {
	return WasmLogger{}
}
