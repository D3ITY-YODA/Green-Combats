package connectors

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type FetchRequest struct {
	PlaceID uuid.UUID
	Lat     float64
	Lon     float64
	From    time.Time
	To      time.Time
}

type RawObservation struct {
	ObservedAt  time.Time
	SourceURL   string
	ContentType string
	Payload     json.RawMessage
}

// Connector fetches provider values without normalizing their meaning or units.
// Every returned observation must use the provider's own timestamp.
type Connector interface {
	Code() string
	Fetch(ctx context.Context, request FetchRequest) ([]RawObservation, error)
}
