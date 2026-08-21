package ingestion

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

const (
	StatusRunning   = "running"
	StatusSucceeded = "succeeded"
	StatusFailed    = "failed"
)

var ErrInvalidData = errors.New("invalid ingestion data")

type Run struct {
	ID                uuid.UUID
	SourceID          uuid.UUID
	PlaceID           uuid.UUID
	Status            string
	RequestedFrom     *time.Time
	RequestedTo       *time.Time
	StartedAt         time.Time
	FinishedAt        *time.Time
	LandedRecordCount int
	ErrorMessage      *string
}

type RawRecord struct {
	ID               uuid.UUID
	IngestionRunID   uuid.UUID
	SourceID         uuid.UUID
	PlaceID          uuid.UUID
	SourceObservedAt time.Time
	FetchedAt        time.Time
	SourceURL        string
	ContentType      string
	Payload          json.RawMessage
	PayloadChecksum  string
	CreatedAt        time.Time
}

func Checksum(payload []byte) string {
	digest := sha256.Sum256(payload)
	return hex.EncodeToString(digest[:])
}
