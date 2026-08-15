package postgres

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"kos-api/internal/domain"
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
func (r *NodeRepository) List(ctx context.Context) ([]domain.NodeDomain, error) {
	rows, err := r.database.Query(ctx, `
		SELECT id::text, title, content, is_stub, created_at, updated_at
		FROM nodes
		ORDER BY updated_at DESC, id DESC`)
	if err != nil {
		return nil, fmt.Errorf("query nodes: %w", err)
	}
	defer rows.Close()

	nodes := make([]domain.NodeDomain, 0)
	for rows.Next() {
		var node domain.NodeDomain
		if err := rows.Scan(&node.ID, &node.Title, &node.Content, &node.IsStub, &node.CreatedAt, &node.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan node: %w", err)
		}
		nodes = append(nodes, node)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate nodes: %w", err)
	}
	return nodes, nil
}

// Create inserts a node and returns the persisted record.
func (r *NodeRepository) Create(ctx context.Context, node domain.NodeDomain) (domain.NodeDomain, error) {
	content := json.RawMessage(node.Content)
	err := r.database.QueryRow(ctx, `
		INSERT INTO nodes (title, content, is_stub)
		VALUES ($1, $2, $3)
		RETURNING id::text, title, content, is_stub, created_at, updated_at`,
		node.Title, content, node.IsStub,
	).Scan(&node.ID, &node.Title, &node.Content, &node.IsStub, &node.CreatedAt, &node.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.NodeDomain{}, fmt.Errorf("insert node returned no row: %w", err)
		}
		return domain.NodeDomain{}, fmt.Errorf("insert node: %w", err)
	}
	return node, nil
}
