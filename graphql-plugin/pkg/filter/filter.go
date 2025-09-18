package filter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/envoyproxy/envoy/contrib/golang/common/go/api"
	"github.com/wundergraph/graphql-go-tools/execution/engine"
	"github.com/wundergraph/graphql-go-tools/execution/graphql"
	"github.com/wundergraph/graphql-go-tools/v2/pkg/engine/resolve"
)

var filterInstance *GraphqlFederationFilter

type GraphqlFederationFilter struct {
	api.PassThroughStreamFilter

	callbacks api.FilterCallbackHandler
	config    *graphqlFederationConfig
	bodys     []byte
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
	//path := header.Path()
	//method := header.Method()
	//api.LogDebugf("[graphql-federation]request.path:%s,method:%s", path, method)
	//if method != "POST" {
	//	api.LogDebugf("[graphql-federation]request.method:%s is not POST, skip graphql filter", method)
	//	return api.Continue
	//}
	//
	//if f.path == "/localreply_by_config" {
	//	return f.sendLocalReplyInternal()
	//}
	api.LogDebugf("[graphql-federation] 进入 decode header...")
	return api.Continue
	/*
		// If the code is time-consuming, to avoid blocking the Envoy,
		// we need to run the code in a background goroutine
		// and suspend & resume the GraphqlFederationFilter
		go func() {
			defer f.callbacks.DecoderFilterCallbacks().RecoverPanic()
			// do time-consuming jobs

			// resume the GraphqlFederationFilter
			f.callbacks.DecoderFilterCallbacks().Continue(status)
		}()

		// suspend the GraphqlFederationFilter
		return api.Running
	*/
}

// DecodeData might be called multiple times during handling the request body.
// The endStream is true when handling the last piece of the body.
func (f *GraphqlFederationFilter) DecodeData(buffer api.BufferInstance, endStream bool) api.StatusType {
	api.LogDebugf("[graphql-federation] 进入 decode data,endStream:%v", endStream)
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
	resultWriter := graphql.NewEngineResultWriterFromBuffer(buf)
	if err = executionEngine.Execute(context.TODO(), &gqlRequest, &resultWriter, opts...); err != nil {
		api.LogErrorf("[graphql-federation]graphql execute error:%v", err)
		return f.sendError(err)
	} else {
		return f.sendLocalResponse(resultWriter.Bytes())
	}
	// support suspending & resuming the GraphqlFederationFilter in a background goroutine
	return api.Continue
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
	f.callbacks.DecoderFilterCallbacks().SendLocalReply(501, err.Error(), nil, 0, "")
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
	}
	return filterInstance
}
