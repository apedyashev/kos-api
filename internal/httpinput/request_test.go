package httpinput

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"

	"kos-api/internal/validation"
)

func TestBindAndValidateParameterizedRule(t *testing.T) {
	e := echo.New()
	e.Validator = validation.New()
	request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"display_name":"long"}`))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response := httptest.NewRecorder()
	context := e.NewContext(request, response)
	var input struct {
		Name string `json:"display_name,omitempty" validate:"max=3"`
	}

	valid, err := BindAndValidate(context, &input, false)
	if valid || err != nil {
		t.Fatalf("valid = %t, error = %v", valid, err)
	}
	if response.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusUnprocessableEntity)
	}
	var body validationResponse
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(body.Errors) != 1 || body.Errors[0].Field != "display_name" || body.Errors[0].Message != `failed validation rule "max=3"` {
		t.Fatalf("unexpected errors: %+v", body.Errors)
	}
}
