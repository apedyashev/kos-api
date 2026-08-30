package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"kos-api/domains/node/persistence/entities"
	"kos-api/domains/node/service/domain"
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

// List returns nodes, newest first, optionally filtered by a case-insensitive title prefix.
func (r *NodeRepository) List(ctx context.Context, titleStartsWith string) ([]entities.NodeEntity, error) {
	var rows pgx.Rows
	var err error
	if titleStartsWith == "" {
		rows, err = r.database.Query(ctx, `
			SELECT id::text, title, content, is_stub, created_at, updated_at
			FROM nodes
			ORDER BY updated_at DESC, id DESC`)
	} else {
		rows, err = r.database.Query(ctx, `
			SELECT id::text, title, content, is_stub, created_at, updated_at
			FROM nodes
			WHERE title ILIKE $1 ESCAPE '\'
			ORDER BY updated_at DESC, id DESC`,
			escapeLikePattern(titleStartsWith)+"%")
	}
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

// ListInSpace returns nodes in a space, newest first, optionally filtered by a case-insensitive title prefix.
func (r *NodeRepository) ListInSpace(ctx context.Context, spaceSlug string, titleStartsWith string) ([]entities.NodeEntity, error) {
	query := `
		SELECT nodes.id::text, nodes.title, nodes.content, nodes.is_stub, nodes.created_at, nodes.updated_at
		FROM nodes
		JOIN knowledge_spaces ON knowledge_spaces.id = nodes.space_id
		WHERE knowledge_spaces.slug = $1`
	arguments := []any{spaceSlug}
	if titleStartsWith != "" {
		query += "\nAND nodes.title ILIKE $2 ESCAPE '\\'"
		arguments = append(arguments, escapeLikePattern(titleStartsWith)+"%")
	}
	query += "\nORDER BY nodes.updated_at DESC, nodes.id DESC"

	rows, err := r.database.Query(ctx, query, arguments...)
	if err != nil {
		return nil, fmt.Errorf("query nodes in space: %w", err)
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

// escapeLikePattern escapes ILIKE wildcard characters so user input is matched literally.
func escapeLikePattern(value string) string {
	replacer := strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`)
	return replacer.Replace(value)
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

// GetByID returns a single node by its id within the specified space.
func (r *NodeRepository) GetByID(ctx context.Context, id string, spaceSlug string) (entities.NodeEntity, error) {
	var entity entities.NodeEntity
	err := r.database.QueryRow(ctx, `
		SELECT nodes.id::text, nodes.title, nodes.content, nodes.is_stub, nodes.created_at, nodes.updated_at
		FROM nodes
		JOIN knowledge_spaces ON knowledge_spaces.id = nodes.space_id
		WHERE nodes.id = $1 AND knowledge_spaces.slug = $2`, id, spaceSlug,
	).Scan(&entity.ID, &entity.Title, &entity.Content, &entity.IsStub, &entity.CreatedAt, &entity.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entities.NodeEntity{}, domain.ErrNodeNotFound
		}
		return entities.NodeEntity{}, fmt.Errorf("query node: %w", err)
	}
	return entity, nil
}

// Update persists changes to an existing node and returns the updated record.
func (r *NodeRepository) Update(ctx context.Context, entity entities.NodeEntity) (entities.NodeEntity, error) {
	content := json.RawMessage(entity.Content)
	err := r.database.QueryRow(ctx, `
		UPDATE nodes
		SET title = $2, content = $3, is_stub = $4, updated_at = now()
		WHERE id = $1
		RETURNING id::text, title, content, is_stub, created_at, updated_at`,
		entity.ID, entity.Title, content, entity.IsStub,
	).Scan(&entity.ID, &entity.Title, &entity.Content, &entity.IsStub, &entity.CreatedAt, &entity.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entities.NodeEntity{}, domain.ErrNodeNotFound
		}
		return entities.NodeEntity{}, fmt.Errorf("update node: %w", err)
	}
	return entity, nil
}
