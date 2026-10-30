package space

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"

	"kos-api/domains/space/api/handlers"
	"kos-api/domains/space/service/domain"
	"kos-api/internal/validation"
)

type fakeSpaceService struct {
	updates []domain.UpdateSpaceInput
	err     error
}

func (service *fakeSpaceService) ListSpaces(ctx context.Context, input domain.ListSpacesInput) ([]domain.SpaceDomain, error) {
	return nil, nil
}

func (service *fakeSpaceService) UpdateSpace(ctx context.Context, input domain.UpdateSpaceInput) (domain.SpaceDomain, error) {
	service.updates = append(service.updates, input)
	if service.err != nil {
		return domain.SpaceDomain{}, service.err
	}
	return domain.SpaceDomain{ID: input.ID, Name: input.Name, Slug: input.Slug}, nil
}

func TestUpdateSpaceRequestValidation(t *testing.T) {
	e := echo.New()
	e.Validator = validation.New()
	spaceService := &fakeSpaceService{}
	module := Module{handler: handlers.NewSpacesHandler(spaceService, false)}
	module.RegisterRoutes(e)

	type fieldError struct {
		Field   string `json:"field"`
		Message string `json:"message"`
	}
	tests := []struct {
		name   string
		body   string
		status int
		errors []fieldError
	}{
		{name: "missing name", body: `{"slug":"example"}`, status: http.StatusUnprocessableEntity, errors: []fieldError{{Field: "name", Message: "is required"}}},
		{name: "missing slug", body: `{"name":"Example"}`, status: http.StatusUnprocessableEntity, errors: []fieldError{{Field: "slug", Message: "is required"}}},
		{name: "missing both", body: `{}`, status: http.StatusUnprocessableEntity, errors: []fieldError{{Field: "name", Message: "is required"}, {Field: "slug", Message: "is required"}}},
		{name: "one character name", body: `{"name":"A","slug":"example"}`, status: http.StatusUnprocessableEntity, errors: []fieldError{{Field: "name", Message: `failed validation rule "min=2"`}}},
		{name: "two character name", body: `{"name":"Ab","slug":"example"}`, status: http.StatusOK},
		{name: "name exceeds maximum", body: `{"name":"abcdefghijklmnopqrstu","slug":"example"}`, status: http.StatusUnprocessableEntity, errors: []fieldError{{Field: "name", Message: `failed validation rule "max=20"`}}},
		{name: "one character slug", body: `{"name":"Example","slug":"a"}`, status: http.StatusUnprocessableEntity, errors: []fieldError{{Field: "slug", Message: `failed validation rule "min=2"`}}},
		{name: "slug exceeds maximum", body: `{"name":"Example","slug":"abcdefghijklmnopqrstu"}`, status: http.StatusUnprocessableEntity, errors: []fieldError{{Field: "slug", Message: `failed validation rule "max=20"`}}},
		{name: "minimum slug length", body: `{"name":"Example","slug":"ab"}`, status: http.StatusOK},
		{name: "maximum field lengths", body: `{"name":"abcdefghijklmnopqrst","slug":"abcdefghijklmnopqrst"}`, status: http.StatusOK},
		{name: "invalid JSON", body: `{`, status: http.StatusBadRequest},
		{name: "wrong field type", body: `{"name":42,"slug":"example"}`, status: http.StatusBadRequest},
		{name: "valid request", body: `{"name":"Example","slug":"example"}`, status: http.StatusOK},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			spaceService.updates = nil
			request := httptest.NewRequest(http.MethodPut, "/api/spaces/space-1", strings.NewReader(test.body))
			request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			response := httptest.NewRecorder()
			e.ServeHTTP(response, request)
			if response.Code != test.status {
				t.Fatalf("status = %d, want %d", response.Code, test.status)
			}
			if called := len(spaceService.updates) > 0; called != (test.status == http.StatusOK) {
				t.Fatalf("service called = %t for status %d", called, test.status)
			}
			if test.status == http.StatusOK && spaceService.updates[0].ID != "space-1" {
				t.Fatalf("space ID = %q, want space-1", spaceService.updates[0].ID)
			}
			if test.status == http.StatusOK {
				var body struct {
					ID   string `json:"id"`
					Name string `json:"name"`
					Slug string `json:"slug"`
				}
				if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
					t.Fatalf("decode response: %v", err)
				}
				if body.ID != "space-1" || body.Name != spaceService.updates[0].Name || body.Slug != spaceService.updates[0].Slug {
					t.Fatalf("unexpected response: %+v", body)
				}
			}
			if test.errors != nil {
				var body struct {
					Errors []fieldError `json:"errors"`
				}
				if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
					t.Fatalf("decode response: %v", err)
				}
				if !reflect.DeepEqual(body.Errors, test.errors) {
					t.Fatalf("errors = %+v, want %+v", body.Errors, test.errors)
				}
			}
		})
	}

	for _, debugErrors := range []bool{false, true} {
		t.Run(fmt.Sprintf("validator not configured debug=%t", debugErrors), func(t *testing.T) {
			e := echo.New()
			module := Module{handler: handlers.NewSpacesHandler(nil, debugErrors)}
			module.RegisterRoutes(e)
			request := httptest.NewRequest(http.MethodPut, "/api/spaces/space-1", strings.NewReader(`{"name":"Example","slug":"example"}`))
			request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			response := httptest.NewRecorder()
			e.ServeHTTP(response, request)
			if response.Code != http.StatusInternalServerError {
				t.Fatalf("status = %d, want %d", response.Code, http.StatusInternalServerError)
			}
			message := response.Body.String()
			if !debugErrors && message != "failed to validate request" {
				t.Fatalf("unexpected response: %q", message)
			}
			if debugErrors && !strings.HasPrefix(message, "failed to validate request: ") {
				t.Fatalf("missing debug details: %q", message)
			}
		})
	}
}

func TestUpdateSpaceServiceErrors(t *testing.T) {
	for _, test := range []struct {
		name   string
		err    error
		status int
		body   string
	}{
		{name: "not found", err: fmt.Errorf("update: %w", domain.ErrSpaceNotFound), status: http.StatusNotFound, body: "space not found"},
		{name: "unexpected error", err: fmt.Errorf("database unavailable"), status: http.StatusInternalServerError, body: "failed to update space"},
	} {
		t.Run(test.name, func(t *testing.T) {
			e := echo.New()
			e.Validator = validation.New()
			module := Module{handler: handlers.NewSpacesHandler(&fakeSpaceService{err: test.err}, false)}
			module.RegisterRoutes(e)
			request := httptest.NewRequest(http.MethodPut, "/api/spaces/space-1", strings.NewReader(`{"name":"Example","slug":"example"}`))
			request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			response := httptest.NewRecorder()
			e.ServeHTTP(response, request)
			if response.Code != test.status {
				t.Fatalf("status = %d, want %d", response.Code, test.status)
			}
			if response.Body.String() != test.body {
				t.Fatalf("body = %q, want %q", response.Body.String(), test.body)
			}
		})
	}
}
