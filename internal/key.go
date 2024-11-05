package internal

import "github.com/gin-gonic/gin"

// operationKey ginのRouteInfoと同型です。
// DocのもつハンドラとRouteInfoがもつハンドラと同一かどうか判定するためのキーです。
type operationKey struct {
	method  string
	path    string
	handler string
}

// makeKey ginのRouteInfoからoperationKeyを生成します。
func makeKey(r gin.RouteInfo) operationKey {
	return operationKey{
		method:  r.Method,
		path:    r.Path,
		handler: r.Handler,
	}
}
