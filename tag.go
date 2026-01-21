package gindocnic

import (
	"fmt"
	"reflect"
	"strings"
)

// makeOpenAPITag go-playground/validatorのタグからswaggest/jsonschema-goのタグを生成します。
func makeOpenAPITag(sf reflect.StructTag, ignoreParams map[string]bool) reflect.StructTag {
	var res strings.Builder
	if v, ok := sf.Lookup("uri"); ok && !ignoreParams[v] {
		fmt.Fprintf(&res, `path:"%s" `, v)
	}
	for _, keyword := range []string{"query", "json", "form", "header", "cookie", "example", "pattern", "description"} {
		if v, ok := sf.Lookup(keyword); ok {
			fmt.Fprintf(&res, `%s:"%s" `, keyword, v)
		}
	}

	if binding, ok := sf.Lookup("binding"); ok {
		res.WriteString(makeValidation(binding))
	}

	return reflect.StructTag(strings.TrimSpace(res.String()))
}

func makeValidation(binding string) string {
	if strings.Contains(binding, "required") {
		return `required:"true" `
	}
	return ""
}
