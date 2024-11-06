package gindocnic

import (
	og "github.com/swaggest/openapi-go"
	"github.com/swaggest/openapi-go/openapi31"
)

// PathItemSpec [path-item-object]の生成に必要な情報をもちます。
// もともとはOperationOptionsという名前でしたが、[コードレビューのコメント]とoperationより上位の
// path-itemのスキーマにあるhttp methodやpathを含む構造体なので、PathItemSpecに名前を変えました。
// [コードレビューのコメント]: https://github.com/alpdr/data-platform/pull/481#discussion_r1804170012
// [path-item-object]: https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.1.0.md#path-item-object
type PathItemSpec struct {
	httpMethod          string
	path                string
	summary             string
	requestBodyRequired *bool
	requests            []requestOptions
	responses           []responseOptions
	id                  string
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

func (o *PathItemSpec) makeOperation(r openapi31.Reflector) (og.OperationContext, error) {
	openAPIPath := makeGinToOpenAPIPath(o.path)
	oc, err := r.NewOperationContext(o.httpMethod, openAPIPath)
	if err != nil {
		return nil, err
	}

	// サマリーとIDを設定しないとredocの警告が出ます。
	// ? Open APIの仕様で必須
	oc.SetSummary(o.summary)
	oc.SetID(o.id)

	starParams := findStarParams(openAPIPath)
	for _, req := range o.requests {
		convertedIn, err := convertStruct(req.in, starParams)
		if err != nil {
			return nil, err
		}
		oc.AddReqStructure(convertedIn, func(cu *og.ContentUnit) {
			if req.contentType != "" {
				cu.ContentType = req.contentType
			}
		})
	}

	for _, resp := range o.responses {
		convertedResp, err := convertStruct(resp.body, starParams)
		if err != nil {
			return nil, err
		}
		oc.AddRespStructure(convertedResp, og.WithHTTPStatus(resp.status))
	}

	return oc, nil
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
	for _, req := range pathItemSpec.requests {
		convertedIn, err := convertStruct(req.in, starParams)
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
		convertedResp, err := convertStruct(resp.body, starParams)
		if err != nil {
			return err
		}
		oc.AddRespStructure(convertedResp, og.WithHTTPStatus(resp.status))
	}

	if err := reflector.AddOperation(oc); err != nil {
		return nil
	}
	return nil

}
