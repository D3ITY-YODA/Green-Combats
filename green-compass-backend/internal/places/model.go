package places

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrNotFound    = errors.New("place not found")
	ErrInvalidData = errors.New("invalid place data")
	ErrNotAllowed  = errors.New("not allowed to modify this place")
)

const (
	TypeCommunity = "community"
	TypeWard      = "ward"
	TypeDistrict  = "district"
	TypeCustom    = "custom"
)

type Place struct {
	ID           uuid.UUID
	Name         string
	PlaceType    string
	Lat          float64
	Lon          float64
	ExternalCode *string
	CreatedBy    *uuid.UUID
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type WithDistance struct {
	Place
	DistanceMeters float64
}

func IsValidType(t string) bool {
	switch t {
	case TypeCommunity, TypeWard, TypeDistrict, TypeCustom:
		return true
	}
	return false
}

func IsOfficialType(t string) bool {
	switch t {
	case TypeCommunity, TypeWard, TypeDistrict:
		return true
	}
	return false
}
