package preferences_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/google/uuid"

	"green-compass-backend/internal/places"
	"green-compass-backend/internal/preferences"
	"green-compass-backend/internal/testutil"
	"green-compass-backend/internal/users"
)

func TestService_SavedPlaces_Integration(t *testing.T) {
	pool := testutil.TestPool(t)
	ctx := context.Background()
	userSvc := users.NewService(users.NewRepository(pool))
	user, err := userSvc.Register(ctx, users.RegisterInput{Email: "saved-" + uuid.NewString() + "@test.invalid", DisplayName: "Saver", Password: "long-enough-password"})
	if err != nil {
		t.Fatalf("register user: %v", err)
	}
	placeRepo := places.NewRepository(pool)
	createPlace := func(name string) *places.Place {
		p := &places.Place{Name: name, PlaceType: places.TypeCommunity, Lat: -1.2, Lon: 36.8}
		if err := placeRepo.Create(ctx, p); err != nil {
			t.Fatalf("create place: %v", err)
		}
		return p
	}
	first, second := createPlace("First "+uuid.NewString()), createPlace("Second "+uuid.NewString())
	svc := preferences.NewService(preferences.NewRepository(pool))

	label := "  Home  "
	if err := svc.Save(ctx, user.ID, first.ID, &label); err != nil {
		t.Fatalf("save first: %v", err)
	}
	updatedLabel := "Work"
	if err := svc.Save(ctx, user.ID, first.ID, &updatedLabel); err != nil {
		t.Fatalf("upsert first: %v", err)
	}
	if err := svc.Save(ctx, user.ID, second.ID, nil); err != nil {
		t.Fatalf("save second: %v", err)
	}
	if err := svc.SetPrimary(ctx, user.ID, second.ID); err != nil {
		t.Fatalf("set primary: %v", err)
	}

	items, err := svc.List(ctx, user.ID)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(items) != 2 || items[0].ID != second.ID || !items[0].IsPrimary || items[1].Label == nil || *items[1].Label != "Work" {
		t.Errorf("saved places = %+v, want primary second and updated first label", items)
	}
	if err := svc.Unsave(ctx, user.ID, second.ID); err != nil {
		t.Fatalf("unsave primary: %v", err)
	}
	if err := svc.Unsave(ctx, user.ID, second.ID); !errors.Is(err, preferences.ErrNotSaved) {
		t.Errorf("repeat unsave = %v, want ErrNotSaved", err)
	}
	if err := svc.Save(ctx, user.ID, uuid.New(), nil); !errors.Is(err, preferences.ErrPlaceNotFound) {
		t.Errorf("unknown place = %v, want ErrPlaceNotFound", err)
	}

	for i := 0; i < preferences.MaxSavedPlaces-1; i++ {
		p := createPlace(fmt.Sprintf("Capacity %d %s", i, uuid.NewString()))
		if err := svc.Save(ctx, user.ID, p.ID, nil); err != nil {
			t.Fatalf("save capacity %d: %v", i, err)
		}
	}
	extra := createPlace("Over capacity " + uuid.NewString())
	if err := svc.Save(ctx, user.ID, extra.ID, nil); !errors.Is(err, preferences.ErrSavedPlaceLimit) {
		t.Errorf("over limit = %v, want ErrSavedPlaceLimit", err)
	}
}
