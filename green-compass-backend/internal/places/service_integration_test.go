package places_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"

	"green-compass-backend/internal/places"
	"green-compass-backend/internal/testutil"
	"green-compass-backend/internal/users"
)

type svcFixture struct {
	svc   *places.Service
	admin *places.Caller
	user  *places.Caller
	other *places.Caller
}

func newSvcFixture(t *testing.T) *svcFixture {
	t.Helper()
	pool := testutil.TestPool(t)
	ctx := context.Background()

	repo := places.NewRepository(pool)
	svc := places.NewService(repo)

	userSvc := users.NewService(users.NewRepository(pool))
	register := func(name string) *places.Caller {
		u, err := userSvc.Register(ctx, users.RegisterInput{
			Email:       name + "-" + uuid.NewString() + "@greencompass.test",
			DisplayName: name,
			Password:    "long-enough-password",
		})
		if err != nil {
			t.Fatalf("seed %s: %v", name, err)
		}
		return &places.Caller{UserID: u.ID}
	}

	admin := register("Platform Admin")
	if _, err := pool.Exec(ctx, `UPDATE users SET is_platform_admin = TRUE WHERE id = $1`, admin.UserID); err != nil {
		t.Fatalf("promote admin: %v", err)
	}

	return &svcFixture{
		svc:   svc,
		admin: admin,
		user:  register("Regular User"),
		other: register("Other User"),
	}
}

func TestService_Create_Authorization_Integration(t *testing.T) {
	f := newSvcFixture(t)
	ctx := context.Background()

	custom, err := f.svc.Create(ctx, places.CreateInput{
		Name: "My Farm", PlaceType: places.TypeCustom, Lat: -1.1, Lon: 36.9, Caller: *f.user,
	})
	if err != nil {
		t.Fatalf("custom place by regular user: %v", err)
	}
	if custom.CreatedBy == nil || *custom.CreatedBy != f.user.UserID {
		t.Error("CreatedBy should record the caller")
	}

	if _, err := f.svc.Create(ctx, places.CreateInput{
		Name: "Ward X", PlaceType: places.TypeWard, Lat: -1.1, Lon: 36.9, Caller: *f.user,
	}); !errors.Is(err, places.ErrNotAllowed) {
		t.Errorf("official place by regular user err = %v, want ErrNotAllowed", err)
	}

	ward, err := f.svc.Create(ctx, places.CreateInput{
		Name: "Ward Y", PlaceType: places.TypeWard, Lat: -1.2, Lon: 36.8, Caller: *f.admin,
	})
	if err != nil {
		t.Fatalf("official place by admin: %v", err)
	}
	if ward.PlaceType != places.TypeWard {
		t.Errorf("place type = %q, want ward", ward.PlaceType)
	}
}

