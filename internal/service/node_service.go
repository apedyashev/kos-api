package service

import (
	"context"
	"fmt"

	"kos-api/internal/domain"
)

type nodeRepository interface {
	List(ctx context.Context) ([]domain.NodeDomain, error)
	Create(ctx context.Context, node domain.NodeDomain) (domain.NodeDomain, error)
}

// NodeService orchestrates node use cases.
type NodeService struct {
	repository nodeRepository
}

// NewNodeService creates a node service with its persistence dependency.
func NewNodeService(repository nodeRepository) *NodeService {
	return &NodeService{repository: repository}
}

// ListNodes returns nodes ordered by the repository's persistence policy.
func (s *NodeService) ListNodes(ctx context.Context) ([]domain.NodeDomain, error) {
	nodes, err := s.repository.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list nodes: %w", err)
	}
	return nodes, nil
}

// CreateNode validates and persists a new node.
func (s *NodeService) CreateNode(ctx context.Context, input domain.CreateNodeInput) (domain.NodeDomain, error) {
	if err := input.Validate(); err != nil {
		return domain.NodeDomain{}, err
	}

	node := domain.NodeDomain{
		Title:   input.Title,
		Content: input.Content,
		IsStub:  input.IsStub,
	}
	createdNode, err := s.repository.Create(ctx, node)
	if err != nil {
		return domain.NodeDomain{}, fmt.Errorf("create node: %w", err)
	}
	return createdNode, nil
}
