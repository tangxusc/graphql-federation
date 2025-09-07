package config

import (
	"net/http"
)

type HttpClientRoundTripper struct {
}

func (h HttpClientRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	//使用proxywasm 发送http 请求
	// proxywasm.DispatchHttpCall(
	// 	req.URL.String(),
	// 	req.Header,
	// 	req.Body,
	// 	req.Trailer,
	// 	req.ContentLength,
	// 	func(headers http.Header, bodyData []byte, trailers http.Header) bool {})
	// return proxywasm.SendHttpRequest(req, deadline)

	return nil, nil
}
