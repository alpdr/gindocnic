package gindocnic

import (
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gin-gonic/gin"
	"io"
	"os"
	"testing"
)

func TestRequestBody(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		request any
		assert  func(s *openapi3.SchemaRef) error
	}{
		{name: `The binding:"required" tag marks the corresponding field as required in the generated OpenAPI specification.`, request: struct {
			Title string `json:"title" binding:"required"`
		}{
			Title: "title",
		}, assert: func(s *openapi3.SchemaRef) error {
			required := s.Value.Required
			if len(required) != 1 {
				return fmt.Errorf("unexpected required: %v", required)
			}
			if required[0] != "title" {
				return fmt.Errorf("unexpected required: %v", required)
			}
			return nil
		}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			gin.DefaultWriter = io.Discard
			defer func() {
				gin.DefaultWriter = os.Stdout
			}()
			r := gin.Default()
			sut := NewDoc()
			spec := func(p *PathItemSpec) {
				p.AddRequest(test.request)
			}
			r.POST("/posts", sut.Operation(func(c *gin.Context) {}, spec))
			if err := sut.AssocRoutesInfo(r.Routes()); err != nil {
				t.Error(err)
			}
			yml, err := sut.MarshalYAML()
			if err != nil {
				t.Error(err)
			}
			doc, err := openapi3.NewLoader().LoadFromData(yml)
			if err != nil {
				t.Error(err)
			}
			schema := doc.Paths.Find("/posts").Post.RequestBody.Value.Content.Get("application/json").Schema
			err = test.assert(schema)
			if err != nil {
				t.Error(err)
			}
		})
	}

}
