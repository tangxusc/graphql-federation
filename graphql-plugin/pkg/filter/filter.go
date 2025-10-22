package filter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/envoyproxy/envoy/contrib/golang/common/go/api"
	"github.com/wundergraph/graphql-go-tools/execution/engine"
	"github.com/wundergraph/graphql-go-tools/execution/graphql"
	"github.com/wundergraph/graphql-go-tools/v2/pkg/engine/resolve"
)

var filterInstance *GraphqlFederationFilter

// 支持的 GraphQL Content-Type 格式
var graphqlContentTypes = []string{
	"application/json",
	"application/graphql",
	"application/graphql+json",
	"application/graphql-request+json",
	"application/vnd.graphql+json",
	"application/vnd.graphql",
}

type GraphqlFederationFilter struct {
	api.PassThroughStreamFilter

	callbacks api.FilterCallbackHandler
	config    *graphqlFederationConfig
	bodys     []byte
	isGraphQL bool // 标记是否为 GraphQL 请求
}

func (f *GraphqlFederationFilter) sendLocalReplyInternal() api.StatusType {
	response, err := http.Get("http://www.baidu.com")
	api.LogErrorf("get baidu error: %v", err)
	api.LogErrorf("get baidu response: %v", response)

	body := fmt.Sprintf("%s, path: %s\r\n", "1", "2")
	f.callbacks.DecoderFilterCallbacks().SendLocalReply(200, body, nil, 0, "")
	// Remember to return LocalReply when the request is replied locally
	return api.LocalReply
}

// Callbacks which are called in request path
// The endStream is true if the request doesn't have body
func (f *GraphqlFederationFilter) DecodeHeaders(header api.RequestHeaderMap, endStream bool) api.StatusType {
	path := header.Path()
	method := header.Method()
	contentType, _ := header.Get("Content-Type")

	api.LogDebugf("[graphql-federation] request.path:%s, method:%s, contentType:%s", path, method, contentType)

	// 检查是否为 GraphQL 请求
	f.isGraphQL = f.isGraphQLRequest(method, path, contentType)

	if !f.isGraphQL {
		api.LogDebugf("[graphql-federation] 非 GraphQL 请求，跳过处理: path=%s, method=%s", path, method)
		return api.Continue
	}

	api.LogDebugf("[graphql-federation] GraphQL 请求，开始处理")
	return api.Continue
}

// DecodeData might be called multiple times during handling the request body.
// The endStream is true when handling the last piece of the body.
func (f *GraphqlFederationFilter) DecodeData(buffer api.BufferInstance, endStream bool) api.StatusType {
	api.LogDebugf("[graphql-federation] 进入 decode data,endStream:%v", endStream)

	// 如果不是 GraphQL 请求，直接跳过处理
	if !f.isGraphQL {
		api.LogDebugf("[graphql-federation] 非 GraphQL 请求，跳过数据处理")
		return api.Continue
	}

	var gqlRequest graphql.Request
	var err error
	f.bodys = append(f.bodys, buffer.Bytes()...)
	if !endStream {
		api.LogDebugf("[graphql-federation]request.body:%s,endStream:%v", buffer.Bytes(), endStream)
		return api.StopAndBuffer
	}

	if err = json.Unmarshal(f.bodys, &gqlRequest); err != nil {
		api.LogDebugf("[graphql-federation]request.body:%s,error:%v", f.bodys, err)
		return f.sendError(err)
	}
	buf := bytes.NewBuffer(make([]byte, 0, 4096))
	var opts []engine.ExecutionOptions

	tracingOpts := resolve.TraceOptions{}
	tracingOpts.EnableAll()

	opts = append(opts, engine.WithRequestTraceOptions(tracingOpts))
	resultWriter := graphql.NewEngineResultWriterFromBuffer(buf)
	if err = executionEngine.Execute(context.TODO(), &gqlRequest, &resultWriter, opts...); err != nil {
		api.LogErrorf("[graphql-federation]graphql execute error:%v", err)
		return f.sendError(err)
	} else {
		return f.sendLocalResponse(resultWriter.Bytes())
	}
	// support suspending & resuming the GraphqlFederationFilter in a background goroutine

}

func (f *GraphqlFederationFilter) DecodeTrailers(trailers api.RequestTrailerMap) api.StatusType {
	api.LogDebugf("[graphql-federation] 进入 decode trailers...")
	// support suspending & resuming the GraphqlFederationFilter in a background goroutine
	return api.Continue
}

// Callbacks which are called in response path
// The endStream is true if the response doesn't have body
func (f *GraphqlFederationFilter) EncodeHeaders(header api.ResponseHeaderMap, endStream bool) api.StatusType {
	//if f.path == "/update_upstream_response" {
	//	header.Set("Content-Length", strconv.Itoa(len(UpdateUpstreamBody)))
	//}
	api.LogDebugf("[graphql-federation] 进入 encode headers...")
	header.Set("Rsp-Header-From-Go", "bar-test")
	// support suspending & resuming the GraphqlFederationFilter in a background goroutine
	return api.Continue
}

