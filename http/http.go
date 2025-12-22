package nghttp

import (
	"net/http"
)

type (
	HttpResponse interface {
		// return HTTP status code
		StatusCode() int

		// response body
		Response() any
	}
)

var _ interface {
	error
	HttpResponse
} = (*Response)(nil)

// error response
func NewError(code Code, statusCode int, message string) *Response {
	return &Response{
		Code:       code,
		statusCode: statusCode,
		Message:    &message,
	}
}

// success response
func NewResponse(data any, opts ...Option) *Response {
	return (&Response{
		Code:       CodeOk,
		statusCode: http.StatusOK,
		Data:       data,
	}).Update(opts...)
}

func WrapResponse(val any) HttpResponse {
	if val == nil {
		return EmptyResponse()
	}

	switch t := val.(type) {
	case HttpResponse:
		return t
	default:
		return NewErrUnknown().Update(Metadata("raw", val))
	}
}

func RawResponse(val HttpResponse) any {
	if resp, ok := val.(*Response); ok {
		if raw, exists := resp.GetMetadata("raw"); exists {
			return raw
		}
	}
	return val
}

/*
	{
	  "code": OK,
	  "message": "ok",
	}
*/
func EmptyResponse() *Response {
	return &Response{
		Code:       CodeOk,
		statusCode: http.StatusOK,
	}
}
