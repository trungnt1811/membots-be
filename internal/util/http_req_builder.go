package util

import (
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/imroc/req/v3"
)

var caller *resty.Client

func init() {
	caller = resty.New()                                    // Create new caller
	caller.SetRetryCount(3)                                 // Retry 3 times
	caller.SetTimeout(time.Duration(30 * int(time.Second))) // Set timeout 30s
}

type HttpRequestBuilder struct {
	headers map[string]string
	req     *resty.Request
}

func NewHttpRequestBuilder() *HttpRequestBuilder {
	req := caller.R()
	return &HttpRequestBuilder{
		headers: map[string]string{},
		req:     req,
	}
}

func (builder *HttpRequestBuilder) setHeaders() {
	req.SetHeaders(builder.headers)
}

func (builder *HttpRequestBuilder) SetHeaders(headers map[string]string) *HttpRequestBuilder {
	builder.headers = headers
	return builder
}

func (builder *HttpRequestBuilder) SetBody(body interface{}) *HttpRequestBuilder {
	builder.req.SetBody(body)
	return builder
}

func (builder *HttpRequestBuilder) SetParams(params map[string]string) *HttpRequestBuilder {
	builder.req.SetQueryParams(params)
	return builder
}

func (builder *HttpRequestBuilder) Build() *resty.Request {
	builder.setHeaders()
	return builder.req
}
