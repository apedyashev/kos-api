package service

import (
	"kos-api/domains/node/persistence/entities"
	"kos-api/domains/node/service/domain"
)

// MapEntityToDomain converts persistence data into a domain model.
func MapEntityToDomain(entity entities.NodeEntity) domain.NodeDomain {
	return domain.NodeDomain{
		ID: entity.ID, Title: entity.Title, Content: entity.Content, IsStub: entity.IsStub,
		CreatedAt: entity.CreatedAt, UpdatedAt: entity.UpdatedAt,
	}
}

// MapDomainToEntity converts a domain model into persistence data.
func MapDomainToEntity(currentNode domain.NodeDomain) entities.NodeEntity {
	return entities.NodeEntity{
		ID: currentNode.ID, Title: currentNode.Title, Content: currentNode.Content, IsStub: currentNode.IsStub,
		CreatedAt: currentNode.CreatedAt, UpdatedAt: currentNode.UpdatedAt,
	}
}
