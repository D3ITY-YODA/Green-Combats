package observations

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"

	"green-compass-backend/pkg/storage"
)

func TestService_Create(t *testing.T) {
	userID := uuid.New()
	now := time.Now().UTC()

	tests := []struct {
		name    string
		input   CreateInput
		setup   func(*mockRepo)
		want    *Observation
		wantErr error
	}{
		{
			name: "success with place_id",
			input: CreateInput{
				PlaceID:    &userID,
				Category:   CategoryFlood,
				Description: "River rising",
				Caller:     Caller{UserID: userID},
			},
			setup: func(r *mockRepo) {
				r.create = &Observation{ID: uuid.New(), ReporterID: userID, Category: CategoryFlood, Description: "River rising", Status: StatusPending, CreatedAt: now, UpdatedAt: now}
			},
			want: &Observation{ReporterID: userID, Category: CategoryFlood, Description: "River rising", Status: StatusPending, CreatedAt: now, UpdatedAt: now},
		},
		{
			name: "success with lat lon",
			input: CreateInput{
				Lat:        floatPtr(-1.2),
				Lon:        floatPtr(36.8),
				Category:   CategoryDrought,
				Description: "Dry spell",
				Caller:     Caller{UserID: userID},
			},
			setup: func(r *mockRepo) {
				r.create = &Observation{ID: uuid.New(), ReporterID: userID, Lat: floatPtr(-1.2), Lon: floatPtr(36.8), Category: CategoryDrought, Description: "Dry spell", Status: StatusPending, CreatedAt: now, UpdatedAt: now}
			},
			want: &Observation{ReporterID: userID, Lat: floatPtr(-1.2), Lon: floatPtr(36.8), Category: CategoryDrought, Description: "Dry spell", Status: StatusPending, CreatedAt: now, UpdatedAt: now},
		},
		{
			name: "invalid category",
			input: CreateInput{
				Category:   "unknown",
				Description: "test",
				Caller:     Caller{UserID: userID},
			},
			wantErr: ErrInvalidData,
		},
		{
			name: "missing description",
			input: CreateInput{
				Category:   CategoryFlood,
				Description: "   ",
				Caller:     Caller{UserID: userID},
			},
			wantErr: ErrInvalidData,
		},
		{
			name: "description too long",
			input: CreateInput{
				Category:   CategoryFlood,
				Description: longString(2001),
				Caller:     Caller{UserID: userID},
			},
			wantErr: ErrInvalidData,
		},
		{
			name: "missing location",
			input: CreateInput{
				Category:   CategoryFlood,
				Description: "test",
				Caller:     Caller{UserID: userID},
			},
			wantErr: ErrInvalidData,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockRepo{}
			if tt.setup != nil {
				tt.setup(repo)
			}
			svc := NewService(repo, storage.NewMemoryStore(), nil)
			got, err := svc.Create(context.Background(), tt.input)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("error = %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr != nil {
				return
			}
			if got.Status != tt.want.Status {
				t.Errorf("Status = %v, want %v", got.Status, tt.want.Status)
			}
			if got.Description != tt.want.Description {
				t.Errorf("Description = %v, want %v", got.Description, tt.want.Description)
			}
		})
	}
}

func TestService_ByID(t *testing.T) {
	obsID := uuid.New()
	reporterID := uuid.New()
	adminID := uuid.New()

	tests := []struct {
		name    string
		id      uuid.UUID
		caller  Caller
		setup   func(*mockRepo)
		want    *Observation
		wantErr error
	}{
		{
			name:   "owner can view",
			id:     obsID,
			caller: Caller{UserID: reporterID},
			setup: func(r *mockRepo) {
				r.byID = &Observation{ID: obsID, ReporterID: reporterID}
			},
			want: &Observation{ID: obsID, ReporterID: reporterID},
		},
		{
			name:   "admin can view",
			id:     obsID,
			caller: Caller{UserID: adminID, IsPlatformAdmin: true},
			setup: func(r *mockRepo) {
				r.byID = &Observation{ID: obsID, ReporterID: reporterID}
			},
			want: &Observation{ID: obsID, ReporterID: reporterID},
		},
		{
			name:    "not found",
			id:      obsID,
			caller:  Caller{UserID: reporterID},
			setup:   func(r *mockRepo) { r.byIDErr = ErrNotFound },
			wantErr: ErrNotFound,
		},
		{
			name:   "forbidden for other user",
			id:     obsID,
			caller: Caller{UserID: adminID},
			setup: func(r *mockRepo) {
				r.byID = &Observation{ID: obsID, ReporterID: reporterID}
			},
			wantErr: ErrNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockRepo{}
			if tt.setup != nil {
				tt.setup(repo)
			}
			svc := NewService(repo, storage.NewMemoryStore(), nil)
			got, err := svc.ByID(context.Background(), tt.id, tt.caller)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("error = %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr != nil {
				return
			}
			if got.ID != tt.want.ID {
				t.Errorf("ID = %v, want %v", got.ID, tt.want.ID)
			}
		})
	}
}

