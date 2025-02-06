package gindocnic

import (
	"fmt"
	"github.com/gin-gonic/gin"
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
	return Doc{
		reflector:          openapi31.NewReflector(),
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
		// ハンドラ名から名前をつけるとソースコードの情報が露出するのでパスに由来する名前にします。
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
