package service

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"kos-api/domains/space/persistence/entities"
	"kos-api/domains/space/persistence/repository"
	"kos-api/domains/space/service/domain"
)

type fakeRepository struct {
	updated entities.SpaceEntity
	result  entities.SpaceEntity
	err     error
}

func (fake *fakeRepository) List(context.Context, domain.ListSpacesInput) ([]entities.SpaceEntity, error) {
	return nil, nil
}

func (fake *fakeRepository) Update(_ context.Context, entity entities.SpaceEntity) (entities.SpaceEntity, error) {
	fake.updated = entity
	return fake.result, fake.err
}

func TestUpdateSpace(t *testing.T) {
	input := domain.UpdateSpaceInput{ID: "space-1", Name: "Example", Slug: "example"}
	t.Run("maps update and result", func(t *testing.T) {
		fake := &fakeRepository{result: entities.SpaceEntity{ID: "space-1", Name: "Saved", Slug: "saved"}}
		result, err := New(fake).UpdateSpace(context.Background(), input)
		if err != nil {
			t.Fatal(err)
		}
		if fake.updated.ID != input.ID || fake.updated.Name != input.Name || fake.updated.Slug != input.Slug {
			t.Fatalf("unexpected entity: %+v", fake.updated)
		}
		if result.ID != "space-1" || result.Name != "Saved" || result.Slug != "saved" {
			t.Fatalf("unexpected domain result: %+v", result)
		}
	})
	t.Run("maps not found", func(t *testing.T) {
		fake := &fakeRepository{err: fmt.Errorf("database: %w", repository.ErrSpaceNotFound)}
		_, err := New(fake).UpdateSpace(context.Background(), input)
		if !errors.Is(err, domain.ErrSpaceNotFound) {
			t.Fatalf("expected domain not found, got %v", err)
		}
	})
	t.Run("preserves unexpected errors", func(t *testing.T) {
		failure := errors.New("database unavailable")
		fake := &fakeRepository{err: failure}
		_, err := New(fake).UpdateSpace(context.Background(), input)
		if !errors.Is(err, failure) {
			t.Fatalf("expected database error, got %v", err)
		}
	})
}
