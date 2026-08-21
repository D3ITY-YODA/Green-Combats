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
	runID := uuid.New()
	lat := 55.0 + float64(runID[0])/10
	lon := -150.0 + float64(runID[1])/10

	seed := []places.Place{
		{Name: uniqueName(t, "nearest"), PlaceType: places.TypeCommunity, Lat: lat + 0.005, Lon: lon},
		{Name: uniqueName(t, "middle"), PlaceType: places.TypeWard, Lat: lat + 0.015, Lon: lon},
		{Name: uniqueName(t, "farthest"), PlaceType: places.TypeDistrict, Lat: lat + 0.030, Lon: lon},
		{Name: uniqueName(t, "custom"), PlaceType: places.TypeCustom, Lat: lat + 0.006, Lon: lon},
	}
	for i := range seed {
		if err := repo.Create(ctx, &seed[i]); err != nil {
			t.Fatalf("seed %s: %v", seed[i].Name, err)
		}
	}

	t.Run("orders by real distance and excludes custom places", func(t *testing.T) {
		results, err := repo.Nearby(ctx, lat, lon, 5_000, 100)
		if err != nil {
			t.Fatalf("Nearby: %v", err)
		}

		seeded := map[string]bool{}
		for _, p := range seed {
			seeded[p.Name] = true
		}

		var order []string
		for _, r := range results {
			if !seeded[r.Name] {
				continue
			}
			if r.PlaceType == places.TypeCustom {
				t.Errorf("custom place %q leaked into nearby results", r.Name)
			}
			if r.DistanceMeters <= 0 {
				t.Errorf("place %q has non-positive distance %f", r.Name, r.DistanceMeters)
			}
			order = append(order, r.Name)
		}

		nearestIdx, middleIdx, farthestIdx := -1, -1, -1
		for i, n := range order {
			switch {
			case contains(n, "nearest"):
				nearestIdx = i
			case contains(n, "middle"):
				middleIdx = i
			case contains(n, "farthest"):
				farthestIdx = i
			}
		}
		if nearestIdx == -1 || middleIdx == -1 || farthestIdx == -1 {
			t.Fatalf("missing seeded places in results: %v", order)
		}
		if nearestIdx > middleIdx || middleIdx > farthestIdx {
			t.Errorf("results not ordered by distance: %v", order)
		}
	})

	t.Run("radius filter excludes far places", func(t *testing.T) {
		results, err := repo.Nearby(ctx, lat, lon, 2_000, 100)
		if err != nil {
			t.Fatalf("Nearby: %v", err)
		}
		for _, r := range results {
			if contains(r.Name, "farthest") && seededName(seed, r.Name) {
				t.Error("farthest seeded place returned within 2km radius")
			}
			if r.DistanceMeters > 2_000 {
				t.Errorf("%q distance %f exceeds radius", r.Name, r.DistanceMeters)
			}
		}
	})

	t.Run("limit is respected", func(t *testing.T) {
		results, err := repo.Nearby(ctx, lat, lon, 5_000, 2)
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

func seededName(seed []places.Place, name string) bool {
	for _, p := range seed {
		if p.Name == name {
			return true
		}
	}
	return false
}
