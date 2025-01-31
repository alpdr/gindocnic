package gindocnic

func ResponseStatus(status int) responseOption {
	return func(r *responseOptions) {
		r.status = status
	}
}

func ResponseDescription(description string) responseOption {
	return func(r *responseOptions) {
		r.description = description
	}
}

func (o *PathItemSpec) AddResponse(body any, opts ...responseOption) {

	r := responseOptions{
		status: 200,
		body:   body,
		description: "",
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
	description string
}
