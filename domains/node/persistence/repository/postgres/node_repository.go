package postgres

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"kos-api/domains/node/persistence/entities"
)

// NodeRepository persists nodes in PostgreSQL.
type NodeRepository struct {
	database *pgxpool.Pool
}

// NewNodeRepository creates a PostgreSQL connection pool from a Supabase URL.
func NewNodeRepository(databaseURL string) (*NodeRepository, error) {
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
	return &NodeRepository{database: database}, nil
}

// Close releases the PostgreSQL connection pool.
func (r *NodeRepository) Close() {
	r.database.Close()
}

// List returns all nodes, newest first.
func (r *NodeRepository) List(ctx context.Context) ([]entities.NodeEntity, error) {
	rows, err := r.database.Query(ctx, `
		SELECT id::text, title, content, is_stub, created_at, updated_at
		FROM nodes
		ORDER BY updated_at DESC, id DESC`)
	if err != nil {
		return nil, fmt.Errorf("query nodes: %w", err)
	}
	defer rows.Close()

	nodes := make([]entities.NodeEntity, 0)
	for rows.Next() {
		var entity entities.NodeEntity
		if err := rows.Scan(&entity.ID, &entity.Title, &entity.Content, &entity.IsStub, &entity.CreatedAt, &entity.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan node: %w", err)
		}
		nodes = append(nodes, entity)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate nodes: %w", err)
	}
	return nodes, nil
}

// Create inserts a node and returns the persisted record.
func (r *NodeRepository) Create(ctx context.Context, entity entities.NodeEntity) (entities.NodeEntity, error) {
	content := json.RawMessage(entity.Content)
	err := r.database.QueryRow(ctx, `
		INSERT INTO nodes (title, content, is_stub)
		VALUES ($1, $2, $3)
		RETURNING id::text, title, content, is_stub, created_at, updated_at`,
		entity.Title, content, entity.IsStub,
	).Scan(&entity.ID, &entity.Title, &entity.Content, &entity.IsStub, &entity.CreatedAt, &entity.UpdatedAt)
	if err != nil {
		return entities.NodeEntity{}, fmt.Errorf("insert node: %w", err)
	}
	return entity, nil
}
