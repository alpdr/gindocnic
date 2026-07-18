package gindocnic

import (
	"fmt"

	og "github.com/swaggest/openapi-go"
	"github.com/swaggest/openapi-go/openapi31"
)

// PathItemSpec represents the fields of a Path Item Object.
// [path-item-object]: https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.1.0.md#path-item-object
type PathItemSpec struct {
	httpMethod  string
	path        string
	summary     string
	description string
	requests    []requestOptions
	responses   []responseOptions
	id          string
}

func (o *PathItemSpec) SetSummary(s string) {
	o.summary = s
}

func (o *PathItemSpec) SetDescription(s string) {
	o.description = s
}

func (o *PathItemSpec) SetMethod(method string) {
	o.httpMethod = method
}

func (o *PathItemSpec) SetPath(path string) {
	o.path = path
}

func (o *PathItemSpec) SetId(id string) {
	o.id = id
}

func (o *PathItemSpec) setMethodIfUndefined(httpMethod string) {
	if o.httpMethod == "" {
		o.httpMethod = httpMethod
	}
}
func (o *PathItemSpec) setPathIfUndefined(path string) {
	if o.path == "" {
		o.path = path
	}
}

func (o *PathItemSpec) setIdIfUndefined(id string) {
	if o.id == "" {
		o.id = id
	}
}

// PathItemSpecFunc
type PathItemSpecFunc func(o *PathItemSpec)

// OperationSummary
func OperationSummary(summary string) PathItemSpecFunc {
	return func(o *PathItemSpec) {
		o.summary = summary
	}
}

// OperationDescription
func OperationDescription(description string) PathItemSpecFunc {
	return func(o *PathItemSpec) {
		o.description = description
	}
}

// OperationMethod
func OperationMethod(method string) PathItemSpecFunc {
	return func(o *PathItemSpec) {
		o.httpMethod = method
	}
}

// PathItemSpecPath
func PathItemSpecPath(path string) PathItemSpecFunc {
	return func(o *PathItemSpec) {
		o.path = path
	}
}

// addPathItem adds a Path Item Object to the OpenAPI document.
func addPathItem(reflector *openapi31.Reflector, pathItemSpec PathItemSpec) error {
	starParams := findStarParams(pathItemSpec.path)
	if len(starParams) > 0 {
		return fmt.Errorf("path parameters with '*' are not supported: %v", starParams)
	}
	openAPIPath := makeGinToOpenAPIPath(pathItemSpec.path)

	oc, err := reflector.NewOperationContext(pathItemSpec.httpMethod, openAPIPath)
	if err != nil {
		return err
	}

	oc.SetSummary(pathItemSpec.summary)
	oc.SetID(pathItemSpec.id)
	if pathItemSpec.description != "" {
		oc.SetDescription(pathItemSpec.description)
	}

	for _, req := range pathItemSpec.requests {
		if err != nil {
			return err
		}
		mapping, err := makePathFieldMapping(req.in)
		if err != nil {
			return err
		}
		oc.AddReqStructure(req.in, func(cu *og.ContentUnit) {
			if req.contentType != "" {
				cu.ContentType = req.contentType
			}
			cu.SetFieldMapping(og.InPath, mapping)
		})
	}

	for _, resp := range pathItemSpec.responses {

		options := make([]og.ContentOption, 0)
		options = append(options, og.WithHTTPStatus(resp.status))
		if resp.description != "" {
			options = append(options, withDescription(resp.description))
		}

		oc.AddRespStructure(resp.body, options...)
	}

	if err := reflector.AddOperation(oc); err != nil {
		return err
	}
	return nil

}

func withDescription(description string) func(cu *og.ContentUnit) {
	return func(cu *og.ContentUnit) {
		cu.Description = description
	}
}
