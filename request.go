package gindocnic

// AddRequest configures operation request schema.
func (o *PathItemSpec) AddRequest(body any, opts ...requestOption) {
	r := requestOptions{
		in: body,
	}
	for _, opt := range opts {
		opt(&r)
	}
	o.requests = append(o.requests, r)
}

// RequestContentType リクエストのContent-Typeを設定します。
// AddRequestの可変長引数に渡してください。
func RequestContentType(contentType string) func(ro *requestOptions) {
	return func(ro *requestOptions) {
		ro.contentType = contentType
	}
}

type requestOptions struct {
	in          any
	contentType string
}

type requestOption func(r *requestOptions)

func (r requestOptions) convertStruct(starParams map[string]bool) (any, error) {
	return convertStruct(r.in, starParams)
}