func TestService_Create_Validation_Integration(t *testing.T) {
	f := newSvcFixture(t)
	ctx := context.Background()

	tests := []struct {
		name    string
		input   places.CreateInput
		wantErr error
	}{
		{
			name:    "empty name",
			input:   places.CreateInput{Name: "  ", PlaceType: places.TypeCustom, Lat: 0, Lon: 0, Caller: *f.user},
			wantErr: places.ErrInvalidData,
		},
		{
			name:    "unknown type",
			input:   places.CreateInput{Name: "X", PlaceType: "province", Lat: 0, Lon: 0, Caller: *f.user},
			wantErr: places.ErrInvalidData,
		},
		{
			name:    "latitude out of range",
			input:   places.CreateInput{Name: "X", PlaceType: places.TypeCustom, Lat: 91, Lon: 0, Caller: *f.user},
			wantErr: places.ErrInvalidData,
		},
		{
			name:    "longitude out of range",
			input:   places.CreateInput{Name: "X", PlaceType: places.TypeCustom, Lat: 0, Lon: -181, Caller: *f.user},
			wantErr: places.ErrInvalidData,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := f.svc.Create(ctx, tt.input); !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestService_Update_Authorization_Integration(t *testing.T) {
	f := newSvcFixture(t)
	ctx := context.Background()

	ownCustom, err := f.svc.Create(ctx, places.CreateInput{
		Name: "Own Pin", PlaceType: places.TypeCustom, Lat: -1.1, Lon: 36.9, Caller: *f.user,
	})
	if err != nil {
		t.Fatalf("seed own custom: %v", err)
	}
	otherCustom, err := f.svc.Create(ctx, places.CreateInput{
		Name: "Other Pin", PlaceType: places.TypeCustom, Lat: -1.1, Lon: 36.9, Caller: *f.other,
	})
	if err != nil {
		t.Fatalf("seed other custom: %v", err)
	}
	official, err := f.svc.Create(ctx, places.CreateInput{
		Name: "Official Ward", PlaceType: places.TypeDistrict, Lat: -1.1, Lon: 36.9, Caller: *f.admin,
	})
	if err != nil {
		t.Fatalf("seed official: %v", err)
	}

	newName := "Renamed"
	newLat := -1.15

	if _, err := f.svc.Update(ctx, ownCustom.ID, places.UpdateInput{Name: &newName}, *f.user); err != nil {
		t.Errorf("creator editing own custom: %v", err)
	}
	if _, err := f.svc.Update(ctx, otherCustom.ID, places.UpdateInput{Name: &newName}, *f.user); !errors.Is(err, places.ErrNotAllowed) {
		t.Errorf("editing stranger's custom err = %v, want ErrNotAllowed", err)
	}
	if _, err := f.svc.Update(ctx, official.ID, places.UpdateInput{Name: &newName}, *f.user); !errors.Is(err, places.ErrNotAllowed) {
		t.Errorf("non-admin editing official err = %v, want ErrNotAllowed", err)
	}

	updated, err := f.svc.Update(ctx, official.ID, places.UpdateInput{Lat: &newLat}, *f.admin)
	if err != nil {
		t.Fatalf("admin editing official: %v", err)
	}
	if updated.Lat != newLat || updated.Name != official.Name {
		t.Errorf("partial update mismatch: name=%q lat=%v", updated.Name, updated.Lat)
	}

	if _, err := f.svc.Update(ctx, uuid.New(), places.UpdateInput{Name: &newName}, *f.admin); !errors.Is(err, places.ErrNotFound) {
		t.Errorf("unknown id err = %v, want ErrNotFound", err)
	}
}

func TestService_Update_Validation_Integration(t *testing.T) {
	f := newSvcFixture(t)
	ctx := context.Background()

	p, err := f.svc.Create(ctx, places.CreateInput{
		Name: "Valid", PlaceType: places.TypeCustom, Lat: -1.1, Lon: 36.9, Caller: *f.user,
	})
	if err != nil {
		t.Fatalf("seed: %v", err)
	}

	badName := "   "
	badLat := 95.0

	if _, err := f.svc.Update(ctx, p.ID, places.UpdateInput{Name: &badName}, *f.user); !errors.Is(err, places.ErrInvalidData) {
		t.Errorf("blank name err = %v, want ErrInvalidData", err)
	}
	if _, err := f.svc.Update(ctx, p.ID, places.UpdateInput{Lat: &badLat}, *f.user); !errors.Is(err, places.ErrInvalidData) {
		t.Errorf("bad lat err = %v, want ErrInvalidData", err)
	}
}

func TestService_Nearby_Validation_Integration(t *testing.T) {
	f := newSvcFixture(t)
	ctx := context.Background()

	if _, err := f.svc.Nearby(ctx, 91, 0, 1000, 10); !errors.Is(err, places.ErrInvalidData) {
		t.Errorf("bad lat err = %v, want ErrInvalidData", err)
	}
	if _, err := f.svc.Nearby(ctx, 0, 0, -1, 10); !errors.Is(err, places.ErrInvalidData) {
		t.Errorf("negative radius err = %v, want ErrInvalidData", err)
	}
	if _, err := f.svc.Nearby(ctx, 0, 0, 1000, -1); !errors.Is(err, places.ErrInvalidData) {
		t.Errorf("negative limit err = %v, want ErrInvalidData", err)
	}
}
