package gindocnic

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func Example() {
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
