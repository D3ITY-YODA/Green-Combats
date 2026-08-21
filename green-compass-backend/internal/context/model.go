package context

import (
	"errors"

	"github.com/google/uuid"

	"green-compass-backend/internal/updates"
)

var (
	ErrNoPlace       = errors.New("no saved place found for user")
	ErrPlaceNotFound = errors.New("place not found")
)

type PlaceInfo struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	PlaceType string    `json:"place_type"`
	Lat       float64   `json:"lat"`
	Lon       float64   `json:"lon"`
	Label     *string   `json:"label,omitempty"`
	IsPrimary bool      `json:"is_primary"`
}

type ContextResponse struct {
	Place  PlaceInfo       `json:"place"`
	Update *updates.Update `json:"update,omitempty"`
}

type ResolveContextRequest struct {
	UserID  uuid.UUID
	PlaceID *uuid.UUID
}
