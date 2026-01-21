package gindocnic

import (
	"fmt"
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
			res, err := convertStruct(tt.v, tt.ignore, nil)
			tt.assertion(res, err)
		})
	}
}

func TestGinStructToJsonSchemaGo(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		input    any
		expected any
	}{
		{
			name: "Make uri path",
			input: struct {
				ID string `uri:"id"`
			}{},
			expected: struct {
				ID int `path:"id"`
			}{},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			actual, err := convertStruct(testCase.input, nil, nil)
			if err != nil {
				t.Errorf("unexpected error: %#v", err)
			}
			if err := checkEqual(t, testCase.expected, actual); err != nil {
				t.Error(err)
			}

		})
	}

}
func checkEqual(t *testing.T, expected, actual any) error {
	t.Helper()

	expectedVal := reflect.ValueOf(expected)
	actualVal := reflect.ValueOf(actual)
	if expectedVal.Kind() != actualVal.Kind() {
		return fmt.Errorf("kind mismatch: expected %s, got %s", expectedVal.Kind(), actualVal.Kind())
	}

	if expectedVal.Kind() == reflect.Struct {
		return checkStructEqual(t, expected, actual)
	}

	panic("not implemented")
}

func checkStructEqual(t *testing.T, expected, actual any) error {
	t.Helper()
	expectedType := reflect.TypeOf(expected)
	actualType := reflect.TypeOf(actual)

	v := reflect.ValueOf(expected)
	v.Interface()

	numFields := expectedType.NumField()
	if numFields != actualType.NumField() {
		return fmt.Errorf("field count mismatch: expected %d, got %d", expectedType.NumField(), actualType.NumField())
	}

	for i := range numFields {
		expectedField := expectedType.Field(i)
		actualField := actualType.Field(i)

	}

}
