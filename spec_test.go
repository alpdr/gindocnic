package gindocnic

import (
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
)

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
