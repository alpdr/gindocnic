package gindocnic

// API全体の情報(ライセンスやAPIの説明)を指定するためのAPIを提供します。
// これらのAPIは、指定しないとredocのlintが警告するので用意したであり、
// 警告を消すためだけに用意しただけなので、GoからOpen APIの定義を生成できるか調べる上では重要な機能ではないです。
// openapi31がメソッドチェーン形式のAPIなので、ひとまずメソッドチェーンの形式でAPIを実装しています。
import (
	"github.com/swaggest/openapi-go/openapi31"
)

// Server represents [server].
//
// [server]: https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.1.0.md#server-object
type Server struct {
	URL string
}

// WithServer sets a [server].
//
// [server]: https://spec.openapis.org/oas/v3.1.0.html#server-object
func (d *Doc) WithServer(server Server) *Doc {
	d.reflector.Spec.WithServers(server.swaggestServer())
	return d
}

// WithNoneSecurities includes an empty security requirement ({}) in [Security Scheme Object].
//
// [Security Scheme Object]: https://spec.openapis.org/oas/v3.1.0.html#server-object
func (d *Doc) WithNoneSecurities() *Doc {
	d.reflector.Spec.WithSecurity(make(map[string][]string))
	return d
}

func (s Server) swaggestServer() openapi31.Server {
	return openapi31.Server{URL: s.URL}
}
