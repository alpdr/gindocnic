package gindocnic

import (
	"github.com/getkin/kin-openapi/openapi3"
	"testing"
)

func TestGenerateOpenAPI31(t *testing.T) {
	sut := NewDoc()

	bytes, err := sut.MarshalYAML()
	if err != nil {
		t.Errorf("failed to marshal %#v: %#v.", sut, err)
	}
	loader := openapi3.NewLoader()

	actual, err := loader.LoadFromData(bytes)
	if err != nil {
		t.Errorf("failed to unmarshal %#v: %#v", string(bytes), err)
	}

	expected := "3.1.0"
	if actual.OpenAPI != "3.1.0" {
		t.Errorf("the version number of the OpenAPI Specification is not %s but %#v", expected, actual.OpenAPI)
	}

}
