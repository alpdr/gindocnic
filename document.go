package gindocnic

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/swaggest/jsonschema-go"
	"github.com/swaggest/openapi-go/openapi31"
)

// Doc represents the root of an OpenAPIv3.1 document.
type Doc struct {
	reflector *openapi31.Reflector
	// maps a pathItemSpecKey in handlerToPathItems to the corresponding PathItemSpec.
	pathItemSpecs map[pathItemSpecKey]PathItemSpec
	// Associates each handler with the path items that share the handler.
	handlerToPathItems map[string][]pathItemSpecKey
}

// MakeDoc returns [Doc].
func MakeDoc() Doc {
	reflector := openapi31.NewReflector()

	reflector.DefaultOptions = append(reflector.DefaultOptions, jsonschema.InterceptProp(func(params jsonschema.InterceptPropParams) error {
		if !params.Processed {
			return nil
		}

		fmt.Printf("aaa %v, %v, %s\n", params.Context.Path, params.Field.Name, params.Field.Tag)
		if binding, ok := params.Field.Tag.Lookup("binding"); ok {
			elements := strings.Split(binding, ",")
			foundDive := false
			for _, element := range elements {
				element = strings.TrimSpace(element)
				if "required" == element {
					params.ParentSchema.Required = append(params.ParentSchema.Required, params.Name)
					if params.PropertySchema.Type.SimpleTypes != nil {
						continue
					}
					types := params.PropertySchema.Type.SliceOfSimpleTypeValues
					fmt.Printf("%v\n", types)
					newTypes := make([]jsonschema.SimpleType, 0)
					for _, t := range types {
						if t != jsonschema.Null {
							newTypes = append(newTypes, t)
						}
					}
					params.PropertySchema.Type.SliceOfSimpleTypeValues = newTypes
					continue
				}

				if element == "dive" {
					foundDive = true
					continue
				}
				if !foundDive && strings.HasPrefix(element, "oneof=") {
					choices := strings.TrimPrefix(element, "oneof=")
					choicesList := strings.Fields(choices)
					for _, choice := range choicesList {
						params.PropertySchema.Enum = append(params.PropertySchema.Enum, choice)
					}
				}
			}
		}

		return nil
	}))
	return Doc{
		reflector:          reflector,
		pathItemSpecs:      make(map[pathItemSpecKey]PathItemSpec),
		handlerToPathItems: make(map[string][]pathItemSpecKey),
	}
}

// AssocRoutesInfo associates HTTP paths and methods with their corresponding handlers to generate path item objects.
func (d Doc) AssocRoutesInfo(routes gin.RoutesInfo) error {
	for i, route := range routes {

		keys, ok := d.handlerToPathItems[route.Handler]
		if !ok || len(keys) == 0 {
			// Skip the handlers that do not contain Open API spec.
			continue
		}
		var key pathItemSpecKey

		if len(keys) > 1 {
			key = makeKey(route)
		} else {
			key = keys[0]
		}

		pathItemSpec, ok := d.pathItemSpecs[key]
		if !ok {
			return fmt.Errorf("the operation options for %#v was not found", route)
		}

		pathItemSpec.setMethodIfUndefined(route.Method)
		pathItemSpec.setPathIfUndefined(route.Path)
		pathItemSpec.setIdIfUndefined(filterNonAlphaNumeric(pathItemSpec.path) + fmt.Sprintf("%d", i))

		if err := addPathItem(d.reflector, pathItemSpec); err != nil {
			return err
		}
	}
	return nil
}

// MarshalYAML returns the YAML encoding of [Doc].
func (d Doc) MarshalYAML() ([]byte, error) {
	return d.reflector.Spec.MarshalYAML()
}
