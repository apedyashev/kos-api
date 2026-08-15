package node

import (
	"fmt"
	"net/http"

	"kos-api/domains/node/api/handlers"
	"kos-api/domains/node/persistence/repository/postgres"
	"kos-api/domains/node/service"
	"kos-api/internal/config"
)

// Module owns node-domain composition and route registration.
type Module struct {
	handler http.Handler
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
func (m *Module) RegisterRoutes(mux *http.ServeMux) {
	mux.Handle("/api/nodes", m.handler)
}

// Close releases the node domain resources.
func (m *Module) Close() {
	m.close()
}
