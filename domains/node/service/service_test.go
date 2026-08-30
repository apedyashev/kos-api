package service

import (
	"context"
	"errors"
	"testing"

	"kos-api/domains/node/persistence/entities"
	"kos-api/domains/node/service/domain"
)

type fakeRepository struct {
	created      entities.NodeEntity
	updated      entities.NodeEntity
	getSpaceSlug string
}

func (f *fakeRepository) List(context.Context, string) ([]entities.NodeEntity, error) {
	return nil, nil
}

func (f *fakeRepository) Create(_ context.Context, entity entities.NodeEntity) (entities.NodeEntity, error) {
	f.created = entity
	return entity, nil
}

func (f *fakeRepository) GetByID(_ context.Context, id string, spaceSlug string) (entities.NodeEntity, error) {
	f.getSpaceSlug = spaceSlug
	if id != f.created.ID {
		return entities.NodeEntity{}, domain.ErrNodeNotFound
	}
	return f.created, nil
}

func (f *fakeRepository) Update(_ context.Context, entity entities.NodeEntity) (entities.NodeEntity, error) {
	f.updated = entity
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

func TestGetNodePassesSpaceSlugToRepository(t *testing.T) {
	repository := &fakeRepository{created: entities.NodeEntity{ID: "node-1", Title: "Welcome"}}
	_, err := New(repository).GetNode(context.Background(), "node-1", "personal")
	if err != nil {
		t.Fatalf("get node failed: %v", err)
	}
	if repository.getSpaceSlug != "personal" {
		t.Fatalf("expected space slug personal, got %q", repository.getSpaceSlug)
	}
}
