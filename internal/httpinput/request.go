package httpinput

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"

	"kos-api/internal/logger"
)

type fieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type validationResponse struct {
	Errors []fieldError `json:"errors"`
}

func BindAndValidate(c echo.Context, request interface{}, debugErrors bool) (bool, error) {
	// Binding can fail for malformed JSON or values with incompatible types.
	if err := c.Bind(request); err != nil {
		return false, c.String(http.StatusBadRequest, "invalid request")
	}
	// Validation checks the bound DTO against its validation tags.
	if err := c.Validate(request); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			response := validationResponse{Errors: make([]fieldError, 0, len(validationErrors))}
			for _, currentError := range validationErrors {
				message := "is required"
				if currentError.Tag() != "required" {
					rule := currentError.Tag()
					if currentError.Param() != "" {
						rule += "=" + currentError.Param()
					}
					message = fmt.Sprintf("failed validation rule %q", rule)
				}
				response.Errors = append(response.Errors, fieldError{Field: currentError.Field(), Message: message})
			}
			return false, c.JSON(http.StatusUnprocessableEntity, response)
		}
		logger.Error("validate request failed", err)
		message := "failed to validate request"
		if debugErrors {
			message = fmt.Sprintf("%s: %v", message, err)
		}
		return false, c.String(http.StatusInternalServerError, message)
	}
	return true, nil
}
