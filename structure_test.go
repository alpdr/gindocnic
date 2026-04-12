package gindocnic

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestGinStructToJsonSchemaGo(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		input    any
		expected any
	}{
		{
			name: "Convert uri tag to path tag",
			input: struct {
				ID string `uri:"id"`
			}{},
			expected: struct {
				ID string `path:"id"`
			}{},
		},
		{
			name: "Pattern is supported",
			input: struct {
				Message string `json:"message" binding:"required" pattern:"^[a-z]{4}$"`
			}{},
			expected: struct {
				Message string `json:"message" pattern:"^[a-z]{4}$" required:"true" nullable:"false"`
			}{},
		},
		{
			name: "Pattern is supported",
			input: struct {
				Message string `json:"message" binding:"required" pattern:"^[a-z]{4}$"`
			}{},
			expected: struct {
				Message string `json:"message" pattern:"^[a-z]{4}$" required:"true" nullable:"false"`
			}{},
		},
		{
			name: "Convert oneof to enum",
			input: struct {
				Message string `json:"message" binding:"required,oneof=active inactive pending"`
			}{},
			expected: struct {
				Message string `json:"message" required:"true" nullable:"false" enum:"[\"active\",\"inactive\",\"pending\"]"`
			}{},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			actual, err := convertStruct(testCase.input, nil, nil)
			if err != nil {
				t.Errorf("unexpected error: %#v", err)
			}
			if err := checkEqual(t, testCase.expected, actual); err != nil {
				t.Error(err)
			}

		})
	}

}
func checkEqual(t *testing.T, expected, actual any) error {
	t.Helper()

	expectedVal := reflect.ValueOf(expected)
	actualVal := reflect.ValueOf(actual)
	if expectedVal.Kind() != actualVal.Kind() {
		return fmt.Errorf("kind mismatch: expected %s, got %s", expectedVal.Kind(), actualVal.Kind())
	}

	if expectedVal.Kind() == reflect.Struct {
		return checkStructEqual(t, expected, actual)
	}

	panic("not implemented")
}

func checkStructEqual(t *testing.T, expected, actual any) error {
	t.Helper()

	expectedVal := reflect.ValueOf(expected)
	actualVal := reflect.ValueOf(actual)

	expectedNumFields := expectedVal.NumField()
	actualNumFields := actualVal.NumField()
	if expectedNumFields != actualNumFields {
		return fmt.Errorf("field count mismatch: expected %d, got %d", expectedNumFields, actualNumFields)
	}

	for i := range expectedNumFields {
		expectedFieldVal := expectedVal.Field(i)
		actualFieldVal := actualVal.Field(i)
		if expectedFieldVal.Kind() != actualFieldVal.Kind() {
			return fmt.Errorf("field kind mismatch at index %d: expected %s, got %s", i, expectedFieldVal.Kind(), actualFieldVal.Kind())
		}
		expectedField := expectedVal.Type().Field(i)
		actualField := actualVal.Type().Field(i)
		if expectedField.Name != actualField.Name {
			return fmt.Errorf("field name mismatch at index %d: expected %s, got %s", i, expectedField.Name, actualField.Name)
		}
		if string(expectedField.Tag) != string(actualField.Tag) {
			a, ok := expectedField.Tag.Lookup("enum")
			if ok {
				fmt.Printf("actual enum tag: %s\n", a)
			}
			return fmt.Errorf("field tag mismatch at index %d: expected %s, got %s", i, expectedField.Tag, actualField.Tag)
		}

		if expectedFieldVal.Kind() == reflect.Struct {
			if err := checkStructEqual(t, expectedFieldVal.Interface(), actualFieldVal.Interface()); err != nil {
				return err
			}
		}
	}
	return nil

}

type AddPetRequest struct {
	//ID         int    `json:"id" binding:"required"`
	Name       string   `json:"name" binding:"required" pattern:"^[a-zA-Z]+$"`
	Sex        string   `json:"sex" binding:"oneof=male female"`
	Emails     []string `json:"emails" binding:"required"`
	CustomerID string   `header:"customerId" description:"identifies a customer"`
	TrackingID string   `cookie:"trackingId"`
}

type Response struct {
	Id int `json:"id" binding:"required"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func (h *Handler) addSpec(p *PathItemSpec) {
	p.SetSummary("Add a new pet to the store")
	p.AddRequest(GetPetRequest{})
	p.AddRequest(AddPetRequest{})
	p.AddResponse(Response{}, ResponseStatus(http.StatusCreated))
	p.AddResponse(ErrorResponse{}, ResponseStatus(http.StatusBadRequest))
}
func (h *Handler) addPet(c *gin.Context) {}

type GetPetRequest struct {
	ID int `uri:"id"`
}

func (h *Handler) getSpec(p *PathItemSpec) {
	p.SetSummary("Find a pet")
	p.AddRequest(GetPetRequest{})
	p.AddResponse(Response{})
	p.AddResponse(ErrorResponse{}, ResponseStatus(http.StatusNotFound))
}
func (h *Handler) getPet(c *gin.Context) {}

type SearchPetsRequest struct {
	Name string `query:"name"`
}

func (h *Handler) searchSpec(p *PathItemSpec) {
	p.SetSummary("Search for pets")
	p.AddRequest(SearchPetsRequest{})
	p.AddResponse(Response{})
	p.AddResponse(ErrorResponse{}, ResponseStatus(http.StatusNotFound))
}
func (h *Handler) searchPets(c *gin.Context) {}

type Handler struct{}

// Example
func TestExample(t *testing.T) {
	doc := MakeDoc().
		WithServer(Server{URL: "https://github.com/alpdr/gindocnic"}).
		WithoutSecurities().
		WithSummary("example API").
		WithLicense(License{Name: "Proprietary", URL: "https://uzabase.com"})

	gin.DefaultWriter = io.Discard
	defer func() {
		gin.DefaultWriter = os.Stdout
	}()
	r := gin.Default()

	handler := Handler{}

	r.POST("/pets/{id}", doc.Operation(handler.addPet, handler.addSpec))
	//r.GET("/pets/{id}", doc.Operation(handler.getPet, handler.getSpec))
	//r.GET("/pets", doc.Operation(handler.searchPets, handler.searchSpec))
	if err := doc.AssocRoutesInfo(r.Routes()); err != nil {
		log.Fatalf("%#v", err)
	}
	fmt.Printf("%#v\n", doc.reflector.Spec.Paths)

	yml, err := doc.MarshalYAML()
	if err != nil {
		log.Fatalf("%#v", err)
	}

	fmt.Println(string(yml))
	t.Fail()
}
