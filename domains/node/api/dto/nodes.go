package dto

import "encoding/json"

// NodeResponseDTO is the public representation of a node.
type NodeResponseDTO struct {
	ID        string          `json:"id"`
	Title     string          `json:"title"`
	Content   json.RawMessage `json:"content"`
	IsStub    bool            `json:"isStub"`
	CreatedAt string          `json:"createdAt"`
	UpdatedAt string          `json:"updatedAt"`
}

// CreateNodeRequestDTO is the transport payload for creating a node.
type CreateNodeRequestDTO struct {
	Title   string          `json:"title"`
	Content json.RawMessage `json:"content"`
	IsStub  bool            `json:"isStub"`
}

// UpdateNodeRequestDTO is the transport payload for updating a node.
type UpdateNodeRequestDTO struct {
	Title   string          `json:"title"`
	Content json.RawMessage `json:"content"`
	IsStub  bool            `json:"isStub"`
}
