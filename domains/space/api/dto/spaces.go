package dto

// ListSpacesRequestDTO is the transport payload for listing spaces.
type ListSpacesRequestDTO struct{}

// ListSpacesResponseDTO is the public representation of a space.
type ListSpacesResponseDTO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type UpdateSpaceRequestDTO struct {
	Name string `json:"name" validate:"required,min=2,max=20"`
	Slug string `json:"slug" validate:"required,min=2,max=20"`
}
