package sources

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

const (
	CodeOpenMeteo  = "open_meteo"
	CodeNASAPower  = "nasa_power"
	CodeFreshwater = "freshwater"
	CodeSatellite  = "satellite"
	CodeCommunity  = "community"
)

var ErrNotFound = errors.New("data source not found")

type Source struct {
	ID           uuid.UUID
	Code         string
	DisplayName  string
	PollInterval time.Duration
	Enabled      bool
	Config       []byte
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
