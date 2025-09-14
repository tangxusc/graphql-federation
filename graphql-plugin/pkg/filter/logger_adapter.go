package filter

import (
	"github.com/envoyproxy/envoy/contrib/golang/common/go/api"
	"github.com/jensneuse/abstractlogger"
)

type GraphqlFederationLoggerAdapter struct {
}

func (g *GraphqlFederationLoggerAdapter) Println(v ...interface{}) {
	api.LogDebugf("%v", v)
}

func (g *GraphqlFederationLoggerAdapter) Printf(format string, v ...interface{}) {
	api.LogDebugf(format, v...)
}

func (g *GraphqlFederationLoggerAdapter) Debug(msg string, fields ...abstractlogger.Field) {
	temp := make([]interface{}, len(fields), len(fields))
	for i, field := range fields {
		temp[i] = field
	}
	api.LogDebugf(msg, temp...)
}

func (g *GraphqlFederationLoggerAdapter) Info(msg string, fields ...abstractlogger.Field) {
	temp := make([]interface{}, len(fields), len(fields))
	for i, field := range fields {
		temp[i] = field
	}
	api.LogInfof(msg, temp...)
}

func (g *GraphqlFederationLoggerAdapter) Warn(msg string, fields ...abstractlogger.Field) {
	temp := make([]interface{}, len(fields), len(fields))
	for i, field := range fields {
		temp[i] = field
	}
	api.LogWarnf(msg, temp...)
}

func (g *GraphqlFederationLoggerAdapter) Error(msg string, fields ...abstractlogger.Field) {
	temp := make([]interface{}, len(fields), len(fields))
	for i, field := range fields {
		temp[i] = field
	}
	api.LogErrorf(msg, temp...)
}

func (g *GraphqlFederationLoggerAdapter) Fatal(msg string, fields ...abstractlogger.Field) {
	temp := make([]interface{}, len(fields), len(fields))
	for i, field := range fields {
		temp[i] = field
	}
	api.LogErrorf(msg, temp...)
}

func (g *GraphqlFederationLoggerAdapter) Panic(msg string, fields ...abstractlogger.Field) {
	temp := make([]interface{}, len(fields), len(fields))
	for i, field := range fields {
		temp[i] = field
	}
	api.LogErrorf(msg, temp...)
}

func (g *GraphqlFederationLoggerAdapter) LevelLogger(level abstractlogger.Level) abstractlogger.LevelLogger {
	return g
}
