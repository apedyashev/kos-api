package service

import (
	"context"
	"errors"
	"testing"

	"kos-api/internal/domain"
)

type fakeNodeRepository struct {
	created domain.NodeDomain
	err     error
}

func (f *fakeNodeRepository) List(context.Context) ([]domain.NodeDomain, error) {
	return nil, f.err
}

func (f *fakeNodeRepository) Create(_ context.Context, node domain.NodeDomain) (domain.NodeDomain, error) {
	if f.err != nil {
		return domain.NodeDomain{}, f.err
	}
	f.created = node
	return node, nil
}

func TestCreateNodeValidatesBeforePersistence(t *testing.T) {
	repository := &fakeNodeRepository{}
	service := NewNodeService(repository)

	_, err := service.CreateNode(context.Background(), domain.CreateNodeInput{Content: []byte(`{"type":"doc"}`)})
	if !errors.Is(err, domain.ErrInvalidNode) {
		t.Fatalf("expected invalid node error, got %v", err)
	}
	if repository.created.Title != "" {
		t.Fatal("repository should not be called for invalid input")
	}
}

func TestCreateNodePersistsValidatedInput(t *testing.T) {
	repository := &fakeNodeRepository{}
	service := NewNodeService(repository)
	content := []byte(`{"type":"doc","content":[]}`)

	created, err := service.CreateNode(context.Background(), domain.CreateNodeInput{
		Title:   "Welcome",
		Content: content,
		IsStub:  true,
	})
	if err != nil {
		t.Fatalf("create node failed: %v", err)
	}
	if string(repository.created.Content) != string(content) || !repository.created.IsStub {
		t.Fatalf("repository received unexpected node: %+v", repository.created)
	}
	if created.Title != "Welcome" {
		t.Fatalf("created node title = %q, want Welcome", created.Title)
	}
}
