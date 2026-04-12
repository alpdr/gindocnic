package gindocnic

import (
	"io"
	"os"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gin-gonic/gin"
)

func TestGenerateOpenAPI31(t *testing.T) {
	// Arrange
	t.Parallel()

	sut := MakeDoc()

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

func TestSchema(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		act    func(t *testing.T, engine *gin.Engine, doc Doc)
		assert func(t *testing.T, actual openapi3.T)
	}{{
		name: "uri tag is converted to path tag",
		act: func(t *testing.T, engine *gin.Engine, doc Doc) {
			engine.GET("/pets/:id", doc.Operation(func(*gin.Context) {}, func(p *PathItemSpec) {
				p.AddRequest(struct {
					ID string `uri:"id"`
				}{})
			}))
		},
		assert: func(t *testing.T, actual openapi3.T) {
			//yml, err:= actual.MarshalYAML()
			path := actual.Paths.Map()["/pets/{id}"]
			if path == nil {
				t.Errorf("path /pets/{id} not found")
			}
			if path.Get.Parameters.GetByInAndName("path", "id") == nil {
				t.Errorf("parameter id was not bound")
			}
		},
	},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			sut := MakeDoc().
				WithServer(Server{URL: "https://github.com/alpdr/gindocnic"}).
				WithoutSecurities().
				WithSummary("example API").
				WithLicense(License{Name: "Proprietary", URL: "https://uzabase.com"})

			gin.DefaultWriter = io.Discard
			defer func() {
				gin.DefaultWriter = os.Stdout
			}()
			r := gin.Default()

			tt.act(t, r, sut)
			if err := sut.AssocRoutesInfo(r.Routes()); err != nil {
				t.Errorf("failed to incorporate routes %s", err)
			}
			yml, err := sut.MarshalYAML()
			if err != nil {
				t.Errorf("failed to generate an Open API document: %s", err)
			}

			actual, err := openapi3.NewLoader().LoadFromData(yml)
			if err != nil {
				t.Errorf("failed to unmarshal %s: %s", string(yml), err)
			}
			tt.assert(t, *actual)
		})
	}

}
