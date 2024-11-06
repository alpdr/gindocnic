package gindocnic_test

import (
	"fmt"
	"github.com/alpdr/gindocnic"
	"github.com/gin-gonic/gin"
	"log"
)

type Request struct {
	Id int `json:"id" binding:"required"`
}

// Example
func Example() {
	doc := gindocnic.NewDoc()
	r := gin.Default()

	request := Request{}
	spec := func(p *gindocnic.PathItemSpec) {
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
	// Output:
	// openapi: 3.1.0
	// info:
	//   title: ""
	//   version: ""
	// paths:
	//   /pets:
	//     post:
	//       operationId: pets0
	//       requestBody:
	//         content:
	//           application/json:
	//             schema:
	//               properties:
	//                 id:
	//                   type: integer
	//               required:
	//               - id
	//               type: object
	//       responses:
	//         "204":
	//           description: No Content
	//       summary: ""
}
