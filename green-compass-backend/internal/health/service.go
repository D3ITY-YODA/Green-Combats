package health

import (
	"context"
	"time"

	"green-compass-backend/pkg/clock"
)

type Service struct {
	version   string
	startedAt time.Time
	clock     clock.Clock
}

func NewService(version string, clk clock.Clock) *Service {
	return &Service{
		version:   version,
		startedAt: clk.Now(),
		clock:     clk,
	}
}

func (s *Service) Status(ctx context.Context) StatusResponse {
	return StatusResponse{
		Status:        StatusOK,
		Version:       s.version,
		UptimeSeconds: int64(s.clock.Now().Sub(s.startedAt).Seconds()),
	}
}
