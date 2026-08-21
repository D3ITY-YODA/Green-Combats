package observations

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrNotFound    = errors.New("observation not found")
	ErrInvalidData = errors.New("invalid observation data")
	ErrNotAllowed  = errors.New("not allowed to access this observation")
)

type Category string

const (
	CategoryFlood        = "flood"
	CategoryDrought      = "drought"
	CategoryWaterQuality = "water_quality"
	CategoryCropDamage   = "crop_damage"
	CategoryAirQuality   = "air_quality"
	CategoryOther        = "other"
)

type Status string

const (
	StatusPending  = "pending"
	StatusVerified = "verified"
	StatusRejected = "rejected"
	StatusFlagged  = "flagged"
)

type Observation struct {
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

type Caller struct {
	UserID          uuid.UUID
	IsPlatformAdmin bool
}

func IsValidCategory(c string) bool {
	switch c {
	case CategoryFlood, CategoryDrought, CategoryWaterQuality, CategoryCropDamage, CategoryAirQuality, CategoryOther:
		return true
	}
	return false
}

func IsValidStatus(s string) bool {
	switch s {
	case StatusPending, StatusVerified, StatusRejected, StatusFlagged:
		return true
	}
	return false
}
