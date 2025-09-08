package config

import (
	"bytes"
	"io"
	"net/http"
	"strconv"

	"github.com/higress-group/proxy-wasm-go-sdk/proxywasm"
	"github.com/higress-group/wasm-go/pkg/wrapper"
)

type HttpClientRoundTripper struct {
}

func (h HttpClientRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	portStr := req.URL.Port()
	if len(portStr) == 0 {
		portStr = "80"
	}
	port, err := strconv.ParseInt(portStr, 10, 64)
	if err != nil {
		return nil, err
	}
	proxywasm.LogDebugf("host:%v,port:%v,method:%s,url:%s", req.Host, port, req.Method, req.URL.Path)
	client := wrapper.NewClusterClient(wrapper.FQDNCluster{
		FQDN: req.Host,
		Port: port,
	})

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	response := new(http.Response)
	response.Request = req

	//TODO:如何处理异步?
	doneRequest := make(chan interface{})
	go client.Call(req.Method, req.URL.Path, nil, body, func(statusCode int, responseHeaders http.Header, responseBody []byte) {
		proxywasm.LogDebugf("statusCode:%v,header:%+v,responseBody:%s", statusCode, responseHeaders, responseBody)
		response.Body = newResponse(responseBody)
		response.Header = responseHeaders
		response.Status = strconv.Itoa(statusCode)
		response.StatusCode = statusCode
		close(doneRequest)
	})
	<-doneRequest
	return response, nil
}

type ResponseBody struct {
	*bytes.Reader
}

func newResponse(body []byte) *ResponseBody {
	return &ResponseBody{
		Reader: bytes.NewReader(body),
	}
}

func (r *ResponseBody) Read(p []byte) (n int, err error) {
	return r.Read(p)
}

func (r *ResponseBody) Close() (_ error) {
	return nil
}
