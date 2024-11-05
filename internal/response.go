package internal

func ResponseStatus(status int) responseOption {
	return func(r *responseOptions) {
		r.status = status
	}
}

func (o *PathItemSpec) AddResponse(body any, opts ...responseOption) {

	r := responseOptions{
		status: 200,
		body:   body,
	}
	for _, opt := range opts {
		opt(&r)
	}
	o.responses = append(o.responses, r)

}

type responseOption func(r *responseOptions)

type responseOptions struct {
	body   any
	status int
}
