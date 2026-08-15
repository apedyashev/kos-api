package repository

import (
	"context"

	"kos-api/domains/node/persistence/entities"
)

// NodeRepository defines persistence operations needed by the node service.
type NodeRepository interface {
	List(ctx context.Context) ([]entities.NodeEntity, error)
	Create(ctx context.Context, node entities.NodeEntity) (entities.NodeEntity, error)
}
