package reports

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"

	"green-compass-backend/internal/observations"
)

func TestService_ListPending(t *testing.T) {
	adminID := uuid.New()

	tests := []struct {
		name    string
		caller  Caller
		limit   int
		setup   func(*mockRepo)
		wantLen int
		wantErr error
	}{
		{
			name:   "admin can list pending",
			caller: Caller{UserID: adminID, IsPlatformAdmin: true},
			limit:  0,
			setup: func(r *mockRepo) {
				r.listPending = []Report{{ID: uuid.New(), Status: observations.StatusPending}}
			},
			wantLen: 1,
		},
		{
			name:    "non-admin forbidden",
			caller:  Caller{UserID: adminID},
			limit:   10,
			wantErr: ErrNotAllowed,
		},
		{
			name:    "limit clamped",
			caller:  Caller{UserID: adminID, IsPlatformAdmin: true},
			limit:   200,
			setup:   func(r *mockRepo) { r.listPending = []Report{} },
			wantLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockRepo{}
			if tt.setup != nil {
				tt.setup(repo)
			}
			svc := NewService(repo, nil)
			got, err := svc.ListPending(context.Background(), tt.caller, tt.limit)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("error = %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr != nil {
				return
			}
			if len(got) != tt.wantLen {
				t.Errorf("len = %d, want %d", len(got), tt.wantLen)
			}
		})
	}
}

func TestService_Verify(t *testing.T) {
	obsID := uuid.New()
	adminID := uuid.New()

	tests := []struct {
		name    string
		id      uuid.UUID
		caller  Caller
		status  string
		setup   func(*mockRepo)
		want    *Report
		wantErr error
	}{
		{
			name:   "admin verifies",
			id:     obsID,
			caller: Caller{UserID: adminID, IsPlatformAdmin: true},
			status: observations.StatusVerified,
			setup: func(r *mockRepo) {
				r.verify = &Report{ID: obsID, Status: observations.StatusVerified}
			},
			want: &Report{ID: obsID, Status: observations.StatusVerified},
		},
		{
			name:   "admin rejects",
			id:     obsID,
			caller: Caller{UserID: adminID, IsPlatformAdmin: true},
			status: observations.StatusRejected,
			setup: func(r *mockRepo) {
				r.verify = &Report{ID: obsID, Status: observations.StatusRejected}
			},
			want: &Report{ID: obsID, Status: observations.StatusRejected},
		},
		{
			name:    "non-admin forbidden",
			id:      obsID,
			caller:  Caller{UserID: uuid.New()},
			status:  observations.StatusVerified,
			wantErr: ErrNotAllowed,
		},
		{
			name:    "invalid status",
			id:      obsID,
			caller:  Caller{UserID: adminID, IsPlatformAdmin: true},
			status:  "invalid",
			wantErr: ErrInvalidData,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockRepo{}
			if tt.setup != nil {
				tt.setup(repo)
			}
			svc := NewService(repo, nil)
			got, err := svc.Verify(context.Background(), tt.id, tt.caller, tt.status)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("error = %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr != nil {
				return
			}
			if got.Status != tt.want.Status {
				t.Errorf("Status = %v, want %v", got.Status, tt.want.Status)
			}
		})
	}
}

type mockRepo struct {
	listPending []Report
	listPendErr error
	byID        *Report
	byIDErr     error
	verify      *Report
	verifyErr   error
}

func (m *mockRepo) ListPending(ctx context.Context, limit int) ([]Report, error) {
	return m.listPending, m.listPendErr
}
func (m *mockRepo) ByID(ctx context.Context, id uuid.UUID) (*Report, error) {
	return m.byID, m.byIDErr
}
func (m *mockRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status string, verifiedBy uuid.UUID) (*Report, error) {
	if m.verifyErr != nil {
		return nil, m.verifyErr
	}
	m.verify.Status = status
	return m.verify, nil
}
