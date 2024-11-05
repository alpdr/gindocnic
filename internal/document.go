package internal

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/swaggest/openapi-go/openapi31"
)

// Doc OpenAPI定義書の全体の情報をもちます。
type Doc struct {
	reflector *openapi31.Reflector
	// handlerNameToOptionsにあるoperationKeyに対応するPathItemSpecをもちます。
	operationOptions map[operationKey]PathItemSpec
	// どのルーティング先がどのハンドラを共有していて、
	// ハンドラ名だけでルーティング先を特定できないかを管理しています。
	// 値の要素が1つであれば、ハンドラ名だけでルーティング先を特定できます。
	handlerNameToOptions map[string][]operationKey
}

// NewDoc Docを作り、そのポインタを返します。
func NewDoc() *Doc {
	return &Doc{
		reflector:            openapi31.NewReflector(),
		operationOptions:     make(map[operationKey]PathItemSpec),
		handlerNameToOptions: make(map[string][]operationKey),
	}
}

// AssocRoutesInfo リクエストハンドラをパスやHTTPメソッドと関連づけて、Open API仕様書に必要な情報をそろえます。
func (d *Doc) AssocRoutesInfo(routes gin.RoutesInfo) error {
	for i, route := range routes {

		keys, ok := d.handlerNameToOptions[route.Handler]
		if !ok || len(keys) == 0 {
			// Open API定義に記述しないハンドラを無視します。
			continue
		}
		var key operationKey

		if len(keys) > 1 {
			key = makeKey(route)
		} else {
			key = keys[0]
		}

		options, ok := d.operationOptions[key]
		if !ok {
			return fmt.Errorf("the operation options for %#v was not found", route)
		}

		options.setMethodIfUndefined(route.Method)
		options.setPathIfUndefined(route.Path)
		// ハンドラ名から名前をつけるとソースコードの情報が露出するのでパスに由来する名前にします。
		options.setIdIfUndefined(filterNonAlphaNumeric(options.path) + fmt.Sprintf("%d", i))

		oc, err := options.newOperation(*d.reflector)
		if err != nil {
			return err
		}

		if err := d.reflector.AddOperation(oc); err != nil {
			return err
		}
	}
	return nil
}

// MarshalYAML yaml形式のOpen API定義をバイト列で返します。
func (d *Doc) MarshalYAML() ([]byte, error) {
	return d.reflector.Spec.MarshalYAML()
}