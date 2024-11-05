package internal

// API全体の情報(ライセンスやAPIの説明)を指定するためのAPIを提供します。
// これらのAPIは、指定しないとredocのlintが警告するので用意したであり、
// 警告を消すためだけに用意しただけなので、GoからOpen APIの定義を生成できるか調べる上では重要な機能ではないです。
// openapi31がメソッドチェーン形式のAPIなので、ひとまずメソッドチェーンの形式でAPIを実装しています。
import (
	"github.com/swaggest/openapi-go/openapi31"
)

// WithServer [server]を登録します。
// [server]: https://spec.openapis.org/oas/v3.1.0.html#server-object
func (d *Doc) WithServer(url string) *Doc {
	d.reflector.Spec.WithServers(openapi31.Server{URL: url})
	return d
}

// WithSecurity [Security Scheme Object]を登録します。
// 一部のAPIで求められる認証の形式を宣言します。Basic認証やOpen ID Connectなどを指定できます。
// Open APIの仕様上必須のフィールドらしいです。
// [Security Scheme Object]: https://spec.openapis.org/oas/v3.1.0.html#security-scheme-object
func (d *Doc) WithSecurity(securities map[string][]string) *Doc {
	d.reflector.Spec.WithSecurity(securities)
	return d
}

// WithSummary [info-object]のsummaryを登録します。
// [info-object]: https://spec.openapis.org/oas/v3.1.0.html#info-object
func (d *Doc) WithSummary(summary string) *Doc {
	d.reflector.Spec.Info.WithSummary(summary)
	return d
}

// WithSummary [info-object]のlicenseを登録します。
// [info-object]: https://spec.openapis.org/oas/v3.1.0.html#info-object
func (d *Doc) WithLicense(name, url string) *Doc {
	d.reflector.Spec.Info.WithLicense(openapi31.License{Name: name, URL: &url})
	return d
}
