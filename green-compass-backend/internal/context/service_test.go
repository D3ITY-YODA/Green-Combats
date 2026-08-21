package context

import (
	"context"
	"errors"
	"log/slog"
	"testing"
	"time"

	"github.com/google/uuid"

	"green-compass-backend/internal/places"
	"green-compass-backend/internal/preferences"
	"green-compass-backend/internal/updates"
)

type stubPreferences struct {
	places []preferences.SavedPlace
	err    error
}

func (s *stubPreferences) List(_ context.Context, _ uuid.UUID) ([]preferences.SavedPlace, error) {
	return s.places, s.err
}

type stubPlaces struct {
	place *places.Place
	err   error
}

func (s *stubPlaces) ByID(_ context.Context, _ uuid.UUID) (*places.Place, error) {
	return s.place, s.err
}

type stubToday struct {
	update *updates.Update
	err    error
}

func (s *stubToday) GetToday(_ context.Context, _ uuid.UUID) (*updates.Update, error) {
	return s.update, s.err
}

func TestResolveContext_PrimaryPlace(t *testing.T) {
	placeID := uuid.New()
	userID := uuid.New()
	label := "Home"

	stub := &stubPreferences{
		places: []preferences.SavedPlace{
			{
				Place:     places.Place{ID: placeID, Name: "Nairobi", PlaceType: "community", Lat: -1.2, Lon: 36.8},
				Label:     &label,
				IsPrimary: true,
				SavedAt:   time.Now(),
			},
		},
	}

	today := &stubToday{
		update: &updates.Update{
			ID:          uuid.New(),
			PlaceID:     placeID,
			PlaceName:   "Nairobi",
			ContentType: "today",
			Headline:    "Today: Rain expected in Nairobi",
			BodyText:    "Conditions in Nairobi require attention.",
			GeneratedAt: time.Now(),
		},
	}

	svc := NewService(stub, &stubPlaces{}, today, slog.Default())

	resp, err := svc.ResolveContext(context.Background(), ResolveContextRequest{UserID: userID})
	if err != nil {
		t.Fatalf("ResolveContext failed: %v", err)
	}

	if resp.Place.ID != placeID {
		t.Errorf("Place.ID: got %v, want %v", resp.Place.ID, placeID)
	}
	if resp.Place.Name != "Nairobi" {
		t.Errorf("Place.Name: got %q, want %q", resp.Place.Name, "Nairobi")
	}
	if resp.Place.Label == nil || *resp.Place.Label != "Home" {
		t.Errorf("Place.Label: got %v, want %q", resp.Place.Label, "Home")
	}
	if !resp.Place.IsPrimary {
		t.Error("Place.IsPrimary: got false, want true")
	}
	if resp.Update == nil {
		t.Fatal("Update: got nil, want non-nil")
	}
}

func TestResolveContext_ExplicitPlaceID(t *testing.T) {
	placeID := uuid.New()
	userID := uuid.New()

	stub := &stubPreferences{}
	placeStub := &stubPlaces{
		place: &places.Place{ID: placeID, Name: "Mombasa", PlaceType: "community", Lat: -4.0, Lon: 39.6},
	}
	today := &stubToday{
		update: &updates.Update{
			ID:          uuid.New(),
			PlaceID:     placeID,
			PlaceName:   "Mombasa",
			ContentType: "today",
			Headline:    "Today in Mombasa",
		},
	}

	svc := NewService(stub, placeStub, today, slog.Default())

	resp, err := svc.ResolveContext(context.Background(), ResolveContextRequest{
		UserID:  userID,
		PlaceID: &placeID,
	})
	if err != nil {
		t.Fatalf("ResolveContext failed: %v", err)
	}

	if resp.Place.ID != placeID {
		t.Errorf("Place.ID: got %v, want %v", resp.Place.ID, placeID)
	}
	if resp.Place.Name != "Mombasa" {
		t.Errorf("Place.Name: got %q, want %q", resp.Place.Name, "Mombasa")
	}
	if resp.Place.Label != nil {
		t.Errorf("Place.Label: got %v, want nil", resp.Place.Label)
	}
	if resp.Place.IsPrimary {
		t.Error("Place.IsPrimary: got true, want false")
	}
}

func TestResolveContext_NoPlace(t *testing.T) {
	stub := &stubPreferences{places: []preferences.SavedPlace{}}
	svc := NewService(stub, &stubPlaces{}, &stubToday{}, slog.Default())

	_, err := svc.ResolveContext(context.Background(), ResolveContextRequest{UserID: uuid.New()})

	if !errors.Is(err, ErrNoPlace) {
		t.Errorf("ResolveContext: got %v, want ErrNoPlace", err)
	}
}

func TestResolveContext_NoContent(t *testing.T) {
	placeID := uuid.New()
	stub := &stubPreferences{
		places: []preferences.SavedPlace{
			{Place: places.Place{ID: placeID, Name: "Kisumu", PlaceType: "community"}, IsPrimary: true},
		},
	}
	today := &stubToday{err: updates.ErrNotFound}
	svc := NewService(stub, &stubPlaces{}, today, slog.Default())

	resp, err := svc.ResolveContext(context.Background(), ResolveContextRequest{UserID: uuid.New()})
	if err != nil {
		t.Fatalf("ResolveContext failed: %v", err)
	}

	if resp.Place.Name != "Kisumu" {
		t.Errorf("Place.Name: got %q, want %q", resp.Place.Name, "Kisumu")
	}
	if resp.Update != nil {
		t.Errorf("Update: got %v, want nil", resp.Update)
	}
}

func TestResolveContext_PlaceNotFound(t *testing.T) {
	placeID := uuid.New()
	placeStub := &stubPlaces{err: places.ErrNotFound}
	svc := NewService(&stubPreferences{}, placeStub, &stubToday{}, slog.Default())

	_, err := svc.ResolveContext(context.Background(), ResolveContextRequest{
		UserID:  uuid.New(),
		PlaceID: &placeID,
	})

	if !errors.Is(err, ErrPlaceNotFound) {
		t.Errorf("ResolveContext: got %v, want ErrPlaceNotFound", err)
	}
}
