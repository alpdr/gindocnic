package internal

// go-playground/validatorのタグのついた構造体からswaggest/jsonschema-goのタグのついた構造体をつくるAPIです。
import (
	"fmt"
	"reflect"
)

// go-playground/validatorのタグのついた構造体からswaggest/jsonschema-goのタグのついた構造体のゼロ値を生成します。
// ignoreParamsに指定されたuriのパラメタは無視されます。
func convertStruct(s any, ignoreParams map[string]bool) (any, error) {
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		return convertStruct(v.Elem().Interface(), ignoreParams)
	}
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("the kind of %#v was not struct", s)
	}
	res, err := makeStruct(s, ignoreParams)
	if err != nil {
		return nil, err
	}
	return reflect.New(res).Elem().Interface(), nil
}

func makeStruct(s any, ignoreParams map[string]bool) (reflect.Type, error) {
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		return makeStruct(v.Elem().Interface(), ignoreParams)
	}
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("the kind of %#v was not struct", s)
	}
	n := v.NumField()
	fields := make([]reflect.StructField, 0, n)
	for i := range n {
		fs := v.Type().Field(i)
		if f, err := makeStructField(v.Field(i), fs, ignoreParams); err != nil {
			return nil, err
		} else {
			fields = append(fields, f)
			//埋め込み型の場合、flattenにします。
			if fs.Type.Kind() == reflect.Struct && fs.Anonymous {
				for j := range f.Type.NumField() {
					fields = append(fields, fs.Type.Field(j))
				}
			}

		}
	}
	return reflect.StructOf(fields), nil

}

func makeStructField(v reflect.Value, sf reflect.StructField, ignoreParams map[string]bool) (reflect.StructField, error) {
	if v.Kind() == reflect.Ptr {
		return makeStructField(v.Elem(), sf, ignoreParams)
	}

	if v.Kind() == reflect.Struct {
		if t, err := makeStruct(v.Interface(), ignoreParams); err != nil {
			return reflect.StructField{}, err
		} else {
			return reflect.StructField{
				Name: sf.Name,
				Type: t,
				Tag:  sf.Tag,
			}, nil
		}
	}
	return reflect.StructField{
		Name: sf.Name,
		Type: sf.Type,
		Tag:  makeOpenAPITag(sf.Tag, ignoreParams),
	}, nil
}
