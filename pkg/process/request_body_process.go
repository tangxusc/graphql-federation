package process

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/higress-group/proxy-wasm-go-sdk/proxywasm"
	"github.com/tangxusc/higress-graphql-federation/pkg/config"
	"github.com/wundergraph/graphql-go-tools/execution/engine"
	"github.com/wundergraph/graphql-go-tools/execution/graphql"
	"github.com/wundergraph/graphql-go-tools/v2/pkg/astparser"
	"github.com/wundergraph/graphql-go-tools/v2/pkg/engine/resolve"

	"github.com/higress-group/proxy-wasm-go-sdk/proxywasm/types"
	"github.com/higress-group/wasm-go/pkg/wrapper"
)

// GraphQLRequest 表示 GraphQL 请求
type GraphQLRequest struct {
	Query         string                 `json:"query"`
	Variables     map[string]interface{} `json:"variables,omitempty"`
	OperationName string                 `json:"operationName,omitempty"`
}

func processRequestBody(ctx wrapper.HttpContext, cfg config.FederationConfig, body []byte) types.Action {
	//TODO:post method , content-type:graphql/json , path is /graphql
	var gqlRequest graphql.Request
	var err error
	if err = json.Unmarshal(body, &gqlRequest); err != nil {
		_ = proxywasm.SendHttpResponse(http.StatusBadRequest, nil, []byte(err.Error()), -1)
		return types.ActionContinue
	}
	var opts []engine.ExecutionOptions

	tracingOpts := resolve.TraceOptions{
		Enable:                                 true,
		ExcludePlannerStats:                    false,
		ExcludeRawInputData:                    false,
		ExcludeInput:                           false,
		ExcludeOutput:                          false,
		ExcludeLoadStats:                       false,
		EnablePredictableDebugTimings:          false,
		IncludeTraceOutputInResponseExtensions: true,
	}

	opts = append(opts, engine.WithRequestTraceOptions(tracingOpts))


	//TODO: 实现此处执行graphql逻辑
	//buf := bytes.NewBuffer(make([]byte, 0, 4096))
	//resultWriter := graphql.NewEngineResultWriterFromBuffer(buf)
	//if err = cfg.ExecutionEngine.Execute(context.TODO(), &gqlRequest, &resultWriter, opts...); err != nil {
	//	proxywasm.LogErrorf("Failed to execute GraphQL query: %v", err)
	//	_ = proxywasm.SendHttpResponse(http.StatusInternalServerError, nil, []byte(err.Error()), -1)
	//	return types.ActionContinue
	//}

	return types.ActionContinue
}

var QueryNotFoundError = errors.New("query not found")
var ParseGraphqlError = errors.New("parse graphql error")

// processGraphQLRequest 处理 GraphQL 请求
func processGraphQLRequest(request GraphQLRequest) error {
	// 使用 wundergraph/graphql-go-tools 解析查询
	document, parseReport := astparser.ParseGraphqlDocumentString(request.Query)
	proxywasm.LogDebugf("parseReport:%+v", parseReport)
	proxywasm.LogDebugf("document:%+v", document)
	if parseReport.HasErrors() {
		proxywasm.LogErrorf("Failed to parse GraphQL query: %v", parseReport.Error())
		return ParseGraphqlError
	}

	return nil
}

// parseGraphQLRequest 解析 GraphQL 请求
func parseGraphQLRequest(body []byte) (GraphQLRequest, error) {
	if len(body) == 0 {
		return GraphQLRequest{}, fmt.Errorf("empty request body")
	}

	var request GraphQLRequest
	if err := json.Unmarshal(body, &request); err != nil {
		return GraphQLRequest{}, err
	}
	// 验证请求
	if strings.TrimSpace(request.Query) == "" {
		return GraphQLRequest{}, QueryNotFoundError
	}
	return request, nil
}
