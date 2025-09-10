package process

import (
	"bytes"
	"io"
	"net/http"
	"strconv"
	"time"

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
		FQDN: "httpbin",
		Port: port,
	})

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	response := new(http.Response)
	response.Request = req

	client.Call(req.Method, req.URL.Path, nil, body,
		func(statusCode int, responseHeaders http.Header, responseBody []byte) {
			proxywasm.LogDebugf("statusCode:%v,header:%+v,responseBody:%s", statusCode, responseHeaders, responseBody)
			response.Body = newResponse(responseBody)
			response.Header = responseHeaders
			response.Status = strconv.Itoa(statusCode)
			response.StatusCode = statusCode
		})
		
	time.Sleep(time.Second * 5)
	proxywasm.LogDebugf("返回数据:%+v", response)
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
