package service

import (
	"context"
	"errors"
	"fmt"

	"kos-api/domains/space/persistence/entities"
	"kos-api/domains/space/persistence/repository"
	"kos-api/domains/space/service/domain"
)

// Repository defines the persistence operations required by the space service.
type Repository interface {
	List(ctx context.Context, input domain.ListSpacesInput) ([]entities.SpaceEntity, error)
	Update(ctx context.Context, entity entities.SpaceEntity) (entities.SpaceEntity, error)
}

// Service defines the space use cases exposed to the API layer.
type Service interface {
	ListSpaces(ctx context.Context, input domain.ListSpacesInput) ([]domain.SpaceDomain, error)
	UpdateSpace(ctx context.Context, space domain.UpdateSpaceInput) (domain.SpaceDomain, error)
}

type spaceService struct {
	repository Repository
}

// New creates a space service with its persistence dependency.
func New(repository Repository) Service {
	return &spaceService{repository: repository}
}

// ListSpaces returns a list of spaces.
func (s *spaceService) ListSpaces(ctx context.Context, input domain.ListSpacesInput) ([]domain.SpaceDomain, error) {
	spaces, err := s.repository.List(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("list spaces: %w", err)
	}

	domainSpaces := make([]domain.SpaceDomain, 0)
	for _, space := range spaces {
		domainSpaces = append(domainSpaces, MapEntityToDomain(space))
	}

	return domainSpaces, nil
}

func (s *spaceService) UpdateSpace(ctx context.Context, space domain.UpdateSpaceInput) (domain.SpaceDomain, error) {
	updated, err := s.repository.Update(ctx, MapUpdateInputToEntity(space))
	if errors.Is(err, repository.ErrSpaceNotFound) {
		return domain.SpaceDomain{}, fmt.Errorf("update space: %w", domain.ErrSpaceNotFound)
	}
	if err != nil {
		return domain.SpaceDomain{}, fmt.Errorf("update space: %w", err)
	}
	return MapEntityToDomain(updated), nil
}
