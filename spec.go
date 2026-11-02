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

// WithOpenAPIVersion sets the version of the OpenAPI Specification to [openapi-object.openapi].
//
// [openapi-object.openapi]: https://spec.openapis.org/oas/v3.1.0.html#openapi-object
func (d Doc) WithOpenAPIVersion(version string) Doc {
	d.reflector.Spec.Openapi = version
	return d
}

// OAuthClientCredentials holds the required fields for the [clientCredentials OAuth2 flow].
//
// [clientCredentials OAuth2 flow]: https://spec.openapis.org/oas/v3.1.0.html#oauth-flow-object
type OAuthClientCredentials struct {
	TokenURL string
	Scopes   map[string]string
}

// WithOAuthClientCredentials registers an OAuth2 clientCredentials security scheme under the
// given name in [components/securitySchemes].
//
// [components/securitySchemes]: https://spec.openapis.org/oas/v3.1.0.html#security-scheme-object
func (d Doc) WithOAuthClientCredentials(name string, flow OAuthClientCredentials) Doc {
	scheme := openapi31.SecuritySchemeOrReference{}
	scheme.SecuritySchemeEns().WithOauth2(openapi31.SecuritySchemeOauth2{
		Flows: openapi31.OauthFlows{
			ClientCredentials: &openapi31.OauthFlowsDefsClientCredentials{
				TokenURL: flow.TokenURL,
				Scopes:   flow.Scopes,
			},
		},
	})
	d.reflector.Spec.ComponentsEns().WithSecuritySchemesItem(name, scheme)
	return d
}

// OAuthAuthorizationCode holds the required fields for the [authorizationCode OAuth2 flow].
//
// [authorizationCode OAuth2 flow]: https://spec.openapis.org/oas/v3.1.0.html#oauth-flow-object
type OAuthAuthorizationCode struct {
	AuthorizationURL string
	TokenURL         string
	Scopes           map[string]string
}

// WithOAuthAuthorizationCode registers an OAuth2 authorizationCode security scheme under the
// given name in [components/securitySchemes].
//
// [components/securitySchemes]: https://spec.openapis.org/oas/v3.1.0.html#security-scheme-object
func (d Doc) WithOAuthAuthorizationCode(name string, flow OAuthAuthorizationCode) Doc {
	scheme := openapi31.SecuritySchemeOrReference{}
	scheme.SecuritySchemeEns().WithOauth2(openapi31.SecuritySchemeOauth2{
		Flows: openapi31.OauthFlows{
			AuthorizationCode: &openapi31.OauthFlowsDefsAuthorizationCode{
				AuthorizationURL: flow.AuthorizationURL,
				TokenURL:         flow.TokenURL,
				Scopes:           flow.Scopes,
			},
		},
	})
	d.reflector.Spec.ComponentsEns().WithSecuritySchemesItem(name, scheme)
	return d
}
