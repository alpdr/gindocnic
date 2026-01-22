package gindocnic

import (
	"fmt"
	"reflect"
	"testing"
)

func TestGinStructToJsonSchemaGo(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		input    any
		expected any
	}{
		{
			name: "Convert uri tag to path tag",
			input: struct {
				ID string `uri:"id"`
			}{},
			expected: struct {
				ID string `path:"id"`
			}{},
		},
		{
			name: "Pattern is supported",
			input: struct {
				Message string `json:"message" binding:"required" pattern:"^[a-z]{4}$"`
			}{},
			expected: struct {
				Message string `json:"message" pattern:"^[a-z]{4}$" required:"true"`
			}{},
		},
		{
			name: "Pattern is supported",
			input: struct {
				Message string `json:"message" binding:"required" pattern:"^[a-z]{4}$"`
			}{},
			expected: struct {
				Message string `json:"message" pattern:"^[a-z]{4}$" required:"true"`
			}{},
		},
		{
			name: "Convert oneof to enum",
			input: struct {
				Message string `json:"message" binding:"required,oneof=active inactive pending"`
			}{},
			expected: struct {
				Message string `json:"message" required:"true" enum:"[\"active\",\"inactive\",\"pending\"]"`
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

	expectedVal := reflect.ValueOf(expected)
	actualVal := reflect.ValueOf(actual)

	expectedNumFields := expectedVal.NumField()
	actualNumFields := actualVal.NumField()
	if expectedNumFields != actualNumFields {
		return fmt.Errorf("field count mismatch: expected %d, got %d", expectedNumFields, actualNumFields)
	}

	for i := range expectedNumFields {
		expectedFieldVal := expectedVal.Field(i)
		actualFieldVal := actualVal.Field(i)
		if expectedFieldVal.Kind() != actualFieldVal.Kind() {
			return fmt.Errorf("field kind mismatch at index %d: expected %s, got %s", i, expectedFieldVal.Kind(), actualFieldVal.Kind())
		}
		expectedField := expectedVal.Type().Field(i)
		actualField := actualVal.Type().Field(i)
		if expectedField.Name != actualField.Name {
			return fmt.Errorf("field name mismatch at index %d: expected %s, got %s", i, expectedField.Name, actualField.Name)
		}
		if string(expectedField.Tag) != string(actualField.Tag) {
			a, ok := expectedField.Tag.Lookup("enum")
			if ok {
				fmt.Printf("actual enum tag: %s\n", a)
			}
			return fmt.Errorf("field tag mismatch at index %d: expected %s, got %s", i, expectedField.Tag, actualField.Tag)
		}

		if expectedFieldVal.Kind() == reflect.Struct {
			if err := checkStructEqual(t, expectedFieldVal.Interface(), actualFieldVal.Interface()); err != nil {
				return err
			}
		}
	}
	return nil

}
