package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"kos-api/domains/node/api/dto"
	"kos-api/domains/node/service"
	"kos-api/domains/node/service/domain"
	"kos-api/internal/logger"
)

// NodesHandler handles the node collection HTTP endpoint.
type NodesHandler struct {
	nodeService service.Service
	debugErrors bool
}

// NewNodesHandler creates a node handler backed by the node service contract.
func NewNodesHandler(nodeService service.Service, debugErrors bool) *NodesHandler {
	return &NodesHandler{nodeService: nodeService, debugErrors: debugErrors}
}

// ListNodes handles GET /api/nodes.
func (h *NodesHandler) ListNodes(c echo.Context) error {
	nodes, err := h.nodeService.ListNodes(c.Request().Context(), c.QueryParam("search"))
	if err != nil {
		logger.Error("list nodes failed", err)
		return h.internalError(c, "failed to list nodes", err)
	}

	responses := make([]dto.NodeResponseDTO, 0, len(nodes))
	for _, currentNode := range nodes {
		responses = append(responses, MapDomainToResponse(currentNode))
	}
	return c.JSON(http.StatusOK, responses)
}

// CreateNode handles POST /api/nodes.
func (h *NodesHandler) CreateNode(c echo.Context) error {
	var request dto.CreateNodeRequestDTO
	if err := c.Bind(&request); err != nil {
		return c.String(http.StatusBadRequest, "invalid JSON body")
	}

	createdNode, err := h.nodeService.CreateNode(c.Request().Context(), MapCreateRequestToDomain(request))
	if err != nil {
		if errors.Is(err, domain.ErrInvalidNode) {
			return c.String(http.StatusBadRequest, err.Error())
		}
		logger.Error("create node failed", err)
		return h.internalError(c, "failed to create node", err)
	}

	return c.JSON(http.StatusCreated, MapDomainToResponse(createdNode))
}

// GetNode handles GET /api/space/:spaceSlug/nodes/:id.
func (h *NodesHandler) GetNode(c echo.Context) error {
	spaceSlug := c.Param("spaceSlug")
	foundNode, err := h.nodeService.GetNode(c.Request().Context(), c.Param("id"), spaceSlug)

	if err != nil {
		if errors.Is(err, domain.ErrNodeNotFound) {
			return c.String(http.StatusNotFound, "node not found")
		}
		logger.Error("get node failed", err)
		return h.internalError(c, "failed to get node", err)
	}

	return c.JSON(http.StatusOK, MapDomainToResponse(foundNode))
}

// UpdateNode handles PATCH /api/nodes/:id.
func (h *NodesHandler) UpdateNode(c echo.Context) error {
	var request dto.UpdateNodeRequestDTO
	if err := c.Bind(&request); err != nil {
		return c.String(http.StatusBadRequest, "invalid JSON body")
	}

	updatedNode, err := h.nodeService.UpdateNode(c.Request().Context(), c.Param("id"), MapUpdateRequestToDomain(request))
	if err != nil {
		if errors.Is(err, domain.ErrInvalidNode) {
			return c.String(http.StatusBadRequest, err.Error())
		}
		if errors.Is(err, domain.ErrNodeNotFound) {
			return c.String(http.StatusNotFound, "node not found")
		}
		logger.Error("update node failed", err)
		return h.internalError(c, "failed to update node", err)
	}

	return c.JSON(http.StatusOK, MapDomainToResponse(updatedNode))
}

func (h *NodesHandler) internalError(c echo.Context, message string, err error) error {
	if h.debugErrors {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("%s: %v", message, err))
	}
	return c.String(http.StatusInternalServerError, message)
}
