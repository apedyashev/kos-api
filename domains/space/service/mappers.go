package service

import (
	"kos-api/domains/space/persistence/entities"
	"kos-api/domains/space/service/domain"
)

// MapEntityToDomain converts persistence data into a domain model.
func MapEntityToDomain(entity entities.SpaceEntity) domain.SpaceDomain {
	// TODO: map entity to domain.
	return domain.SpaceDomain{
		ID:   entity.ID,
		Name: entity.Name,
		Slug: entity.Slug,
	}
}

// MapDomainToEntity converts a domain model into persistence data.
func MapDomainToEntity(currentSpace domain.SpaceDomain) entities.SpaceEntity {
	// TODO: map domain to entity.
	return entities.SpaceEntity{}
}
