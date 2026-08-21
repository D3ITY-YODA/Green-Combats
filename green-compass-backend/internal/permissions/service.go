package permissions

import (
	"context"

	"github.com/google/uuid"

	"green-compass-backend/internal/organizations"
)

type Service struct {
	orgs *organizations.Service
}

func New(orgs *organizations.Service) *Service {
	return &Service{orgs: orgs}
}

func (s *Service) OrgRole(ctx context.Context, userID, orgID uuid.UUID) (string, error) {
	role, err := s.orgs.MemberRole(ctx, orgID, userID)
	if err != nil {
		return "", err
	}
	return role, nil
}

func (s *Service) RequireOrgRole(ctx context.Context, userID, orgID uuid.UUID, roles ...string) error {
	role, err := s.OrgRole(ctx, userID, orgID)
	if err != nil {
		return err
	}
	for _, allowed := range roles {
		if role == allowed {
			return nil
		}
	}
	return ErrPermissionDenied
}

func CanVerifyObservations(role string, orgStatus string) bool {
	if orgStatus != organizations.StatusVerified {
		return false
	}
	switch role {
	case organizations.RoleAdmin, organizations.RoleReviewer:
		return true
	}
	return false
}

func (s *Service) CheckCanVerifyObservations(ctx context.Context, userID, orgID uuid.UUID) error {
	org, err := s.orgs.ByID(ctx, orgID)
	if err != nil {
		return err
	}
	role, err := s.OrgRole(ctx, userID, orgID)
	if err != nil {
		return err
	}
	if !CanVerifyObservations(role, org.Status) {
		return ErrPermissionDenied
	}
	return nil
}
