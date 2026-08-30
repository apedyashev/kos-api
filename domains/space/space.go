package space

import (
	"fmt"

	"github.com/labstack/echo/v4"

	"kos-api/domains/space/api/handlers"
	"kos-api/domains/space/persistence/repository/postgres"
	"kos-api/domains/space/service"
	"kos-api/internal/config"
)

// Module owns the space-domain composition and route registration.
type Module struct {
	handler *handlers.SpacesHandler
	close   func()
}

// New constructs the space repository, service, and HTTP handler.
func New(appConfig config.Config) (*Module, error) {
	spaceRepository, err := postgres.NewSpaceRepository(appConfig.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("bootstrap space repository: %w", err)
	}

	return &Module{
		handler: handlers.NewSpacesHandler(service.New(spaceRepository), appConfig.DebugErrors),
		close:   spaceRepository.Close,
	}, nil
}

// RegisterRoutes registers the space domain routes.
func (m *Module) RegisterRoutes(e *echo.Echo) {
	// TODO: register space routes here.
	e.GET("/api/spaces", m.handler.ListSpaces)
}

// Close releases the space domain resources.
func (m *Module) Close() {
	m.close()
}
