package repository

import (
	"context"

	"kos-api/domains/node/persistence/entities"
)

// NodeRepository defines persistence operations needed by the node service.
type NodeRepository interface {
	List(ctx context.Context, titleStartsWith string) ([]entities.NodeEntity, error)
	ListInSpace(ctx context.Context, spaceSlug string, titleStartsWith string) ([]entities.NodeEntity, error)
	Create(ctx context.Context, node entities.NodeEntity) (entities.NodeEntity, error)
	GetByID(ctx context.Context, id string, spaceSlug string) (entities.NodeEntity, error)
	Update(ctx context.Context, node entities.NodeEntity) (entities.NodeEntity, error)
}
