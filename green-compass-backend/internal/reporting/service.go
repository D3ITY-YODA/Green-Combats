package reporting

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type repo interface {
	DeliveryStats(ctx context.Context, placeID uuid.UUID, from, to string) (*DeliveryStats, error)
	ReportStats(ctx context.Context, orgID uuid.UUID, from, to string) (*ReportStats, error)
}

type Service struct {
	repo repo
}

func NewService(repo repo) *Service {
	return &Service{repo: repo}
}

// GetDeliveryStats returns notification delivery metrics for a place.
func (s *Service) GetDeliveryStats(ctx context.Context, placeID uuid.UUID, from, to string) (*DeliveryStats, error) {
	if from == "" || to == "" {
		return nil, fmt.Errorf("period_start and period_end are required")
	}
	return s.repo.DeliveryStats(ctx, placeID, from, to)
}

// GetReportStats returns community report metrics for an organization.
func (s *Service) GetReportStats(ctx context.Context, orgID uuid.UUID, from, to string) (*ReportStats, error) {
	if from == "" || to == "" {
		return nil, fmt.Errorf("period_start and period_end are required")
	}
	return s.repo.ReportStats(ctx, orgID, from, to)
}

// GetDashboard returns a combined dashboard view for an institutional user.
func (s *Service) GetDashboard(ctx context.Context, req QueryRequest) (*DashboardResponse, error) {
	resp := &DashboardResponse{}

	if req.PlaceID != nil {
		from := req.PeriodStart.Format("2006-01-02T15:04:05Z")
		to := req.PeriodEnd.Format("2006-01-02T15:04:05Z")

		delivery, err := s.repo.DeliveryStats(ctx, *req.PlaceID, from, to)
		if err != nil {
			return nil, fmt.Errorf("delivery stats: %w", err)
		}
		resp.Delivery = delivery
	}

	if req.OrgID != nil {
		from := req.PeriodStart.Format("2006-01-02T15:04:05Z")
		to := req.PeriodEnd.Format("2006-01-02T15:04:05Z")

		reports, err := s.repo.ReportStats(ctx, *req.OrgID, from, to)
		if err != nil {
			return nil, fmt.Errorf("report stats: %w", err)
		}
		resp.Reports = reports
	}

	return resp, nil
}
