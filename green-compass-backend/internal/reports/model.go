package reports

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrNotFound    = errors.New("report not found")
	ErrInvalidData = errors.New("invalid report data")
	ErrNotAllowed  = errors.New("not allowed to access reports")
)

type Caller struct {
	UserID          uuid.UUID
	IsPlatformAdmin bool
}

type Report struct {
	ID             uuid.UUID
	ReporterID     uuid.UUID
	PlaceID        *uuid.UUID
	Lat            *float64
	Lon            *float64
	Category       string
	Description    string
	PhotoObjectKey *string
	Status         string
	VerifiedBy     *uuid.UUID
	VerifiedAt     *time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
