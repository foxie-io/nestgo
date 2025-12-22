package nghttp

var internalErr = NewErrInternal()

type ResponseInfo struct {
	HttpResponse HttpResponse
	Raw          any
}

func (r *ResponseInfo) IsInternalError() bool {
	return r.HttpResponse == nil
}

func (r *ResponseInfo) StatusCode() int {
	if r.HttpResponse != nil {
		return r.HttpResponse.StatusCode()
	}

	return internalErr.StatusCode()
}

func (r *ResponseInfo) Response() any {
	if r.HttpResponse != nil {
		return r.HttpResponse.Response()
	}

	return internalErr
}
