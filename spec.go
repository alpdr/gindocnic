package gindocnic

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
func (d Doc) WithServer(server Server) Doc {
	d.reflector.Spec.WithServers(server.swaggestServer())
	return d
}

// WithoutSecurities includes an empty security requirement ({}) in [Security Scheme Object].
//
// [Security Scheme Object]: https://spec.openapis.org/oas/v3.1.0.html#server-object
func (d Doc) WithoutSecurities() Doc {
	d.reflector.Spec.WithSecurity(make(map[string][]string))
	return d
}

func (s Server) swaggestServer() openapi31.Server {
	return openapi31.Server{URL: s.URL}
}
