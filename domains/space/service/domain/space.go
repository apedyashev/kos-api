package domain

import "errors"

var ErrSpaceNotFound = errors.New("space not found")

// SpaceDomain is the business representation of a knowledge space.
type SpaceDomain struct {
	ID   string
	Name string
	Slug string
}

// ListSpacesInput contains the values needed to list spaces.
type ListSpacesInput struct {
	// OwnerID is the ID of the user who owns the spaces.
	OwnerID string
}

// UpdateSpaceInput contains the values needed to update a space.
type UpdateSpaceInput struct {
	// ID is the ID of the space to update.
	ID string
	// Name is the new name of the space.
	Name string
	// Slug is the new slug of the space.
	Slug string
}
