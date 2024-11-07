package gindocnic

import (
	"github.com/swaggest/openapi-go/openapi31"
)

// WithSummary sets a short summary of the API to [info-object.summary].
//
// [info-object.summary]: https://spec.openapis.org/oas/v3.1.0.html#info-object
func (d *Doc) WithSummary(summary string) *Doc {
	d.reflector.Spec.Info.WithSummary(summary)
	return d
}

// License represents [license-object].
//
// [license-object]: https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.1.0.md#license-object
type License struct {
	Name string
	URL  string
}

// WithLicense declares the license information for the exposed API.
// Accepts a [license-object].
//
// [license-object]: https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.1.0.md#license-object
func (d *Doc) WithLicense(l License) *Doc {
	d.reflector.Spec.Info.WithLicense(l.swaggestLicense())
	return d
}

func (l License) swaggestLicense() openapi31.License {
	return openapi31.License{Name: l.Name, URL: &l.URL}
}
