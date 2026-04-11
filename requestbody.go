package gindocnic

import (
	"fmt"

	"github.com/swaggest/openapi-go/openapi31"
)

func setRequestBodyRequired(p PathItemSpec, pathItems map[string]openapi31.PathItem) error {
	pathItem, ok := pathItems[p.path]
	if !ok {
		return fmt.Errorf("the path item of %#v was not found", p.path)
	}
	required := true
	if p.httpMethod == "POST" {
		pathItem.Post.RequestBody.RequestBody.Required = &required
	}
	if p.httpMethod == "PUT" {
		pathItem.Put.RequestBody.RequestBody.Required = &required
	}
	return nil
}
