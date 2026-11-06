package gindocnic

import (
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
)

func TestWithOAuthClientCredentials(t *testing.T) {
	t.Parallel()

	flow := OAuthClientCredentials{
		TokenURL: "https://auth.example.com/oauth/token",
		Scopes:   map[string]string{"read:pets": "read pets"},
	}

	load := func(t *testing.T) openapi3.T {
		t.Helper()
		bytes, err := MakeDoc().WithOAuthClientCredentials("petStoreAuth", flow).MarshalYAML()
		if err != nil {
			t.Fatalf("MarshalYAML() error: %v", err)
		}
		doc, err := openapi3.NewLoader().LoadFromData(bytes)
		if err != nil {
			t.Fatalf("failed to load generated YAML: %v", err)
		}
		return *doc
	}

	t.Run("registers an OAuth2 clientCredentials scheme in components/securitySchemes", func(t *testing.T) {
		t.Parallel()
		doc := load(t)

		scheme, ok := doc.Components.SecuritySchemes["petStoreAuth"]
		if !ok {
			t.Fatalf("WithOAuthClientCredentials(%q): securitySchemes[%q] not found", "petStoreAuth", "petStoreAuth")
		}
		if got, want := scheme.Value.Type, "oauth2"; got != want {
			t.Errorf("WithOAuthClientCredentials(%q): type = %q, want %q", "petStoreAuth", got, want)
		}
		cc := scheme.Value.Flows.ClientCredentials
		if cc == nil {
			t.Fatalf("WithOAuthClientCredentials(%q): flows.clientCredentials is nil", "petStoreAuth")
		}
		if got, want := cc.TokenURL, flow.TokenURL; got != want {
			t.Errorf("WithOAuthClientCredentials(%q): tokenUrl = %q, want %q", "petStoreAuth", got, want)
		}
		if got, want := cc.Scopes["read:pets"], flow.Scopes["read:pets"]; got != want {
			t.Errorf("WithOAuthClientCredentials(%q): scopes[read:pets] = %q, want %q", "petStoreAuth", got, want)
		}
	})

}

func TestWithOAuthAuthorizationCode(t *testing.T) {
	t.Parallel()

	flow := OAuthAuthorizationCode{
		AuthorizationURL: "https://auth.example.com/oauth/authorize",
		TokenURL:         "https://auth.example.com/oauth/token",
		Scopes:           map[string]string{"read:pets": "read pets"},
	}

	load := func(t *testing.T) openapi3.T {
		t.Helper()
		bytes, err := MakeDoc().WithOAuthAuthorizationCode("petStoreAuth", flow).MarshalYAML()
		if err != nil {
			t.Fatalf("MarshalYAML() error: %v", err)
		}
		doc, err := openapi3.NewLoader().LoadFromData(bytes)
		if err != nil {
			t.Fatalf("failed to load generated YAML: %v", err)
		}
		return *doc
	}

	t.Run("registers an OAuth2 authorizationCode scheme in components/securitySchemes", func(t *testing.T) {
		t.Parallel()
		doc := load(t)

		scheme, ok := doc.Components.SecuritySchemes["petStoreAuth"]
		if !ok {
			t.Fatalf("WithOAuthAuthorizationCode(%q): securitySchemes[%q] not found", "petStoreAuth", "petStoreAuth")
		}
		if got, want := scheme.Value.Type, "oauth2"; got != want {
			t.Errorf("WithOAuthAuthorizationCode(%q): type = %q, want %q", "petStoreAuth", got, want)
		}
		ac := scheme.Value.Flows.AuthorizationCode
		if ac == nil {
			t.Fatalf("WithOAuthAuthorizationCode(%q): flows.authorizationCode is nil", "petStoreAuth")
		}
		if got, want := ac.AuthorizationURL, flow.AuthorizationURL; got != want {
			t.Errorf("WithOAuthAuthorizationCode(%q): authorizationUrl = %q, want %q", "petStoreAuth", got, want)
		}
		if got, want := ac.TokenURL, flow.TokenURL; got != want {
			t.Errorf("WithOAuthAuthorizationCode(%q): tokenUrl = %q, want %q", "petStoreAuth", got, want)
		}
		if got, want := ac.Scopes["read:pets"], flow.Scopes["read:pets"]; got != want {
			t.Errorf("WithOAuthAuthorizationCode(%q): scopes[read:pets] = %q, want %q", "petStoreAuth", got, want)
		}
	})

}

func TestWithOpenAPIVersion(t *testing.T) {
	t.Parallel()
	t.Run("the specified OpenAPI version appears in the generated document", func(t *testing.T) {
		t.Parallel()

		sut := MakeDoc().WithOpenAPIVersion("3.2")

		bytes, err := sut.MarshalYAML()
		if err != nil {
			t.Fatalf("failed to marshal: %v", err)
		}

		actual, err := openapi3.NewLoader().LoadFromData(bytes)
		if err != nil {
			t.Fatalf("failed to load generated YAML: %v", err)
		}

		if actual.OpenAPI != "3.2" {
			t.Errorf("WithOpenAPIVersion() = %q, want %q", actual.OpenAPI, "3.2")
		}
	})
}
