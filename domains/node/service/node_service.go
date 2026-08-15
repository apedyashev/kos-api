package service

import (
	"context"
	"errors"
	"fmt"

	"kos-api/domains/node/persistence/entities"
	"kos-api/domains/node/service/domain"
)

// Repository defines the persistence operations required by the node service.
type Repository interface {
	List(ctx context.Context) ([]entities.NodeEntity, error)
	Create(ctx context.Context, node entities.NodeEntity) (entities.NodeEntity, error)
	GetByID(ctx context.Context, id string) (entities.NodeEntity, error)
	Update(ctx context.Context, node entities.NodeEntity) (entities.NodeEntity, error)
}

// Service defines the node use cases exposed to the API layer.
type Service interface {
	ListNodes(ctx context.Context) ([]domain.NodeDomain, error)
	CreateNode(ctx context.Context, input domain.CreateNodeInput) (domain.NodeDomain, error)
	GetNode(ctx context.Context, id string) (domain.NodeDomain, error)
	UpdateNode(ctx context.Context, id string, input domain.UpdateNodeInput) (domain.NodeDomain, error)
}

type nodeService struct {
	repository Repository
}

// New creates a node service with its persistence dependency.
func New(repository Repository) Service {
	return &nodeService{repository: repository}
}

// ListNodes returns nodes ordered by the repository's persistence policy.
func (s *nodeService) ListNodes(ctx context.Context) ([]domain.NodeDomain, error) {
	entitiesList, err := s.repository.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list nodes: %w", err)
	}

	nodes := make([]domain.NodeDomain, 0, len(entitiesList))
	for _, entity := range entitiesList {
		nodes = append(nodes, MapEntityToDomain(entity))
	}
	return nodes, nil
}

// CreateNode validates and persists a new node.
func (s *nodeService) CreateNode(ctx context.Context, input domain.CreateNodeInput) (domain.NodeDomain, error) {
	if err := input.Validate(); err != nil {
		return domain.NodeDomain{}, err
	}

	createdEntity, err := s.repository.Create(ctx, MapDomainToEntity(domain.NodeDomain{
		Title: input.Title, Content: input.Content, IsStub: input.IsStub,
	}))
	if err != nil {
		return domain.NodeDomain{}, fmt.Errorf("create node: %w", err)
	}
	return MapEntityToDomain(createdEntity), nil
}

// GetNode retrieves a single node by id.
func (s *nodeService) GetNode(ctx context.Context, id string) (domain.NodeDomain, error) {
	entity, err := s.repository.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrNodeNotFound) {
			return domain.NodeDomain{}, err
		}
		return domain.NodeDomain{}, fmt.Errorf("get node: %w", err)
	}
	return MapEntityToDomain(entity), nil
}

// UpdateNode validates and persists changes to an existing node.
func (s *nodeService) UpdateNode(ctx context.Context, id string, input domain.UpdateNodeInput) (domain.NodeDomain, error) {
	if err := input.Validate(); err != nil {
		return domain.NodeDomain{}, err
	}

	updatedEntity, err := s.repository.Update(ctx, MapDomainToEntity(domain.NodeDomain{
		ID: id, Title: input.Title, Content: input.Content, IsStub: input.IsStub,
	}))
	if err != nil {
		if errors.Is(err, domain.ErrNodeNotFound) {
			return domain.NodeDomain{}, err
		}
		return domain.NodeDomain{}, fmt.Errorf("update node: %w", err)
	}
	return MapEntityToDomain(updatedEntity), nil
}
