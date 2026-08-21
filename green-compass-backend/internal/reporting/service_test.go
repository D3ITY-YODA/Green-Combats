package reporting

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
)

// --- Mock repository ---

type mockRepo struct {
	deliveryStats map[uuid.UUID]*DeliveryStats
	reportStats   map[uuid.UUID]*ReportStats
}

func newMockRepo() *mockRepo {
	return &mockRepo{
		deliveryStats: make(map[uuid.UUID]*DeliveryStats),
		reportStats:   make(map[uuid.UUID]*ReportStats),
	}
}

func (m *mockRepo) DeliveryStats(_ context.Context, placeID uuid.UUID, from, to string) (*DeliveryStats, error) {
	if s, ok := m.deliveryStats[placeID]; ok {
		return s, nil
	}
	return &DeliveryStats{
		PlaceID:     placeID,
		ByChannel:   make(map[string]int),
		PeriodStart: time.Now().UTC().AddDate(0, 0, -30),
		PeriodEnd:   time.Now().UTC(),
	}, nil
}

func (m *mockRepo) ReportStats(_ context.Context, orgID uuid.UUID, from, to string) (*ReportStats, error) {
	if s, ok := m.reportStats[orgID]; ok {
		return s, nil
	}
	return &ReportStats{
		OrgID:       orgID,
		PeriodStart: time.Now().UTC().AddDate(0, 0, -30),
		PeriodEnd:   time.Now().UTC(),
	}, nil
}

// --- Tests ---

func TestGetDeliveryStats_Success(t *testing.T) {
	repo := newMockRepo()
	svc := NewService(repo)
	placeID := uuid.New()

	now := time.Now().UTC()
	stats, err := svc.GetDeliveryStats(context.Background(), placeID, now.AddDate(0, 0, -30).Format(time.RFC3339), now.Format(time.RFC3339))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stats == nil {
		t.Fatal("expected stats, got nil")
	}
	if stats.PlaceID != placeID {
		t.Errorf("place_id: got %v, want %v", stats.PlaceID, placeID)
	}
}

func TestGetDeliveryStats_MissingPeriod(t *testing.T) {
	svc := NewService(newMockRepo())

	_, err := svc.GetDeliveryStats(context.Background(), uuid.New(), "", "2024-01-01T00:00:00Z")
	if err == nil {
		t.Fatal("expected error for missing period")
	}
}

func TestGetReportStats_Success(t *testing.T) {
	repo := newMockRepo()
	svc := NewService(repo)
	orgID := uuid.New()

	now := time.Now().UTC()
	stats, err := svc.GetReportStats(context.Background(), orgID, now.AddDate(0, 0, -30).Format(time.RFC3339), now.Format(time.RFC3339))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stats == nil {
		t.Fatal("expected stats, got nil")
	}
	if stats.OrgID != orgID {
		t.Errorf("org_id: got %v, want %v", stats.OrgID, orgID)
	}
}

func TestGetReportStats_MissingPeriod(t *testing.T) {
	svc := NewService(newMockRepo())

	_, err := svc.GetReportStats(context.Background(), uuid.New(), "", "")
	if err == nil {
		t.Fatal("expected error for missing period")
	}
}

func TestGetDashboard_WithPlaceID(t *testing.T) {
	svc := NewService(newMockRepo())
	placeID := uuid.New()

	now := time.Now().UTC()
	dashboard, err := svc.GetDashboard(context.Background(), QueryRequest{
		PlaceID:     &placeID,
		PeriodStart: now.AddDate(0, 0, -7),
		PeriodEnd:   now,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dashboard.Delivery == nil {
		t.Error("expected delivery stats")
	}
}

func TestGetDashboard_WithOrgID(t *testing.T) {
	svc := NewService(newMockRepo())
	orgID := uuid.New()

	now := time.Now().UTC()
	dashboard, err := svc.GetDashboard(context.Background(), QueryRequest{
		OrgID:       &orgID,
		PeriodStart: now.AddDate(0, 0, -7),
		PeriodEnd:   now,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dashboard.Reports == nil {
		t.Error("expected report stats")
	}
}

func TestGetDashboard_Empty(t *testing.T) {
	svc := NewService(newMockRepo())

	now := time.Now().UTC()
	dashboard, err := svc.GetDashboard(context.Background(), QueryRequest{
		PeriodStart: now.AddDate(0, 0, -7),
		PeriodEnd:   now,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dashboard.Delivery != nil {
		t.Error("expected nil delivery stats")
	}
	if dashboard.Reports != nil {
		t.Error("expected nil report stats")
	}
}
