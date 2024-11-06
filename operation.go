package gindocnic

import (
	"github.com/gin-gonic/gin"
	og "github.com/swaggest/openapi-go"
	"github.com/swaggest/openapi-go/openapi31"
	"reflect"
	"runtime"
)

// PathItemSpec [path-item-object]の生成に必要な情報をもちます。
// もともとはOperationOptionsという名前でしたが、[コードレビューのコメント]とoperationより上位の
// path-itemのスキーマにあるhttp methodやpathを含む構造体なので、PathItemSpecに名前を変えました。
// [コードレビューのコメント]: https://github.com/alpdr/data-platform/pull/481#discussion_r1804170012
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

func (o *PathItemSpec) newOperation(r openapi31.Reflector) (og.OperationContext, error) {
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

// Operation configures the path item schema for the HTTP path and method that the handler is registered to.
func (d *Doc) Operation(h gin.HandlerFunc, opts ...PathItemSpecFunc) gin.HandlerFunc {
	// r.GET('path', doc.Operation(...))のように書けるようにerrが発生する処理を
	// OperationOptions.newOperationの呼び出しに移譲します。
	ops := &PathItemSpec{}

	for _, opt := range opts {
		opt(ops)
	}
	handlerName := nameHandlerAfterGin(h)
	operationKey := pathItemSpecKey{
		handler: handlerName,
		method:  ops.httpMethod,
		path:    ops.path,
	}
	d.handlerToPathItems[handlerName] = append(d.handlerToPathItems[handlerName], operationKey)
	d.pathItemSpecs[operationKey] = *ops
	return h
}

// nameHandlerAfterGin Ginのutils.goにあるハンドラの名付け方(nameOfFunction関数)と同一の方法で、ハンドラに名前をつけます。
// *gin.Engine.RoutesInfoとDoc構造体でハンドラの名付け方をそろえることで、Doc構造体がRoutesInfoの情報をハンドラ名でひけるようにします。
func nameHandlerAfterGin(f any) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}
