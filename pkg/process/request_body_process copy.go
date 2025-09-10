package process

// import (
// 	"bytes"
// 	"context"
// 	"encoding/json"
// 	"errors"
// 	"net/http"

// 	"github.com/higress-group/proxy-wasm-go-sdk/proxywasm"
// 	"github.com/tangxusc/higress-graphql-federation/pkg/config"
// 	"github.com/wundergraph/graphql-go-tools/execution/engine"
// 	"github.com/wundergraph/graphql-go-tools/execution/graphql"
// 	"github.com/wundergraph/graphql-go-tools/v2/pkg/engine/resolve"

// 	"github.com/higress-group/proxy-wasm-go-sdk/proxywasm/types"
// 	"github.com/higress-group/wasm-go/pkg/wrapper"
// )

// // GraphQLRequest 表示 GraphQL 请求
// type GraphQLRequest struct {
// 	Query         string                 `json:"query"`
// 	Variables     map[string]interface{} `json:"variables,omitempty"`
// 	OperationName string                 `json:"operationName,omitempty"`
// }

// func processRequestBody(ctx wrapper.HttpContext, cfg config.FederationConfig, body []byte) types.Action {
// 	//TODO:post method , content-type:graphql/json , path is /graphql
// 	var gqlRequest graphql.Request
// 	var err error
// 	if err = json.Unmarshal(body, &gqlRequest); err != nil {
// 		proxywasm.LogErrorf("Failed to unmarshal GraphQL request: %v", err)
// 		_ = proxywasm.SendHttpResponse(http.StatusBadRequest, nil, []byte(err.Error()), -1)
// 		return types.ActionContinue
// 	}
// 	var opts []engine.ExecutionOptions

// 	tracingOpts := resolve.TraceOptions{
// 		Enable:                                 true,
// 		ExcludePlannerStats:                    false,
// 		ExcludeRawInputData:                    false,
// 		ExcludeInput:                           false,
// 		ExcludeOutput:                          false,
// 		ExcludeLoadStats:                       false,
// 		EnablePredictableDebugTimings:          false,
// 		IncludeTraceOutputInResponseExtensions: true,
// 	}

// 	opts = append(opts, engine.WithRequestTraceOptions(tracingOpts))
// 	err = initGraphqlFederationEngine(cfg)
// 	if err != nil {
// 		proxywasm.LogErrorf("Failed to init GraphQL engine: %v", err)
// 		_ = proxywasm.SendHttpResponse(http.StatusInternalServerError, nil, []byte(err.Error()), -1)
// 		return types.ActionContinue
// 	}
// 	proxywasm.LogDebug("=====================")
// 	return types.ActionPause
// 	// return types.

// 	buf := bytes.NewBuffer(make([]byte, 0, 4096))
// 	resultWriter := graphql.NewEngineResultWriterFromBuffer(buf)
// 	if err = ExecutionEngine.Execute(context.TODO(), &gqlRequest, &resultWriter, opts...); err != nil {
// 		proxywasm.LogErrorf("Failed to execute GraphQL query: %v", err)
// 		_ = proxywasm.SendHttpResponse(http.StatusInternalServerError, nil, []byte(err.Error()), -1)
// 		return types.ActionContinue
// 	}
// 	proxywasm.LogDebug("=======================")
// 	proxywasm.LogDebug(string(resultWriter.String()))
// 	proxywasm.LogDebug("=======================")
// 	proxywasm.SendHttpResponse(http.StatusOK, nil, resultWriter.Bytes(), -1)

// 	return types.ActionPause
// }

// var QueryNotFoundError = errors.New("query not found")
// var ParseGraphqlError = errors.New("parse graphql error")
