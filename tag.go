package gindocnic

import (
	"fmt"
	"reflect"
	"strings"
)

// makeOpenAPITag go-playground/validatorのタグからswaggest/jsonschema-goのタグを生成します。
func makeOpenAPITag(sf reflect.StructTag, ignoreParams map[string]bool) reflect.StructTag {
	res := ""
	if v, ok := sf.Lookup("uri"); ok && !ignoreParams[v] {
		res += fmt.Sprintf(`path:"%s" `, v)
	}
	for _, keyword := range []string{"query", "json", "form", "header", "cookie", "example", "pattern"} {
		if v, ok := sf.Lookup(keyword); ok {
			res += fmt.Sprintf(`%s:"%s" `, keyword, v)
		}
	}

	if binding, ok := sf.Lookup("binding"); ok {
		res += makeValidation(binding)
	}

	return reflect.StructTag(res)
}

func makeValidation(binding string) string {
	if strings.Contains(binding, "required") {
		return `required:"true" `
	}
	return ""
}
