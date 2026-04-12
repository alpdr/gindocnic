package gindocnic

// go-playground/validatorのタグのついた構造体からswaggest/jsonschema-goのタグのついた構造体をつくるAPIです。
import (
	"fmt"
	"reflect"
)

func makePathFieldMapping(s any) (map[string]string, error) {
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		return makePathFieldMapping(v.Elem().Interface())
	}

	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("the kind of %#v was not struct", s)
	}

	mapping := make(map[string]string)

	n := v.NumField()
	for i := range n {
		fs := v.Type().Field(i)
		tag := fs.Tag
		uriFound, uriOk := tag.Lookup("uri")
		_, pathOk := tag.Lookup("path")
		if uriOk && !pathOk {
			mapping[fs.Name] = uriFound
		}
	}

	return mapping, nil
}
