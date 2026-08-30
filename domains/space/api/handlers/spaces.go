package handlers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"kos-api/domains/space/api/dto"
	"kos-api/domains/space/service"
	"kos-api/internal/logger"
)

// SpacesHandler handles the space collection HTTP endpoint.
type SpacesHandler struct {
	spaceService service.Service
	debugErrors  bool
}

// NewSpacesHandler creates a space handler backed by the space service contract.
func NewSpacesHandler(spaceService service.Service, debugErrors bool) *SpacesHandler {
	return &SpacesHandler{spaceService: spaceService, debugErrors: debugErrors}
}

// ListSpaces handles GET /api/spaces.
func (h *SpacesHandler) ListSpaces(c echo.Context) error {
	var request dto.ListSpacesRequestDTO
	if err := c.Bind(&request); err != nil {
		return c.String(http.StatusBadRequest, "invalid request")
	}

	spaces, err := h.spaceService.ListSpaces(c.Request().Context(), MapListRequestToDomain(request))
	if err != nil {
		logger.Error("list spaces failed", err)
		return h.internalError(c, "failed to list spaces", err)
	}

	responses := make([]dto.ListSpacesResponseDTO, 0, len(spaces))
	for _, currentSpace := range spaces {
		responses = append(responses, MapDomainToResponse(currentSpace))
	}
	return c.JSON(http.StatusOK, responses)
}

func (h *SpacesHandler) internalError(c echo.Context, message string, err error) error {
	if h.debugErrors {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("%s: %v", message, err))
	}
	return c.String(http.StatusInternalServerError, message)
}