func TestService_ListMine(t *testing.T) {
	userID := uuid.New()

	tests := []struct {
		name    string
		userID  uuid.UUID
		limit   int
		setup   func(*mockRepo)
		wantLen int
		wantErr error
	}{
		{
			name:    "default limit",
			userID:  userID,
			limit:   0,
			setup:   func(r *mockRepo) { r.listMine = []Observation{{ID: uuid.New()}} },
			wantLen: 1,
		},
		{
			name:    "custom limit",
			userID:  userID,
			limit:   5,
			setup:   func(r *mockRepo) { r.listMine = []Observation{{ID: uuid.New()}, {ID: uuid.New()}} },
			wantLen: 2,
		},
		{
			name:    "limit clamped",
			userID:  userID,
			limit:   200,
			setup:   func(r *mockRepo) {},
			wantLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockRepo{}
			if tt.setup != nil {
				tt.setup(repo)
			}
			svc := NewService(repo, storage.NewMemoryStore(), nil)
			got, err := svc.ListMine(context.Background(), tt.userID, tt.limit)
			if err != nil {
				t.Fatalf("error = %v", err)
			}
			if len(got) != tt.wantLen {
				t.Errorf("len = %d, want %d", len(got), tt.wantLen)
			}
		})
	}
}

func TestService_Verify(t *testing.T) {
	obsID := uuid.New()

	tests := []struct {
		name    string
		id      uuid.UUID
		caller  Caller
		status  string
		setup   func(*mockRepo)
		want    *Observation
		wantErr error
	}{
		{
			name:   "admin verifies",
			id:     obsID,
			caller: Caller{UserID: uuid.New(), IsPlatformAdmin: true},
			status: StatusVerified,
			setup: func(r *mockRepo) {
				r.verify = &Observation{ID: obsID, Status: StatusVerified}
			},
			want: &Observation{ID: obsID, Status: StatusVerified},
		},
		{
			name:   "admin rejects",
			id:     obsID,
			caller: Caller{UserID: uuid.New(), IsPlatformAdmin: true},
			status: StatusRejected,
			setup: func(r *mockRepo) {
				r.verify = &Observation{ID: obsID, Status: StatusRejected}
			},
			want: &Observation{ID: obsID, Status: StatusRejected},
		},
		{
			name:    "non-admin forbidden",
			id:      obsID,
			caller:  Caller{UserID: uuid.New()},
			status:  StatusVerified,
			wantErr: ErrNotAllowed,
		},
		{
			name:    "invalid status",
			id:      obsID,
			caller:  Caller{UserID: uuid.New(), IsPlatformAdmin: true},
			status:  StatusPending,
			wantErr: ErrInvalidData,
		},
		{
			name:    "unknown status",
			id:      obsID,
			caller:  Caller{UserID: uuid.New(), IsPlatformAdmin: true},
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
			svc := NewService(repo, storage.NewMemoryStore(), nil)
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
	create      *Observation
	createErr   error
	byID        *Observation
	byIDErr     error
	listMine    []Observation
	listMineErr error
	listPending []Observation
	listPendErr error
	verify      *Observation
	verifyErr   error
}

func (m *mockRepo) Create(ctx context.Context, o *Observation) error {
	if m.createErr != nil {
		return m.createErr
	}
	o.ID = uuid.New()
	o.CreatedAt = time.Now().UTC()
	o.UpdatedAt = o.CreatedAt
	m.create = o
	return nil
}
func (m *mockRepo) ByID(ctx context.Context, id uuid.UUID) (*Observation, error) {
	return m.byID, m.byIDErr
}
func (m *mockRepo) ListByReporter(ctx context.Context, reporterID uuid.UUID, limit int) ([]Observation, error) {
	return m.listMine, m.listMineErr
}
func (m *mockRepo) ListPending(ctx context.Context, limit int) ([]Observation, error) {
	return m.listPending, m.listPendErr
}
func (m *mockRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status string, verifiedBy *uuid.UUID, verifiedAt *time.Time) (*Observation, error) {
	if m.verifyErr != nil {
		return nil, m.verifyErr
	}
	m.verify.Status = status
	return m.verify, nil
}

func floatPtr(v float64) *float64 { return &v }
func stringPtr(v string) *string   { return &v }
func longString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = 'x'
	}
	return string(b)
}
