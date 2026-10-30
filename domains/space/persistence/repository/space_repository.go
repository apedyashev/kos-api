package repository

import (
	"context"
	"errors"

	"kos-api/domains/space/persistence/entities"
	"kos-api/domains/space/service/domain"
)

var ErrSpaceNotFound = errors.New("space not found")

// SpaceRepository defines persistence operations needed by the space service.
type SpaceRepository interface {
	List(ctx context.Context, input domain.ListSpacesInput) ([]entities.SpaceEntity, error)
	Update(ctx context.Context, entity entities.SpaceEntity) (entities.SpaceEntity, error)
}
