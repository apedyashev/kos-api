package postgres

import (
	"context"
	"fmt"

	"kos-api/domains/space/persistence/entities"
	"kos-api/domains/space/service/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// SpaceRepository persists spaces in PostgreSQL.
type SpaceRepository struct {
	database *pgxpool.Pool
}

// NewSpaceRepository creates a PostgreSQL connection pool from a Supabase URL.
func NewSpaceRepository(databaseURL string) (*SpaceRepository, error) {
	databaseConfig, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("parse database URL: %w", err)
	}
	// Supabase's transaction pooler can route statements across backend
	// connections, so named prepared statements can collide between sessions.
	databaseConfig.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	database, err := pgxpool.NewWithConfig(context.Background(), databaseConfig)
	if err != nil {
		return nil, fmt.Errorf("create database pool: %w", err)
	}
	return &SpaceRepository{
		database: database,
	}, nil
}

// Close releases the PostgreSQL resources.
func (r *SpaceRepository) Close() {
	r.database.Close()
}

// List returns spaces.
func (r *SpaceRepository) List(ctx context.Context, input domain.ListSpacesInput) ([]entities.SpaceEntity, error) {
	var rows pgx.Rows
	var err error

	if input.OwnerID == "" {
		return nil, fmt.Errorf("owner ID is required")
	}

	rows, err = r.database.Query(ctx, `
		SELECT id::text, name, slug 
		FROM knowledge_spaces
		WHERE owner_id = $1
		ORDER BY created_at DESC, id DESC`, input.OwnerID)

	if err != nil {
		return nil, fmt.Errorf("query spaces: %w", err)
	}
	defer rows.Close()

	spaces := make([]entities.SpaceEntity, 0)
	for rows.Next() {
		var entity entities.SpaceEntity
		if err := rows.Scan(&entity.ID, &entity.Name, &entity.Slug); err != nil {
			return nil, fmt.Errorf("scan knowledge_spaces: %w", err)
		}
		spaces = append(spaces, entity)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate knowledge_spaces: %w", err)
	}
	return spaces, nil
}
