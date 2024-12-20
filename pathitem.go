package gindocnic

import (
	"reflect"

	og "github.com/swaggest/openapi-go"
	"github.com/swaggest/openapi-go/openapi31"
)

// PathItemSpec represents the fields of a Path Item Object.
// [path-item-object]: https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.1.0.md#path-item-object
type PathItemSpec struct {
	httpMethod string
	path       string
	summary    string
	requests   []requestOptions
	responses  []responseOptions
	id         string
}

func (o *PathItemSpec) SetSummary(s string) {
	o.summary = s
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
	openAPIPath := makeGinToOpenAPIPath(pathItemSpec.path)
	oc, err := reflector.NewOperationContext(pathItemSpec.httpMethod, openAPIPath)
	if err != nil {
		return err
	}

	// サマリーとIDを設定しないとredocの警告が出ます。
	// ? Open APIの仕様で必須
	oc.SetSummary(pathItemSpec.summary)
	oc.SetID(pathItemSpec.id)

	starParams := findStarParams(openAPIPath)
	containsRequestBody := false
	hook := func(tag reflect.StructTag) {
		if _, ok := tag.Lookup("json"); ok {
			containsRequestBody = true
		}
	}
	for _, req := range pathItemSpec.requests {
		convertedIn, err := req.convertStruct(starParams, &hook)
		if err != nil {
			return err
		}
		oc.AddReqStructure(convertedIn, func(cu *og.ContentUnit) {
			if req.contentType != "" {
				cu.ContentType = req.contentType
			}
		})
	}

	for _, resp := range pathItemSpec.responses {
		convertedResp, err := convertStruct(resp.body, starParams, nil)
		if err != nil {
			return err
		}
		oc.AddRespStructure(convertedResp, og.WithHTTPStatus(resp.status))
	}

	if err := reflector.AddOperation(oc); err != nil {
		return nil
	}
	if containsRequestBody {
		setRequestBodyRequired(pathItemSpec, reflector.Spec.Paths.MapOfPathItemValues)
	}

	return nil

}
