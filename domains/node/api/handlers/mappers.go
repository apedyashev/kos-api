package handlers

import (
	"encoding/json"
	"strings"
	"time"

	"kos-api/domains/node/api/dto"
	"kos-api/domains/node/service/domain"
)

// MapCreateRequestToDomain converts an API request DTO into domain input.
func MapCreateRequestToDomain(request dto.CreateNodeRequestDTO) domain.CreateNodeInput {
	content := request.Content
	if content == nil {
		content = json.RawMessage(`{"type":"doc","content":[]}`)
	}
	return domain.CreateNodeInput{
		Title: strings.TrimSpace(request.Title), Content: content, IsStub: request.IsStub,
	}
}

// MapDomainToResponse converts a domain node into an API response DTO.
func MapDomainToResponse(currentNode domain.NodeDomain) dto.NodeResponseDTO {
	return dto.NodeResponseDTO{
		ID: currentNode.ID, Title: currentNode.Title, Content: currentNode.Content, IsStub: currentNode.IsStub,
		CreatedAt: currentNode.CreatedAt.UTC().Format(time.RFC3339Nano),
		UpdatedAt: currentNode.UpdatedAt.UTC().Format(time.RFC3339Nano),
	}
}
