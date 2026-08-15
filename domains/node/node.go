package node

import (
	"fmt"

	"github.com/labstack/echo/v4"

	"kos-api/domains/node/api/handlers"
	"kos-api/domains/node/persistence/repository/postgres"
	"kos-api/domains/node/service"
	"kos-api/internal/config"
)

// Module owns node-domain composition and route registration.
type Module struct {
	handler *handlers.NodesHandler
	close   func()
}

// New constructs the node repository, service, and HTTP handler.
func New(appConfig config.Config) (*Module, error) {
	nodeRepository, err := postgres.NewNodeRepository(appConfig.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("bootstrap node repository: %w", err)
	}

	return &Module{
		handler: handlers.NewNodesHandler(service.New(nodeRepository), appConfig.DebugErrors),
		close:   nodeRepository.Close,
	}, nil
}

// RegisterRoutes registers the node domain routes.
func (m *Module) RegisterRoutes(e *echo.Echo) {
	e.GET("/api/nodes", m.handler.ListNodes)
	e.POST("/api/nodes", m.handler.CreateNode)
	e.GET("/api/nodes/:id", m.handler.GetNode)
	e.PATCH("/api/nodes/:id", m.handler.UpdateNode)
}

// Close releases the node domain resources.
func (m *Module) Close() {
	m.close()
}
