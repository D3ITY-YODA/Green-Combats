package audit

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
)

// --- Mock repository ---

type mockRepo struct {
	entries []Entry
}

func newMockRepo() *mockRepo {
	return &mockRepo{}
}

func (m *mockRepo) Create(_ context.Context, e *Entry) error {
	e.ID = uuid.New()
	e.CreatedAt = time.Now()
	m.entries = append(m.entries, *e)
	return nil
}

func (m *mockRepo) List(_ context.Context, req ListRequest) ([]Entry, int, error) {
	var result []Entry
	for _, e := range m.entries {
		if req.ActorID != nil && (e.ActorID == nil || *e.ActorID != *req.ActorID) {
			continue
		}
		if req.Action != nil && e.Action != *req.Action {
			continue
		}
		if req.ResourceType != nil && e.ResourceType != *req.ResourceType {
			continue
		}
		result = append(result, e)
	}
	total := len(result)
	start := (req.Page - 1) * req.Limit
	if start >= total {
		return []Entry{}, total, nil
	}
	end := start + req.Limit
	if end > total {
		end = total
	}
	return result[start:end], total, nil
}

func (m *mockRepo) GetByID(_ context.Context, id uuid.UUID) (*Entry, error) {
	for _, e := range m.entries {
		if e.ID == id {
			return &e, nil
		}
	}
	return nil, nil
}

// --- Tests ---

func TestRecord_Success(t *testing.T) {
	repo := newMockRepo()
	svc := NewService(repo)

	actorID := uuid.New()
	orgID := uuid.New()
	resourceID := uuid.New()

	entry, err := svc.Record(context.Background(), CreateRequest{
		ActorID:      &actorID,
		ActorOrgID:   &orgID,
		Action:       ActionObservationVerified,
		ResourceType: "observation",
		ResourceID:   &resourceID,
		Details:      map[string]interface{}{"notes": "verified by field inspection"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry.ID == uuid.Nil {
		t.Error("expected non-nil ID")
	}
	if entry.Action != ActionObservationVerified {
		t.Errorf("action: got %s, want %s", entry.Action, ActionObservationVerified)
	}
	if len(repo.entries) != 1 {
		t.Errorf("stored count: got %d, want 1", len(repo.entries))
	}
}

func TestRecord_MissingAction(t *testing.T) {
	svc := NewService(newMockRepo())

	_, err := svc.Record(context.Background(), CreateRequest{
		ResourceType: "observation",
	})
	if err == nil {
		t.Fatal("expected error for missing action")
	}
}

func TestRecord_MissingResourceType(t *testing.T) {
	svc := NewService(newMockRepo())

	_, err := svc.Record(context.Background(), CreateRequest{
		Action: ActionObservationVerified,
	})
	if err == nil {
		t.Fatal("expected error for missing resource_type")
	}
}

func TestList_FilterByAction(t *testing.T) {
	repo := newMockRepo()
	svc := NewService(repo)

	// Create entries with different actions
	_, _ = svc.Record(context.Background(), CreateRequest{
		Action:       ActionObservationVerified,
		ResourceType: "observation",
	})
	_, _ = svc.Record(context.Background(), CreateRequest{
		Action:       ActionObservationRejected,
		ResourceType: "observation",
	})
	_, _ = svc.Record(context.Background(), CreateRequest{
		Action:       ActionContentPublished,
		ResourceType: "content",
	})

	action := ActionObservationVerified
	resp, err := svc.List(context.Background(), ListRequest{
		Action: &action,
		Page:   1,
		Limit:  10,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Total != 1 {
		t.Errorf("total: got %d, want 1", resp.Total)
	}
}

func TestList_FilterByActor(t *testing.T) {
	repo := newMockRepo()
	svc := NewService(repo)

	actor1 := uuid.New()
	actor2 := uuid.New()

	_, _ = svc.Record(context.Background(), CreateRequest{
		ActorID:      &actor1,
		Action:       ActionObservationVerified,
		ResourceType: "observation",
	})
	_, _ = svc.Record(context.Background(), CreateRequest{
		ActorID:      &actor2,
		Action:       ActionObservationVerified,
		ResourceType: "observation",
	})

	resp, err := svc.List(context.Background(), ListRequest{
		ActorID: &actor1,
		Page:    1,
		Limit:   10,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Total != 1 {
		t.Errorf("total: got %d, want 1", resp.Total)
	}
}

func TestGetByID(t *testing.T) {
	repo := newMockRepo()
	svc := NewService(repo)

	created, _ := svc.Record(context.Background(), CreateRequest{
		Action:       ActionUserCreated,
		ResourceType: "user",
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
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
