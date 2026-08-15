package entities

import "time"

// NodeEntity represents the nodes table in PostgreSQL.
type NodeEntity struct {
	ID        string
	Title     string
	Content   []byte
	IsStub    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
