package service

import (
	"context"
	"fmt"

	"kos-api/domains/node/persistence/entities"
	"kos-api/domains/node/service/domain"
)

// Repository defines the persistence operations required by the node service.
type Repository interface {
	List(ctx context.Context) ([]entities.NodeEntity, error)
	Create(ctx context.Context, node entities.NodeEntity) (entities.NodeEntity, error)
}

// Service defines the node use cases exposed to the API layer.
type Service interface {
	ListNodes(ctx context.Context) ([]domain.NodeDomain, error)
	CreateNode(ctx context.Context, input domain.CreateNodeInput) (domain.NodeDomain, error)
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
