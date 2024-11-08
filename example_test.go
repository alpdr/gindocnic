package gindocnic_test

import (
	"fmt"
	"github.com/alpdr/gindocnic"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"os"
)

type AddPetRequest struct {
	ID         int    `json:"id" binding:"required"`
	CustomerID string `header:"customerId"`
	TrackingID string `cookie:"trackingId"`
}

type Response struct {
	Id int `json:"id"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func (h *Handler) addSpec(p *gindocnic.PathItemSpec) {
	p.SetSummary("Add a new pet to the store")
	p.AddRequest(AddPetRequest{})
	p.AddResponse(Response{}, gindocnic.ResponseStatus(http.StatusCreated))
	p.AddResponse(ErrorResponse{}, gindocnic.ResponseStatus(http.StatusBadRequest))
}
func (h *Handler) addPet(c *gin.Context) {}

type GetPetRequest struct {
	ID int `uri:"id"`
}

func (h *Handler) getSpec(p *gindocnic.PathItemSpec) {
	p.SetSummary("Find a pet")
	p.AddRequest(GetPetRequest{})
	p.AddResponse(Response{})
	p.AddResponse(ErrorResponse{}, gindocnic.ResponseStatus(http.StatusNotFound))
}
func (h *Handler) getPet(c *gin.Context) {}

type SearchPetsRequest struct {
	Name string `query:"name"`
}

func (h *Handler) searchSpec(p *gindocnic.PathItemSpec) {
	p.SetSummary("Search for pets")
	p.AddRequest(SearchPetsRequest{})
	p.AddResponse(Response{})
	p.AddResponse(ErrorResponse{}, gindocnic.ResponseStatus(http.StatusNotFound))
}
func (h *Handler) searchPets(c *gin.Context) {}

type Handler struct{}

// Example
func Example() {
	doc := gindocnic.NewDoc().
		WithServer(gindocnic.Server{URL: "https://github.com/alpdr/gindocnic"}).
		WithNoneSecurities().
		WithSummary("example API").
		WithLicense(gindocnic.License{Name: "Proprietary", URL: "https://uzabase.com"})

	gin.DefaultWriter = io.Discard
	defer func() {
		gin.DefaultWriter = os.Stdout
	}()
	r := gin.Default()

	handler := Handler{}

	r.POST("/pets", doc.Operation(handler.addPet, handler.addSpec))
	r.GET("/pets/{id}", doc.Operation(handler.getPet, handler.getSpec))
	r.GET("/pets", doc.Operation(handler.searchPets, handler.searchSpec))
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
	//   license:
	//     name: Proprietary
	//     url: https://uzabase.com
	//   summary: example API
	//   title: ""
	//   version: ""
	// servers:
	// - url: https://github.com/alpdr/gindocnic
	// paths:
	//   /pets:
	//     get:
	//       operationId: pets1
	//       parameters:
	//       - in: query
	//         name: name
	//         schema:
	//           type: string
	//       responses:
	//         "200":
	//           content:
	//             application/json:
	//               schema:
	//                 properties:
	//                   id:
	//                     type: integer
	//                 type: object
	//           description: OK
	//         "404":
	//           content:
	//             application/json:
	//               schema:
	//                 properties:
	//                   error:
	//                     type: string
	//                 type: object
	//           description: Not Found
	//       summary: Search for pets
	//     post:
	//       operationId: pets0
	//       parameters:
	//       - in: cookie
	//         name: trackingId
	//         schema:
	//           type: string
	//       - in: header
	//         name: customerId
	//         schema:
	//           type: string
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
	//         required: true
	//       responses:
	//         "201":
	//           content:
	//             application/json:
	//               schema:
	//                 properties:
	//                   id:
	//                     type: integer
	//                 type: object
	//           description: Created
	//         "400":
	//           content:
	//             application/json:
	//               schema:
	//                 properties:
	//                   error:
	//                     type: string
	//                 type: object
	//           description: Bad Request
	//       summary: Add a new pet to the store
	//   /pets/{id}:
	//     get:
	//       operationId: petsid2
	//       parameters:
	//       - in: path
	//         name: id
	//         required: true
	//         schema:
	//           type: integer
	//       responses:
	//         "200":
	//           content:
	//             application/json:
	//               schema:
	//                 properties:
	//                   id:
	//                     type: integer
	//                 type: object
	//           description: OK
	//         "404":
	//           content:
	//             application/json:
	//               schema:
	//                 properties:
	//                   error:
	//                     type: string
	//                 type: object
	//           description: Not Found
	//       summary: Find a pet
	// security:
	// - {}
}
