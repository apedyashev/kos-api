package service

import (
	"context"
	"errors"
	"testing"

	"kos-api/domains/node/persistence/entities"
	"kos-api/domains/node/service/domain"
)

type fakeRepository struct {
	created entities.NodeEntity
}

func (f *fakeRepository) List(context.Context) ([]entities.NodeEntity, error) {
	return nil, nil
}

func (f *fakeRepository) Create(_ context.Context, entity entities.NodeEntity) (entities.NodeEntity, error) {
	f.created = entity
	return entity, nil
}

func TestCreateNodeValidatesBeforePersistence(t *testing.T) {
	repository := &fakeRepository{}
	_, err := New(repository).CreateNode(context.Background(), domain.CreateNodeInput{Content: []byte(`{"type":"doc"}`)})
	if !errors.Is(err, domain.ErrInvalidNode) {
		t.Fatalf("expected invalid node error, got %v", err)
	}
	if repository.created.Title != "" {
		t.Fatal("repository should not be called for invalid input")
	}
}

func TestCreateNodeMapsDomainAndEntity(t *testing.T) {
	repository := &fakeRepository{}
	created, err := New(repository).CreateNode(context.Background(), domain.CreateNodeInput{
		Title: "Welcome", Content: []byte(`{"type":"doc"}`), IsStub: true,
	})
	if err != nil {
		t.Fatalf("create node failed: %v", err)
	}
	if created.Title != "Welcome" || !repository.created.IsStub {
		t.Fatalf("unexpected node mapping: domain=%+v entity=%+v", created, repository.created)
	}
}
