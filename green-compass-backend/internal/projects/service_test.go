package projects

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
)

// --- Mock repository ---

type mockRepo struct {
	projects []Project
}

func newMockRepo() *mockRepo {
	return &mockRepo{}
}

func (m *mockRepo) Create(_ context.Context, p *Project) error {
	p.ID = uuid.New()
	p.Status = "active"
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
	m.projects = append(m.projects, *p)
	return nil
}

func (m *mockRepo) GetByID(_ context.Context, id uuid.UUID) (*Project, error) {
	for _, p := range m.projects {
		if p.ID == id {
			return &p, nil
		}
	}
	return nil, ErrNotFound
}

func (m *mockRepo) Update(_ context.Context, id uuid.UUID, req UpdateRequest) (*Project, error) {
	for i := range m.projects {
		if m.projects[i].ID == id {
			if req.Name != nil {
				m.projects[i].Name = *req.Name
			}
			if req.Description != nil {
				m.projects[i].Description = *req.Description
			}
			if req.Status != nil {
				m.projects[i].Status = *req.Status
			}
			m.projects[i].UpdatedAt = time.Now()
			return &m.projects[i], nil
		}
	}
	return nil, ErrNotFound
}

func (m *mockRepo) List(_ context.Context, req ListRequest) ([]Project, int, error) {
	var result []Project
	for _, p := range m.projects {
		if p.OrgID == req.OrgID {
			result = append(result, p)
		}
	}
	total := len(result)
	start := (req.Page - 1) * req.Limit
	if start >= total {
		return []Project{}, total, nil
	}
	end := start + req.Limit
	if end > total {
		end = total
	}
	return result[start:end], total, nil
}

func (m *mockRepo) Delete(_ context.Context, id uuid.UUID) error {
	for i := range m.projects {
		if m.projects[i].ID == id {
			m.projects = append(m.projects[:i], m.projects[i+1:]...)
			return nil
		}
	}
	return ErrNotFound
}

// --- Tests ---

func TestCreate_Success(t *testing.T) {
	svc := NewService(newMockRepo())
	orgID := uuid.New()

	p, err := svc.Create(context.Background(), CreateRequest{
		OrgID:       orgID,
		Name:        "Climate Resilience Program",
		Description: "Building climate resilience in rural communities",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Name != "Climate Resilience Program" {
		t.Errorf("name: got %q, want %q", p.Name, "Climate Resilience Program")
	}
	if p.Status != "active" {
		t.Errorf("status: got %q, want active", p.Status)
	}
}

func TestCreate_EmptyName(t *testing.T) {
	svc := NewService(newMockRepo())

	_, err := svc.Create(context.Background(), CreateRequest{
		OrgID: uuid.New(),
	})
	if err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestGetByID_Found(t *testing.T) {
	repo := newMockRepo()
	svc := NewService(repo)

	created, _ := svc.Create(context.Background(), CreateRequest{
		OrgID: uuid.New(),
		Name:  "Test Project",
	})

	got, err := svc.GetByID(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != created.ID {
		t.Errorf("ID mismatch: got %v, want %v", got.ID, created.ID)
	}
}

func TestGetByID_NotFound(t *testing.T) {
	svc := NewService(newMockRepo())

	_, err := svc.GetByID(context.Background(), uuid.New())
	if err != ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestUpdate(t *testing.T) {
	repo := newMockRepo()
	svc := NewService(repo)

	created, _ := svc.Create(context.Background(), CreateRequest{
		OrgID: uuid.New(),
		Name:  "Original Name",
	})

	newName := "Updated Name"
	updated, err := svc.Update(context.Background(), created.ID, UpdateRequest{
		Name: &newName,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.Name != "Updated Name" {
		t.Errorf("name: got %q, want %q", updated.Name, "Updated Name")
	}
}

func TestList(t *testing.T) {
	repo := newMockRepo()
	svc := NewService(repo)
	orgID := uuid.New()

	for i := 0; i < 3; i++ {
		_, _ = svc.Create(context.Background(), CreateRequest{
			OrgID: orgID,
			Name:  "Project",
		})
	}

	resp, err := svc.List(context.Background(), ListRequest{
		OrgID: orgID,
		Page:  1,
		Limit: 10,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Total != 3 {
		t.Errorf("total: got %d, want 3", resp.Total)
	}
}

func TestDelete(t *testing.T) {
	repo := newMockRepo()
	svc := NewService(repo)

	created, _ := svc.Create(context.Background(), CreateRequest{
		OrgID: uuid.New(),
		Name:  "To Delete",
	})

	if err := svc.Delete(context.Background(), created.ID); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err := svc.GetByID(context.Background(), created.ID)
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound after delete, got %v", err)
	}
}
