package updates

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

type mockUpdatesRepo struct {
	todayUpdates map[uuid.UUID]*Update
	explored     []ExploreIndicator
}

func (m *mockUpdatesRepo) GetTodayForPlace(ctx context.Context, placeID uuid.UUID) (*Update, error) {
	if update, ok := m.todayUpdates[placeID]; ok {
		return update, nil
	}
	return nil, ErrNotFound
}

func (m *mockUpdatesRepo) ListForUser(ctx context.Context, req ListRequest) ([]Update, int, error) {
	return []Update{}, 0, nil
}

func (m *mockUpdatesRepo) ListIndicators(ctx context.Context, contentID uuid.UUID) ([]Indicator, error) {
	return []Indicator{}, nil
}

func (m *mockUpdatesRepo) ExploreIndicators(ctx context.Context, req ExploreRequest) ([]ExploreIndicator, int, error) {
	return m.explored, len(m.explored), nil
}

func newMockUpdatesRepo() *mockUpdatesRepo {
	return &mockUpdatesRepo{
		todayUpdates: make(map[uuid.UUID]*Update),
		explored:     make([]ExploreIndicator, 0),
	}
}

func TestGetTodayNotFound(t *testing.T) {
	repo := newMockUpdatesRepo()
	svc := NewService(repo, nil)

	placeID := uuid.New()
	_, err := svc.GetToday(context.Background(), placeID)

	if err != ErrNotFound {
		t.Errorf("GetToday: got %v, want ErrNotFound", err)
	}
}

func TestListPaginationValidation(t *testing.T) {
	repo := newMockUpdatesRepo()
	svc := NewService(repo, nil)

	tests := []struct {
		name       string
		inputPage  int
		inputLimit int
		wantPage   int
		wantLimit  int
	}{
		{
			name:       "invalid page (negative)",
			inputPage:  -1,
			inputLimit: 20,
			wantPage:   1,
			wantLimit:  20,
		},
		{
			name:       "invalid limit (too large)",
			inputPage:  1,
			inputLimit: 500,
			wantPage:   1,
			wantLimit:  100, // max limit
		},
		{
			name:       "valid pagination",
			inputPage:  2,
			inputLimit: 25,
			wantPage:   2,
			wantLimit:  25,
		},
		{
			name:       "zero limit defaults to 20",
			inputPage:  1,
			inputLimit: 0,
			wantPage:   1,
			wantLimit:  20,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := svc.List(context.Background(), ListRequest{
				UserID: uuid.New(),
				Page:   tt.inputPage,
				Limit:  tt.inputLimit,
			})

			if err != nil {
				t.Fatalf("List failed: %v", err)
			}

			if resp.Page != tt.wantPage {
				t.Errorf("Page: got %d, want %d", resp.Page, tt.wantPage)
			}

			if resp.Limit != tt.wantLimit {
				t.Errorf("Limit: got %d, want %d", resp.Limit, tt.wantLimit)
			}
		})
	}
}

func TestExplorePaginationValidation(t *testing.T) {
	repo := newMockUpdatesRepo()
	svc := NewService(repo, nil)

	resp, err := svc.Explore(context.Background(), ExploreRequest{
		UserID: uuid.New(),
		Page:   -1,
		Limit:  999,
	})

	if err != nil {
		t.Fatalf("Explore failed: %v", err)
	}

	if resp.Page != 1 {
		t.Errorf("Page: got %d, want 1", resp.Page)
	}

	if resp.Limit != 100 {
		t.Errorf("Limit: got %d, want 100 (max)", resp.Limit)
	}
}
