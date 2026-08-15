package domain

import (
	"bytes"
	"errors"
	"strings"
	"time"
)

// ErrInvalidNode identifies a node that does not satisfy business invariants.
var ErrInvalidNode = errors.New("invalid node")

// NodeDomain is the business representation of a knowledge node.
type NodeDomain struct {
	ID        string
	Title     string
	Content   []byte
	IsStub    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

// CreateNodeInput contains the values needed to create a node.
type CreateNodeInput struct {
	Title   string
	Content []byte
	IsStub  bool
}

// Validate checks business invariants for a new node.
func (input CreateNodeInput) Validate() error {
	if strings.TrimSpace(input.Title) == "" {
		return errors.Join(ErrInvalidNode, errors.New("title is required"))
	}
	if len(bytes.TrimSpace(input.Content)) == 0 {
		return errors.Join(ErrInvalidNode, errors.New("content is required"))
	}
	return nil
}