// EncodeData might be called multiple times during handling the response body.
// The endStream is true when handling the last piece of the body.
func (f *GraphqlFederationFilter) EncodeData(buffer api.BufferInstance, endStream bool) api.StatusType {
	//if f.path == "/update_upstream_response" {
	//	if endStream {
	//		buffer.SetString(UpdateUpstreamBody)
	//	} else {
	//		buffer.Reset()
	//	}
	//}
	api.LogDebugf("[graphql-federation] 进入 encode data,endStream:%v...", endStream)
	// support suspending & resuming the GraphqlFederationFilter in a background goroutine
	return api.Continue
}

func (f *GraphqlFederationFilter) EncodeTrailers(trailers api.ResponseTrailerMap) api.StatusType {
	api.LogDebugf("[graphql-federation] 进入 encode trailers...")
	return api.Continue
}

// OnLog is called when the HTTP stream is ended on HTTP Connection Manager GraphqlFederationFilter.
func (f *GraphqlFederationFilter) OnLog(reqHeader api.RequestHeaderMap, reqTrailer api.RequestTrailerMap, respHeader api.ResponseHeaderMap, respTrailer api.ResponseTrailerMap) {
	api.LogDebugf("[graphql-federation] 进入 on log...")
	code, _ := f.callbacks.StreamInfo().ResponseCode()
	respCode := strconv.Itoa(int(code))
	api.LogDebug(respCode)

	/*
		// It's possible to kick off a goroutine here.
		// But it's unsafe to access the f.callbacks because the FilterCallbackHandler
		// may be already released when the goroutine is scheduled.
		go func() {
			defer func() {
				if p := recover(); p != nil {
					const size = 64 << 10
					buf := make([]byte, size)
					buf = buf[:runtime.Stack(buf, false)]
					fmt.Printf("http: panic serving: %v\n%s", p, buf)
				}
			}()

			// do time-consuming jobs
		}()
	*/
}

// OnLogDownstreamStart is called when HTTP Connection Manager GraphqlFederationFilter receives a new HTTP request
// (required the corresponding access log type is enabled)
func (f *GraphqlFederationFilter) OnLogDownstreamStart(reqHeader api.RequestHeaderMap) {
	// also support kicking off a goroutine here, like OnLog.
}

// OnLogDownstreamPeriodic is called on any HTTP Connection Manager periodic log record
// (required the corresponding access log type is enabled)
func (f *GraphqlFederationFilter) OnLogDownstreamPeriodic(reqHeader api.RequestHeaderMap, reqTrailer api.RequestTrailerMap, respHeader api.ResponseHeaderMap, respTrailer api.ResponseTrailerMap) {
	// also support kicking off a goroutine here, like OnLog.
}

func (f *GraphqlFederationFilter) OnDestroy(reason api.DestroyReason) {
	api.LogDebugf("[graphql-federation] 进入 destroy...")
	// One should not access f.callbacks here because the FilterCallbackHandler
	// is released. But we can still access other Go fields in the GraphqlFederationFilter f.

	// goroutine can be used everywhere.
}

func (f *GraphqlFederationFilter) sendError(err error) api.StatusType {
	payload := map[string]interface{}{
		"data": nil,
		"errors": []map[string]interface{}{
			{
				"message":    err.Error(),
				"extensions": map[string]interface{}{"code": "INTERNAL_SERVER_ERROR"},
			},
		},
	}
	body, marshalErr := json.Marshal(payload)
	if marshalErr != nil {
		// 兜底：如果序列化失败，仍然返回简单错误信息
		f.callbacks.DecoderFilterCallbacks().SendLocalReply(200, "{\"data\":null,\"errors\":[{\"message\":\"internal error\"}]}", map[string][]string{
			"Content-Type": {"application/json; charset=utf-8"},
		}, 0, "")
		return api.LocalReply
	}

	f.callbacks.DecoderFilterCallbacks().SendLocalReply(200, string(body), map[string][]string{
		"Content-Type": {"application/json; charset=utf-8"},
	}, 0, "")
	return api.LocalReply
}

func (f *GraphqlFederationFilter) sendLocalResponse(i []byte) api.StatusType {
	f.callbacks.DecoderFilterCallbacks().SendLocalReply(200, string(i), nil, 0, "")
	return api.LocalReply
}

func NewGraphqlFederationFilter(cfg *graphqlFederationConfig, callbacks api.FilterCallbackHandler) api.StreamFilter {
	api.LogDebugf("[graphql-federation] 创建graphql federation filter...")
	filterInstance = &GraphqlFederationFilter{
		callbacks: callbacks,
		isGraphQL: false, // 初始化为 false
	}
	return filterInstance
}

// isGraphQLRequest 判断请求是否为 GraphQL 请求
func (f *GraphqlFederationFilter) isGraphQLRequest(method, path, contentType string) bool {
	// GraphQL 请求必须满足以下条件：
	// 1. 路径严格为 /graphql
	// 2. POST 请求且 Content-Type 为 GraphQL 相关格式

	if path != "/graphql" {
		return false
	}

	if method == "POST" {
		// 检查 Content-Type 是否为 GraphQL 相关格式
		contentType = strings.ToLower(strings.TrimSpace(contentType))

		// 精确匹配
		for _, ct := range graphqlContentTypes {
			if contentType == ct {
				return true
			}
		}

		// 模糊匹配（包含关键字）
		if strings.Contains(contentType, "application/json") ||
			strings.Contains(contentType, "application/graphql") {
			return true
		}
	}

	return false
}
