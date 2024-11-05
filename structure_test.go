package gindocnic

import (
	"reflect"
	"testing"
)

func TestConvertStruct(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name      string
		v         any
		ignore    map[string]bool
		assertion func(res any, err error)
	}{
		{name: "uriタグをpathタグに変換します", v: struct {
			ID string `uri:"id"`
		}{}, assertion: func(res any, err error) {
			if err != nil {
				t.Errorf("unexpected error: %#v", err)
			}
			typ := reflect.TypeOf(res)
			if typ.NumField() != 1 {
				t.Errorf("unexpected field number: %d", typ.NumField())
			}
			field := typ.Field(0)
			if field.Name != "ID" {
				t.Errorf("unexpected field name: %s", field.Name)
			}
			if field.Type != reflect.TypeOf("") {
				t.Errorf("unexpected type: %s", field.Type)
			}
			if tag := field.Tag.Get("path"); tag != "id" {
				t.Errorf("unexpected tag: %s", tag)
			}
		}},
	}
	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			res, err := convertStruct(tt.v, tt.ignore)
			tt.assertion(res, err)
		})
	}
}
