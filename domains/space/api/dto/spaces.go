package dto

// ListSpacesRequestDTO is the transport payload for listing spaces.
type ListSpacesRequestDTO struct{}

// ListSpacesResponseDTO is the public representation of a space.
type ListSpacesResponseDTO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}
