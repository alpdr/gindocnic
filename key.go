package gindocnic

import "github.com/gin-gonic/gin"

// pathItemSpecKey ginのRouteInfoと同型です。
// DocのもつハンドラとRouteInfoがもつハンドラと同一かどうか判定するためのキーです。
type pathItemSpecKey struct {
	method  string
	path    string
	handler string
}

// makeKey ginのRouteInfoからoperationKeyを生成します。
func makeKey(r gin.RouteInfo) pathItemSpecKey {
	return pathItemSpecKey{
		method:  r.Method,
		path:    r.Path,
		handler: r.Handler,
	}
}
