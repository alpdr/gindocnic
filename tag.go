package gindocnic

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

// makeOpenAPITag converts struct tags from go-playground/validator to tags used by swaggest/jsonschema-go.
func makeOpenAPITag(sf reflect.StructTag, ignoreParams map[string]bool) (reflect.StructTag, error) {
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
		validationStr, err := makeValidation(binding)
		if err != nil {
			return "", err
		}
		res.WriteString(validationStr)
	}

	return reflect.StructTag(strings.TrimSpace(res.String())), nil
}

func makeValidation(binding string) (string, error) {
	validations := strings.Split(binding, ",")
	var res strings.Builder
	for _, validation := range validations {
		validation = strings.TrimSpace(validation)
		if validation == "required" {
			res.WriteString(`required:"true" `)
			continue
		} else if choices, ok := strings.CutPrefix(validation, "oneof="); ok {
			re := regexp.MustCompile(`'[^']*'|\S+`)
			vals := re.FindAllString(choices, -1)
			for i := range vals {
				vals[i] = strings.Trim(vals[i], "'")
			}
			encoded, err := json.Marshal(vals)
			if err != nil {
				return "", err
			}
			fmt.Fprintf(&res, `enum:"%s" `, strings.ReplaceAll(string(encoded), "\"", "\\\""))
		}
	}
	return strings.TrimSpace(res.String()), nil
}
