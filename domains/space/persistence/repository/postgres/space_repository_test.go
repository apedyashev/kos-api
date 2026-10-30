package postgres

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"

	"kos-api/domains/space/persistence/entities"
	"kos-api/domains/space/persistence/repository"
)

type fakeDatabase struct {
	query string
	args  []any
	row   pgx.Row
}

func (database *fakeDatabase) Query(context.Context, string, ...any) (pgx.Rows, error) {
	return nil, errors.New("unexpected query")
}

func (database *fakeDatabase) QueryRow(_ context.Context, query string, args ...any) pgx.Row {
	database.query = query
	database.args = args
	return database.row
}

func (database *fakeDatabase) Close() {}

type fakeRow struct {
	err error
}

func (row fakeRow) Scan(destinations ...any) error {
	if row.err != nil {
		return row.err
	}
	*destinations[0].(*string) = "space-1"
	*destinations[1].(*string) = "Saved"
	*destinations[2].(*string) = "saved"
	return nil
}

func TestUpdate(t *testing.T) {
	failure := errors.New("database unavailable")
	for _, test := range []struct {
		name    string
		rowErr  error
		wantErr error
	}{
		{name: "updated"},
		{name: "not found", rowErr: pgx.ErrNoRows, wantErr: repository.ErrSpaceNotFound},
		{name: "database error", rowErr: failure, wantErr: failure},
	} {
		t.Run(test.name, func(t *testing.T) {
			database := &fakeDatabase{row: fakeRow{err: test.rowErr}}
			repo := &SpaceRepository{database: database}
			result, err := repo.Update(context.Background(), entities.SpaceEntity{ID: "space-1", Name: "Example", Slug: "example"})
			if !errors.Is(err, test.wantErr) {
				t.Fatalf("error = %v, want %v", err, test.wantErr)
			}
			if !strings.Contains(database.query, "UPDATE knowledge_spaces") || !strings.Contains(database.query, "RETURNING id::text, name, slug") {
				t.Fatalf("unexpected query: %s", database.query)
			}
			if !reflect.DeepEqual(database.args, []any{"space-1", "Example", "example"}) {
				t.Fatalf("unexpected query arguments: %v", database.args)
			}
			if test.wantErr == nil && (result.ID != "space-1" || result.Name != "Saved" || result.Slug != "saved") {
				t.Fatalf("unexpected result: %+v", result)
			}
		})
	}
}
