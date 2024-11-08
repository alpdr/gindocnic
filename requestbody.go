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
	if p.httpMethod != "POST" {
		return fmt.Errorf("only post method is supported at the moment.")
	}
	required := true
	pathItem.Post.RequestBody.RequestBody.Required = &required
	return nil
}
