package domain

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
