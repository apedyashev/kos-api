package handlers

import (
	"kos-api/domains/space/api/dto"
	"kos-api/domains/space/service/domain"
)

// MapListRequestToDomain converts an API request DTO into domain input.
func MapListRequestToDomain(request dto.ListSpacesRequestDTO) domain.ListSpacesInput {
	// TODO: map request DTO fields to domain input.
	return domain.ListSpacesInput{
		OwnerID: "fe7bfc3c-2020-4d6e-b73b-ec041efb0be1", // TODO: Replace with actual mapping logic
	}
}

// MapDomainToResponse converts a domain space into an API response DTO.
func MapDomainToResponse(currentSpace domain.SpaceDomain) dto.ListSpacesResponseDTO {
	// TODO: map domain model to response DTO.
	return dto.ListSpacesResponseDTO{
		ID:   currentSpace.ID,
		Name: currentSpace.Name,
		Slug: currentSpace.Slug,
	}
}

// MapUpdateRequestToDomain converts an API update request DTO into domain input.
func MapUpdateRequestToDomain(id string, request dto.UpdateSpaceRequestDTO) domain.UpdateSpaceInput {
	// TODO: map request DTO fields to domain input.
	return domain.UpdateSpaceInput{
		ID:   id,
		Name: request.Name,
		Slug: request.Slug,
	}
}