package places_test

import (
	"context"
	"math"
	"strings"
	"testing"

	"github.com/google/uuid"

	"green-compass-backend/internal/places"
	"green-compass-backend/internal/testutil"
)

func TestRepository_CreateAndRoundtrip_Integration(t *testing.T) {
	pool := testutil.TestPool(t)
	repo := places.NewRepository(pool)
	ctx := context.Background()

	created := &places.Place{
		Name:      "Nairobi CBD",
		PlaceType: places.TypeCommunity,
		Lat:       -1.2864,
		Lon:       36.8172,
	}
	if err := repo.Create(ctx, created); err != nil {
		t.Fatalf("Create: %v", err)
	}
	if created.ID == uuid.Nil {
		t.Fatal("Create did not populate ID")
	}
	if created.CreatedAt.IsZero() || created.UpdatedAt.IsZero() {
		t.Fatal("Create did not populate timestamps")
	}

	got, err := repo.ByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("ByID: %v", err)
	}
	if got.Name != "Nairobi CBD" || got.PlaceType != places.TypeCommunity {
		t.Errorf("roundtrip mismatch: %+v", got)
	}
	if math.Abs(got.Lat-(-1.2864)) > 1e-9 || math.Abs(got.Lon-36.8172) > 1e-9 {
		t.Errorf("coordinates = (%v, %v), want (-1.2864, 36.8172)", got.Lat, got.Lon)
	}
	if got.CreatedBy != nil {
		t.Errorf("CreatedBy = %v, want NULL when unset", *got.CreatedBy)
	}
}

func TestRepository_Nearby_Integration(t *testing.T) {
	pool := testutil.TestPool(t)
	repo := places.NewRepository(pool)
	ctx := context.Background()

	seed := []places.Place{
		{Name: uniqueName(t, "Nairobi CBD"), PlaceType: places.TypeCommunity, Lat: -1.2864, Lon: 36.8172},
		{Name: uniqueName(t, "Kisumu Central"), PlaceType: places.TypeWard, Lat: -0.0917, Lon: 34.7680},
		{Name: uniqueName(t, "Mombasa Island"), PlaceType: places.TypeDistrict, Lat: -4.0435, Lon: 39.6682},
		{Name: uniqueName(t, "Someone's Farm"), PlaceType: places.TypeCustom, Lat: -1.2865, Lon: 36.8173},
	}
	for i := range seed {
		if err := repo.Create(ctx, &seed[i]); err != nil {
			t.Fatalf("seed %s: %v", seed[i].Name, err)
		}
	}

	t.Run("orders by real distance and excludes custom places", func(t *testing.T) {
		results, err := repo.Nearby(ctx, -1.2921, 36.8219, 600_000, 10)
		if err != nil {
			t.Fatalf("Nearby: %v", err)
		}
		if len(results) < 3 {
			t.Fatalf("got %d results, want at least the 3 seeded official places", len(results))
		}

		var names []string
		for _, r := range results {
			names = append(names, r.Name)
			if r.PlaceType == places.TypeCustom {
				t.Errorf("custom place %q leaked into nearby results", r.Name)
			}
			if r.DistanceMeters <= 0 {
				t.Errorf("place %q has non-positive distance %f", r.Name, r.DistanceMeters)
			}
		}

		nairobiIdx, kisumuIdx, mombasaIdx := -1, -1, -1
		for i, n := range names {
			switch {
			case contains(n, "Nairobi CBD"):
				nairobiIdx = i
			case contains(n, "Kisumu"):
				kisumuIdx = i
			case contains(n, "Mombasa"):
				mombasaIdx = i
			}
		}
		if nairobiIdx == -1 || kisumuIdx == -1 || mombasaIdx == -1 {
			t.Fatalf("missing seeded places in results: %v", names)
		}
		if nairobiIdx > kisumuIdx || kisumuIdx > mombasaIdx {
			t.Errorf("results not ordered by distance from Nairobi: %v", names)
		}
	})

	t.Run("radius filter excludes far places", func(t *testing.T) {
		results, err := repo.Nearby(ctx, -1.2921, 36.8219, 300_000, 10)
		if err != nil {
			t.Fatalf("Nearby: %v", err)
		}
		for _, r := range results {
			if contains(r.Name, "Mombasa") {
				t.Error("Mombasa returned within 300km radius of Nairobi")
			}
			if r.DistanceMeters > 300_000 {
				t.Errorf("%q distance %f exceeds radius", r.Name, r.DistanceMeters)
			}
		}
	})

	t.Run("limit is respected", func(t *testing.T) {
		results, err := repo.Nearby(ctx, -1.2921, 36.8219, 600_000, 2)
		if err != nil {
			t.Fatalf("Nearby: %v", err)
		}
		if len(results) > 2 {
			t.Errorf("got %d results with limit 2", len(results))
		}
	})
}

func TestRepository_Update_Integration(t *testing.T) {
	pool := testutil.TestPool(t)
	repo := places.NewRepository(pool)
	ctx := context.Background()

	p := &places.Place{Name: uniqueName(t, "Before"), PlaceType: places.TypeWard, Lat: 0.0, Lon: 37.0}
	if err := repo.Create(ctx, p); err != nil {
		t.Fatalf("Create: %v", err)
	}

	p.Name = uniqueName(t, "After")
	p.Lat = 0.5
	p.Lon = 37.5

	updated, err := repo.Update(ctx, p)
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if updated.Name != p.Name {
		t.Errorf("name = %q, want %q", updated.Name, p.Name)
	}
	if math.Abs(updated.Lat-0.5) > 1e-9 || math.Abs(updated.Lon-37.5) > 1e-9 {
		t.Errorf("location = (%v, %v), want (0.5, 37.5)", updated.Lat, updated.Lon)
	}
	if !updated.UpdatedAt.After(updated.CreatedAt) {
		t.Errorf("updated_at %v should be after created_at %v", updated.UpdatedAt, updated.CreatedAt)
	}
}

func TestRepository_NotFound_Integration(t *testing.T) {
	pool := testutil.TestPool(t)
	repo := places.NewRepository(pool)
	ctx := context.Background()

	if _, err := repo.ByID(ctx, uuid.New()); err != places.ErrNotFound {
		t.Errorf("ByID unknown id: err = %v, want ErrNotFound", err)
	}
	if _, err := repo.ByExternalCode(ctx, "no-such-code-"+uuid.NewString()); err != places.ErrNotFound {
		t.Errorf("ByExternalCode unknown code: err = %v, want ErrNotFound", err)
	}
}

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

func uniqueName(t *testing.T, base string) string {
	t.Helper()
	return base + " " + uuid.NewString()[:8]
}
