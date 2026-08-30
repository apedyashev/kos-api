package repository

import (
	"context"

	"kos-api/domains/space/persistence/entities"
	"kos-api/domains/space/service/domain"
)

// SpaceRepository defines persistence operations needed by the space service.
type SpaceRepository interface {
	List(ctx context.Context, input domain.ListSpacesInput) ([]entities.SpaceEntity, error)
}
