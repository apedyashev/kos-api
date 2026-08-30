package entities

import "time"

// SpaceEntity represents the spaces table in PostgreSQL.
type SpaceEntity struct {
	ID        string
	Name      string
	Slug      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
