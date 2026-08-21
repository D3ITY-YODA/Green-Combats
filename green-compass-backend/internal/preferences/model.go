package preferences

import (
	"errors"
	"time"

	"green-compass-backend/internal/places"
)

var (
	ErrNotSaved        = errors.New("place is not saved")
	ErrSavedPlaceLimit = errors.New("saved place limit reached")
	ErrInvalidData     = errors.New("invalid saved place data")
	ErrPlaceNotFound   = errors.New("place not found")
)

const (
	MaxSavedPlaces = 50
	MaxLabelLen    = 200
)

type SavedPlace struct {
	places.Place
	Label     *string
	IsPrimary bool
	SavedAt   time.Time
}
