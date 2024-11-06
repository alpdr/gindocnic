package gindocnic

import (
	"github.com/gin-gonic/gin"
	"reflect"
	"runtime"
)

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
