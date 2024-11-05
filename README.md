# Gindocnic
<!--[golang-standards/project-layout](https://github.com/golang-standards/project-layout) suggested Go Report Card, Go Reference, and release badges" -->
[![Go Report Card](https://goreportcard.com/badge/github.com/alpdr/gindocnic)](https://goreportcard.com/report/github.com/alpdr/gindocnic)
[![Go Reference](https://pkg.go.dev/badge/github.com/alpdr/gindocnic.svg)](https://pkg.go.dev/github.com/alpdr/gindocnic)
[![release](https://img.shields.io/github/v/release/alpdr/gindocnic.svg?style=flat-square)](https://github.com/alpdr/gindocnic/releases)

A library for generating OpenAPI 3.1 documentation for the Gin Web Framework.

## Usage
A basic example:
```go
doc := NewDoc()
r := gin.Default()

request := struct {
    Id int `json:"id" binding:"required"`
}{}
spec := func(p *PathItemSpec) {
    p.AddRequest(request)
}
r.POST("/pets", doc.Operation(func(c *gin.Context) {}, spec))
if err := doc.AssocRoutesInfo(r.Routes()); err != nil {
    log.Fatalf("%#v", err)
}
yml, err := doc.MarshalYAML()
if err != nil {
    log.Fatalf("%#v", err)
}
fmt.Println(string(yml))
```

Output:
```
openapi: 3.1.0
info:
  title: ""
  version: ""
paths:
  /pets:
    post:
      operationId: pets0
      requestBody:
        content:
          application/json:
            schema:
              properties:
                id:
                  type: integer
              required:
              - id
              type: object
      responses:
        "204":
          description: No Content
      summary: ""
```

Further examples are available in the [API documentation on go.dev](https://pkg.go.dev/github.com/alpdr/gindocnic).

## Contributing
- The layout of this module loosely follows [a basic layout for Go application projects](https://github.com/golang-standards/project-layout).
