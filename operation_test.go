package gindocnic

import (
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gin-gonic/gin"
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
		{name: `pattern:"<regex>" tag indicates the expected regex pattern for the corresponding field in the OpenAPI document.`, request: struct {
			Title string `json:"title" pattern:"^[a-zA-Z]+$"`
		}{
			Title: "title",
		}, assert: func(s *openapi3.SchemaRef) error {
			pattern := s.Value.Properties["title"].Value.Pattern
			if pattern != "^[a-zA-Z]+$" {
				return fmt.Errorf("unexpected pattern: %v", pattern)
			}
			return nil
		},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			gin.DefaultWriter = io.Discard
			defer func() {
				gin.DefaultWriter = os.Stdout
			}()
			r := gin.Default()
			sut := MakeDoc()
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

func TestSetSecurity(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		spec   PathItemSpecFunc
		assert func(t *testing.T, secReqs *openapi3.SecurityRequirements)
	}{
		{
			name: "If no SetSecurity() is called, security should be nil",
			spec: func(o *PathItemSpec) {},
			assert: func(t *testing.T, secReqs *openapi3.SecurityRequirements) {
				if secReqs != nil {
					t.Fatal("Expected security to be nil when no SetSecurity() is called")
				}
			},
		},
		{
			name: "SetSecurity sets an operation-level security requirement",
			spec: SetSecurity("petStoreAuth", []string{"read:pets"}),
			assert: func(t *testing.T, secReqs *openapi3.SecurityRequirements) {
				if secReqs == nil {
					t.Fatal("SetSecurity(): security is nil")
				}
				if got, want := len(*secReqs), 1; got != want {
					t.Fatalf("SetSecurity(): len(security) = %d, want %d", got, want)
				}
				scopes, ok := (*secReqs)[0]["petStoreAuth"]
				if !ok {
					t.Fatal("SetSecurity(): petStoreAuth not found in operation security")
				}
				if got, want := len(scopes), 1; got != want {
					t.Fatalf("SetSecurity(): len(scopes) = %d, want %d", got, want)
				}
				if got, want := scopes[0], "read:pets"; got != want {
					t.Errorf("SetSecurity(): scopes[0] = %q, want %q", got, want)
				}
			},
		},
		{
			name: "SetSecurity with empty name sets an operation-level security requirement as security: []",
			spec: SetSecurity("", []string{}),
			assert: func(t *testing.T, secReqs *openapi3.SecurityRequirements) {
				if secReqs == nil {
					t.Fatal("SetSecurity(): security is nil")
				}
				// Since security is empty, there should be no security requirements
				if got, want := len(*secReqs), 0; got != want {
					t.Fatalf("SetSecurity(): len(security) = %d, want %d", got, want)
				}
			},
		},
		{
			name: "SetEmptySecurity works as security: []",
			spec: SetEmptySecurity(),
			assert: func(t *testing.T, secReqs *openapi3.SecurityRequirements) {
				if secReqs == nil {
					t.Fatal("SetSecurity(): security is nil")
				}
				// Since security is empty, there should be no security requirements
				if got, want := len(*secReqs), 0; got != want {
					t.Fatalf("SetSecurity(): len(security) = %d, want %d", got, want)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			sut := MakeDoc().WithOAuthClientCredentials("petStoreAuth", OAuthClientCredentials{
				TokenURL: "https://auth.example.com/oauth/token",
				Scopes:   map[string]string{"read:pets": "read pets"},
			})
			r := gin.New()
			r.GET("/pets", sut.Operation(func(*gin.Context) {}, tt.spec))
			if err := sut.AssocRoutesInfo(r.Routes()); err != nil {
				t.Fatal(err)
			}
			yml, err := sut.MarshalYAML()
			if err != nil {
				t.Fatal(err)
			}
			doc, err := openapi3.NewLoader().LoadFromData(yml)
			if err != nil {
				t.Fatal(err)
			}

			op := doc.Paths.Find("/pets").Get
			tt.assert(t, op.Security)
		})
	}
}

func TestPathItemSpecFunc(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		opt    PathItemSpecFunc
		assert func(t *testing.T, op *openapi3.Operation)
	}{
		{
			name: "OperationID sets the operation ID of the operation in the generated document",
			opt:  OperationID("createPost"),
			assert: func(t *testing.T, op *openapi3.Operation) {
				got, want := op.OperationID, "createPost"
				if got != want {
					t.Errorf("OperationID(%q) = %q, want %q", want, got, want)
				}
			},
		},
		{
			name: "OperationDescription sets the description of the operation in the generated document",
			opt:  OperationDescription("A detailed description."),
			assert: func(t *testing.T, op *openapi3.Operation) {
				got, want := op.Description, "A detailed description."
				if got != want {
					t.Errorf("OperationDescription(%q) = %q, want %q", want, got, want)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			sut := MakeDoc()
			r := gin.New()
			r.POST("/posts", sut.Operation(func(*gin.Context) {}, tt.opt))
			if err := sut.AssocRoutesInfo(r.Routes()); err != nil {
				t.Fatal(err)
			}
			yml, err := sut.MarshalYAML()
			if err != nil {
				t.Fatal(err)
			}
			doc, err := openapi3.NewLoader().LoadFromData(yml)
			if err != nil {
				t.Fatal(err)
			}

			tt.assert(t, doc.Paths.Find("/posts").Post)
		})
	}
}
